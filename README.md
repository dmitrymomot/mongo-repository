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
go get -u path/to/your/mongodb/repository/package
```

## Usage

Import the package into your Go file:

```go
import "path/to/your/mongodb/repository/package"
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

## Contributing

Contributions to the `mongo-repository` package are welcome! Here are some ways you can contribute:

- Reporting bugs
- Additional tests cases
- Suggesting enhancements
- Submitting pull requests
- Sharing the love by telling others about this project

## License

This project is licensed under the [Apache 2.0](LICENSE) - see the `LICENSE` file for details.


