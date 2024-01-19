package mongorepository

import (
	"go.mongodb.org/mongo-driver/bson"
)

// FilterFunc is a function type that takes a BSON document and modifies it.
type FilterFunc func(bson.D) bson.D

// Eq creates an equality filter
func Eq(field string, value interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: value})
	}
}

// Gt creates a greater-than filter
func Gt(field string, value interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$gt": value}})
	}
}

// Lt creates a less-than filter
func Lt(field string, value interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$lt": value}})
	}
}

// In creates an "in" filter
func In(field string, values interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$in": values}})
	}
}

// And combines multiple filters with a logical AND
func And(filters ...FilterFunc) FilterFunc {
	return func(filter bson.D) bson.D {
		andFilters := make([]bson.E, 0)
		for _, f := range filters {
			andFilter := f(bson.D{})
			andFilters = append(andFilters, andFilter...)
		}
		return append(filter, bson.E{Key: "$and", Value: andFilters})
	}
}

// Or combines multiple filters with a logical OR
func Or(filters ...FilterFunc) FilterFunc {
	return func(filter bson.D) bson.D {
		orFilters := make([]bson.E, 0)
		for _, f := range filters {
			orFilter := f(bson.D{})
			orFilters = append(orFilters, orFilter...)
		}
		return append(filter, bson.E{Key: "$or", Value: orFilters})
	}
}

// Ne creates a not-equal filter
func Ne(field string, value interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$ne": value}})
	}
}

// Lte creates a less-than-or-equal filter
func Lte(field string, value interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$lte": value}})
	}
}

// Gte creates a greater-than-or-equal filter
func Gte(field string, value interface{}) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$gte": value}})
	}
}

// Exists checks if a field exists
func Exists(field string, exists bool) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$exists": exists}})
	}
}

// Regex creates a filter for regular expression matching
func Regex(field string, pattern string, options string) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: field, Value: bson.M{"$regex": pattern, "$options": options}})
	}
}

// TextSearch creates a full-text search filter
func TextSearch(searchTerm string) FilterFunc {
	return func(filter bson.D) bson.D {
		return append(filter, bson.E{Key: "$text", Value: bson.M{"$search": searchTerm}})
	}
}
