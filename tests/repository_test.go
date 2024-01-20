package tests

import (
	"context"
	"testing"

	mongorepository "github.com/dmitrymomot/mongo-repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

func TestCreateAndFindByID(t *testing.T) {
	db := setupMongoDB(t)
	repo := mongorepository.NewMongoRepository[User](db, "users")

	user := User{Name: "John Doe", Email: "john@example.com"}

	// Test Create
	id, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// Test FindByID
	foundUser, err := repo.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
}
