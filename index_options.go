package mongorepository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexOption wraps the MongoDB IndexOptions for extensibility and ease of use
type IndexOption func(*options.IndexOptions)

// Unique specifies the unique option for an index
func Unique(unique bool) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetUnique(unique)
	}
}

// Sparse specifies the sparse option for an index
func Sparse(sparse bool) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetSparse(sparse)
	}
}

// ExpireAfterSeconds specifies the expireAfterSeconds option for an index
func ExpireAfterSeconds(expireAfterSeconds int32) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetExpireAfterSeconds(expireAfterSeconds)
	}
}

// TTL sets the Time-To-Live for an index
func TTL(duration time.Duration) IndexOption {
	return func(opts *options.IndexOptions) {
		ttl := int32(duration.Seconds())
		opts.SetExpireAfterSeconds(ttl)
	}
}

// Name specifies the name option for an index
func Name(name string) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetName(name)
	}
}

// PartialFilterExpression specifies the partialFilterExpression option for an index
func PartialFilterExpression(partialFilterExpression interface{}) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetPartialFilterExpression(partialFilterExpression)
	}
}

// SetCollation specifies the collation option for an index
func SetCollation(collation *options.Collation) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetCollation(collation)
	}
}

// SetWildcardProjection specifies the wildcardProjection option for an index
func SetWildcardProjection(wildcardProjection interface{}) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetWildcardProjection(wildcardProjection)
	}
}

// SetHidden specifies the hidden option for an index
func SetHidden(hidden bool) IndexOption {
	return func(opts *options.IndexOptions) {
		opts.SetHidden(hidden)
	}
}

// TextIndexConfig configures the text index
type TextIndexConfig struct {
	Fields      map[string]int32 // Fields with weights
	DefaultLang string           // Default language for the index
	Name        string           // Optional custom name for the index
}

// NewTextIndexConfig creates a new text index config with specified fields and weights
// Panics if no fields are provided
func NewTextIndexConfig(fields map[string]int32) TextIndexConfig {
	if len(fields) == 0 {
		panic("fields must be provided")
	}

	return TextIndexConfig{
		Fields:      fields,
		DefaultLang: "english",
		Name:        "default_fts_index",
	}
}

// TextIndex creates a text index option with specified fields and weights
func TextIndex(config TextIndexConfig) IndexOption {
	return func(opts *options.IndexOptions) {
		// Set weights for fields
		weights := bson.D{}
		for field, weight := range config.Fields {
			weights = append(weights, bson.E{Key: field, Value: weight})
		}

		opts.SetWeights(weights)

		// Set default language if provided
		if config.DefaultLang != "" {
			opts.SetDefaultLanguage(config.DefaultLang)
		}

		// Set custom index name if provided
		if config.Name != "" {
			opts.SetName(config.Name)
		}
	}
}
