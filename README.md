# mongo-repository

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/mongo-repository)](https://github.com/dmitrymomot/mongo-repository)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/mongo-repository.svg)](https://pkg.go.dev/github.com/dmitrymomot/mongo-repository)
[![License](https://img.shields.io/github/license/dmitrymomot/mongo-repository)](https://github.com/dmitrymomot/mongo-repository/blob/main/LICENSE)

[![Tests](https://github.com/dmitrymomot/mongo-repository/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/mongo-repository/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/mongo-repository/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/mongo-repository/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/mongo-repository/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/mongo-repository/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/mongo-repository)](https://goreportcard.com/report/github.com/dmitrymomot/mongo-repository)

This Go package provides a generic, extensible MongoDB repository with advanced CRUD operations, flexible querying capabilities, and support for common MongoDB patterns. It is designed to simplify the interaction with MongoDB using Go's standard library and the official MongoDB Go driver.

## Features

- Generic CRUD operations for MongoDB documents.
- Query helpers for complex filters and full-text search.
- Support for finding documents by multiple IDs.
- Custom index creation with various index options.
- Batch update capabilities.

## Installation

To use this MongoDB repository package, you need to have Go installed on your machine. The package can be installed using the following Go command:

```bash
go get -u github.com/dmitrymomot/mongo-repository
```

## Usage

Import the package into your Go file:

```go
import "github.com/dmitrymomot/mongo-repository"
```

### Basic CRUD Operations
Here's a quick example of how you can use the repository for CRUD operations:

```go
func main() {
    // Initialize your MongoDB client and context
    // ...

    // Create a new repository for your model
    userRepo := repository.NewMongoRepository[User](db, "users")

    // Use the repository for various operations
    // ...
}
```

### Advanced Querying
The package includes a filter builder to create complex queries easily:

```go
func findUsers(repo repository.Repository[User]) {
    ctx := context.TODO()
    filter := repository.And(
        repository.Gt("age", 30),
        repository.Eq("status", "active"),
    )
    users, err := repo.FindManyByFilter(ctx, 0, 10, filter)
    // Handle err and work with users
}
```

### Full-Text Search

The package includes a full-text search builder to create text queries easily. The text search query uses the [MongoDB text search](https://docs.mongodb.com/manual/text-search/) feature.

#### Creating a Text Index in Your Collection

When you want to create a text index, specify the field and use the TextIndex option. This field typically stores the text you want to search. MongoDB uses this field to determine if a document is a match. You can also specify the weights for each field to control the relative score of each field.

```go
// Create a text index
err := repo.CreateIndex( context.TODO(),  "name",  mongorepository.TextIndex(
    mongorepository.NewTextIndexConfig(
        map[string]int32{
            "name": 10,
            "bio":  5,
            "tags": 1,
        },
    ),
))
```

#### Document Structure

Ensure your documents have a field (like name) that stores the text you want to search. This field is used by MongoDB to determine if a document is a match.

```go
type User struct {
    ID   primitive.ObjectID `bson:"_id,omitempty"`
    Name string             `bson:"name"`
    Bio  string             `bson:"bio"`
    Tags []string           `bson:"tags"`
    // Other fields...
}
```

#### Searching for Documents

To search for documents, use the Text helper to create a text search query. The text search query uses the [MongoDB text search](https://docs.mongodb.com/manual/text-search/) feature.

```go
users, err := repo.FindManyByFilter(ctx, 0, 10, mongorepository.TextSearch("John"))
```

### TTL Index

To create an index with a Time-To-Live (TTL) in MongoDB, which automatically deletes documents after a certain amount of time, you need to specify the TTL when creating the index. MongoDB uses a special background task that runs periodically to remove expired documents.

#### Creating a TTL Index in Your Collection

When you want to create an index with a TTL, specify the field and use the TTL option. This field typically stores the creation time of the document and should be of BSON date type.

```go
// Create an index with a TTL of 24 hours
err := repo.CreateIndex(context.TODO(), "createdAt", mongorepository.TTL(24*time.Hour))
```

#### Document Structure

Ensure your documents have a field (like createdAt) that stores the time when the document was created. This field is used by MongoDB to determine if a document is expired.

```go
type YourDocumentType struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    CreatedAt time.Time          `bson:"createdAt"`
    // Other fields...
}
```

#### Notes

- MongoDB runs a background task every 60 seconds to remove expired documents, so there may be a slight delay before documents are actually deleted.
- This approach is commonly used for data that needs to be retained only for a specific duration, such as logs, temporary data, or session information.

## Contributing

Contributions to the `mongo-repository` package are welcome! Here are some ways you can contribute:

- Reporting bugs
- Additional tests cases
- Suggesting enhancements
- Submitting pull requests
- Sharing the love by telling others about this project

## License

This project is licensed under the [Apache 2.0](LICENSE) - see the `LICENSE` file for details.


