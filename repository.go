package mongorepository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository is an interface that defines the methods for interacting with a MongoDB collection.
type Repository[T any] interface {
	// CreateIndex creates an index in the MongoDB collection based on the specified key and options.
	// It takes a context.Context as the first argument, the key for the index as the second argument,
	// and optional IndexOption(s) as the third argument(s).
	// The function returns an error if the index creation fails.
	CreateIndex(ctx context.Context, key interface{}, opts ...IndexOption) error

	// Create inserts a new document into the MongoDB collection.
	// It takes a context.Context and a model of type T as input parameters.
	// It returns the ID of the newly created document as a string and an error, if any.
	Create(ctx context.Context, model T) (string, error)

	// FindByID retrieves a document from the MongoDB collection by its ID.
	// It takes a context.Context and the ID of the document as parameters.
	// It returns the retrieved document of type T and an error, if any.
	FindByID(ctx context.Context, id string) (T, error)

	// FindByIDs retrieves multiple documents from the MongoDB collection by their IDs.
	// It takes a context.Context and a slice of IDs as parameters.
	// It returns a slice of documents of type T and an error, if any.
	FindByIDs(ctx context.Context, ids ...string) ([]T, error)

	// Update updates a document in the MongoDB collection with the specified ID.
	// It takes a context, ID string, and model as input parameters.
	// It returns the number of modified documents and an error, if any.
	Update(ctx context.Context, id string, model T) (int64, error)

	// UpdateMany updates multiple documents in the MongoDB collection based on the provided filters.
	// It takes a context.Context, a map of update fields, and optional filter functions as parameters.
	// The update fields specify the changes to be made to the documents.
	// The filter functions are used to build the filter for selecting the documents to be updated.
	// It returns the number of documents modified and an error if any.
	UpdateMany(ctx context.Context, update map[string]interface{}, filters ...FilterFunc) (int64, error)

	// Delete deletes a document from the MongoDB collection based on the provided ID.
	// It returns the number of deleted documents and an error, if any.
	Delete(ctx context.Context, id string) (int64, error)

	// DeleteMany deletes multiple documents from the MongoDB collection based on the provided filters.
	// It returns the number of deleted documents and an error, if any.
	DeleteMany(ctx context.Context, filters ...FilterFunc) (int64, error)

	// FindManyByFilter retrieves multiple documents from the collection based on the provided filters.
	// It allows skipping a certain number of documents and limiting the number of documents to be returned.
	// The filters are applied in the order they are passed.
	// If no documents match the filters, it returns an error with the ErrNotFound error code.
	// If an error occurs during the retrieval process, it returns an error with the ErrFailedToFindManyByFilter error code.
	// The function returns a slice of documents of type T and an error.
	FindManyByFilter(ctx context.Context, skip int64, limit int64, filters ...FilterFunc) ([]T, error)

	// FindOneByFilter finds a single document in the collection based on the provided filters.
	// It accepts one or more FilterFunc functions that modify the filter criteria.
	// The function returns the found document of type T and an error, if any.
	// If no document is found, it returns an error of type ErrNotFound.
	// If an error occurs during the find operation, it returns the error.
	FindOneByFilter(ctx context.Context, filters ...FilterFunc) (T, error)

	// Exists checks if a document exists in the collection based on the provided filters.
	// It accepts one or more FilterFunc functions that modify the filter criteria.
	// The function returns true if a document exists and false otherwise.
	// If an error occurs during the find operation, it returns the error.
	Exists(ctx context.Context, filters ...FilterFunc) (bool, error)

	// Count returns the number of documents in the collection based on the provided filters.
	// It accepts one or more FilterFunc functions that modify the filter criteria.
	// The function returns the number of documents and an error, if any.
	Count(ctx context.Context, filters ...FilterFunc) (int64, error)
}

// mongoRepository is a generic struct that represents a MongoDB repository.
// It holds a reference to a mongo.Collection, which is used to interact with the MongoDB database.
type mongoRepository[T any] struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a new instance of the mongoRepository[T] struct.
// It takes a mongo.Database and a collectionName as parameters and returns a pointer to the mongoRepository[T] struct.
// The mongoRepository[T] struct represents a repository for working with a specific MongoDB collection.
// The collection field of the struct is initialized with the specified collectionName from the provided database.
func NewMongoRepository[T any](db *mongo.Database, collectionName string) *mongoRepository[T] {
	return &mongoRepository[T]{collection: db.Collection(collectionName)}
}

// CreateIndex creates an index in the MongoDB collection based on the specified key and options.
// It takes a context.Context as the first argument, the key for the index as the second argument,
// and optional IndexOption(s) as the third argument(s).
// The function returns an error if the index creation fails.
func (r *mongoRepository[T]) CreateIndex(ctx context.Context, key string, opts ...IndexOption) error {
	indexOpts := options.Index()
	for _, opt := range opts {
		opt(indexOpts)
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: key, Value: 1}},
		Options: indexOpts,
	}

	if _, err := r.collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return errors.Join(ErrFailedToCreateIndex, err)
	}
	return nil
}

// Create inserts a new document into the MongoDB collection.
// It takes a context.Context and a model of type T as input parameters.
// It returns the ID of the newly created document as a string and an error, if any.
func (r *mongoRepository[T]) Create(ctx context.Context, model T) (string, error) {
	result, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		// Handle duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return "", errors.Join(ErrFailedToCreate, ErrDuplicate, err)
		}
		return "", errors.Join(ErrFailedToCreate, err)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.Join(ErrInvalidDocumentID, err)
	}
	return oid.Hex(), nil
}

// FindByID retrieves a document from the MongoDB collection by its ID.
// It takes a context.Context and the ID of the document as parameters.
// It returns the retrieved document of type T and an error, if any.
func (r *mongoRepository[T]) FindByID(ctx context.Context, id string) (T, error) {
	var result T
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.Join(ErrFailedToFindByID, ErrInvalidDocumentID, err)
	}
	filter := bson.M{"_id": objID}
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, errors.Join(ErrFailedToFindByID, ErrNotFound, err)
		}
		return result, errors.Join(ErrFailedToFindByID, err)
	}
	return result, nil
}

// FindByIDs retrieves multiple documents from the MongoDB collection by their IDs.
// It takes a context.Context and a slice of IDs as parameters.
// It returns a slice of documents of type T and an error, if any.
func (r *mongoRepository[T]) FindByIDs(ctx context.Context, ids ...string) ([]T, error) {
	// Convert string IDs to ObjectIDs
	objIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.Join(ErrFailedToFindByIDs, ErrInvalidDocumentID, err)
		}
		objIDs[i] = objID
	}

	// Build the query filter
	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	// Find documents
	cursor, err := r.collection.Find(ctx, filter, options.Find())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Join(ErrFailedToFindByIDs, ErrNotFound, err)
		}
		return nil, errors.Join(ErrFailedToFindByIDs, err)
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var element T
		if err := cursor.Decode(&element); err != nil {
			return nil, errors.Join(ErrFailedToFindByIDs, err)
		}
		results = append(results, element)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Join(ErrFailedToFindByIDs, err)
	}
	if len(results) == 0 {
		return nil, errors.Join(ErrFailedToFindByIDs, ErrNotFound)
	}
	return results, nil
}

// Update updates a document in the MongoDB collection with the specified ID.
// It takes a context, ID string, and model as input parameters.
// It returns the number of modified documents and an error, if any.
func (r *mongoRepository[T]) Update(ctx context.Context, id string, model T) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.Join(ErrFailedToFindByID, ErrInvalidDocumentID, err)
	}
	update := bson.M{"$set": model}
	result, err := r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.Join(ErrFailedToUpdate, ErrNotFound, err)
		}
		return 0, errors.Join(ErrFailedToUpdate, err)
	}
	if result.MatchedCount == 0 {
		return 0, errors.Join(ErrFailedToUpdate, ErrNotFound)
	}
	return result.ModifiedCount, nil
}

// UpdateMany updates multiple documents in the MongoDB collection based on the provided filters.
// It takes a context.Context, a map of update fields, and optional filter functions as parameters.
// The update fields specify the changes to be made to the documents.
// The filter functions are used to build the filter for selecting the documents to be updated.
// It returns the number of documents modified and an error if any.
func (r *mongoRepository[T]) UpdateMany(ctx context.Context, update map[string]interface{}, filters ...FilterFunc) (int64, error) {
	// Build the filter
	filter := bson.D{}
	for _, f := range filters {
		filter = f(filter)
	}

	// Prepare the update document
	updateDoc := bson.M{"$set": update}

	// Perform the update
	result, err := r.collection.UpdateMany(ctx, filter, updateDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.Join(ErrFailedToUpdateMany, ErrNotFound, err)
		}
		return 0, errors.Join(ErrFailedToUpdateMany, err)
	}
	return result.ModifiedCount, nil
}

// Delete deletes a document from the MongoDB collection based on the provided ID.
// It returns the number of deleted documents and an error, if any.
func (r *mongoRepository[T]) Delete(ctx context.Context, id string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.Join(ErrFailedToFindByID, ErrInvalidDocumentID, err)
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.Join(ErrFailedToDelete, ErrNotFound, err)
		}
		return 0, errors.Join(ErrFailedToDelete, err)
	}
	if result.DeletedCount == 0 {
		return 0, errors.Join(ErrFailedToDelete, ErrNotFound)
	}
	return result.DeletedCount, nil
}

// DeleteMany deletes multiple documents from the MongoDB collection based on the provided filters.
// It returns the number of deleted documents and an error, if any.
func (r *mongoRepository[T]) DeleteMany(ctx context.Context, filters ...FilterFunc) (int64, error) {
	filter := bson.D{}
	for _, f := range filters {
		filter = f(filter)
	}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.Join(ErrFailedToDeleteMany, ErrNotFound, err)
		}
		return 0, errors.Join(ErrFailedToDeleteMany, err)
	}
	return result.DeletedCount, nil
}

// FindManyByFilter retrieves multiple documents from the collection based on the provided filters.
// It allows skipping a certain number of documents and limiting the number of documents to be returned.
// The filters are applied in the order they are passed.
// If no documents match the filters, it returns an error with the ErrNotFound error code.
// If an error occurs during the retrieval process, it returns an error with the ErrFailedToFindManyByFilter error code.
// The function returns a slice of documents of type T and an error.
func (r *mongoRepository[T]) FindManyByFilter(ctx context.Context, skip int64, limit int64, filters ...FilterFunc) ([]T, error) {
	filter := bson.D{}
	for _, f := range filters {
		filter = f(filter)
	}
	if limit == 0 {
		limit = 10
	}
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)
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

// FindOneByFilter finds a single document in the collection based on the provided filters.
// It accepts one or more FilterFunc functions that modify the filter criteria.
// The function returns the found document of type T and an error, if any.
// If no document is found, it returns an error of type ErrNotFound.
// If an error occurs during the find operation, it returns the error.
func (r *mongoRepository[T]) FindOneByFilter(ctx context.Context, filters ...FilterFunc) (T, error) {
	filter := bson.D{}
	for _, f := range filters {
		filter = f(filter)
	}
	var result T
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, errors.Join(ErrFailedToFindOneByFilter, ErrNotFound, err)
		}
		return result, errors.Join(ErrFailedToFindOneByFilter, err)
	}
	return result, nil
}

// Exists checks if a document exists in the collection based on the provided filters.
// It accepts one or more FilterFunc functions that modify the filter criteria.
// The function returns true if a document exists and false otherwise.
// If an error occurs during the find operation, it returns the error.
func (r *mongoRepository[T]) Exists(ctx context.Context, filters ...FilterFunc) (bool, error) {
	filter := bson.D{}
	for _, f := range filters {
		filter = f(filter)
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.Join(ErrFailedToFindOneByFilter, err)
	}
	return count > 0, nil
}

// Count returns the number of documents in the collection based on the provided filters.
// It accepts one or more FilterFunc functions that modify the filter criteria.
// The function returns the number of documents and an error, if any.
func (r *mongoRepository[T]) Count(ctx context.Context, filters ...FilterFunc) (int64, error) {
	filter := bson.D{}
	for _, f := range filters {
		filter = f(filter)
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.Join(ErrFailedToFindOneByFilter, err)
	}
	return count, nil
}
