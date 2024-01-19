package mongorepository

import "go.mongodb.org/mongo-driver/mongo/options"

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
