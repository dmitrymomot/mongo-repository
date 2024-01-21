package mongorepository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateFullTextIndex creates a full-text index in the MongoDB collection based on the specified key and options.
// It takes a context.Context as the first argument, the key for the index as the second argument,
// and optional IndexOption(s) as the third argument(s).
// The function returns an error if the index creation fails.
func (r *mongoRepository[T]) CreateFullTextIndex(ctx context.Context, keys map[string]int32, lang string) error {
	// Build the index keys and weights
	idxKeys := make(bson.D, 0, len(keys))
	weights := make(bson.D, 0, len(keys))
	for k, w := range keys {
		idxKeys = append(idxKeys, bson.E{Key: k, Value: "text"})
		weights = append(weights, bson.E{Key: k, Value: w})
	}
	if lang == "" {
		lang = "english"
	}

	// Set index options
	idxOpt := options.Index()
	idxOpt.SetWeights(weights)
	idxOpt.SetDefaultLanguage(lang)
	idxOpt.SetName(fmt.Sprintf("%s_fts_index", lang))
	idxOpt.SetSparse(true)

	// Create the index
	indexModel := mongo.IndexModel{
		Keys:    idxKeys,
		Options: idxOpt,
	}

	// Create the index
	if _, err := r.collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return errors.Join(ErrFailedToCreateIndex, err)
	}
	return nil
}

// Search finds documents in the collection based on the provided search term.
// It allows skipping a certain number of documents and limiting the number of documents to be returned.
// The function returns a slice of documents of type T and an error.
func (r *mongoRepository[T]) Search(ctx context.Context, skip, limit int64, searchTerm string) ([]T, error) {
	filter := bson.M{"$text": bson.M{"$search": searchTerm}}
	if limit == 0 {
		limit = 10
	}
	// Set the find options
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}}).
		SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	// Find documents
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Join(ErrFailedToFindManyByFilter, ErrNotFound, err)
		}
		return nil, errors.Join(ErrFailedToFindManyByFilter, err)
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var element T
		if err := cursor.Decode(&element); err != nil {
			return nil, errors.Join(ErrFailedToFindManyByFilter, err)
		}
		results = append(results, element)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Join(ErrFailedToFindManyByFilter, err)
	}
	if len(results) == 0 {
		return nil, errors.Join(ErrFailedToFindManyByFilter, ErrNotFound)
	}

	return results, nil
}
