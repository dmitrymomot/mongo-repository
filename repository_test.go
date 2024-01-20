package mongorepository_test

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

func TestRepository(t *testing.T) {
	db := setupMongoDB(t)
	repo := mongorepository.NewMongoRepository[User](db, "users")

	var id string
	email := "john@example.com"
	user := User{Name: "John Doe", Email: email}

	// Create unique index for email field
	t.Run("CreateIndex", func(t *testing.T) {
		err := repo.CreateIndex(
			context.Background(),
			"email",
			mongorepository.Unique(true),
			mongorepository.Sparse(true),
		)
		require.NoError(t, err)
	})

	// Test Create
	id, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// Test create duplicate
	t.Run("CreateDuplicate", func(t *testing.T) {
		_, err := repo.Create(context.Background(), user)
		require.Error(t, err)
		require.ErrorIs(t, err, mongorepository.ErrDuplicate)
	})

	// Test Exists
	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(context.Background(), mongorepository.Eq("email", email))
		require.NoError(t, err)
		assert.True(t, exists)
	})

	// Test Count
	t.Run("Count", func(t *testing.T) {
		count, err := repo.Count(context.Background(), mongorepository.Eq("email", email))
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	// Test FindByID
	t.Run("FindByID", func(t *testing.T) {
		foundUser, err := repo.FindByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test FindByIDs
	t.Run("FindByIDs", func(t *testing.T) {
		users, err := repo.FindByIDs(context.Background(), id)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, user.Name, users[0].Name)
		assert.Equal(t, user.Email, users[0].Email)
	})

	// Test FindOne
	t.Run("FindOne", func(t *testing.T) {
		foundUser, err := repo.FindOneByFilter(context.Background(), mongorepository.Eq("email", email))
		require.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test Find Many by filter
	t.Run("FindMany", func(t *testing.T) {
		users, err := repo.FindManyByFilter(context.Background(), 0, 0, mongorepository.Eq("email", email))
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, user.Name, users[0].Name)
		assert.Equal(t, user.Email, users[0].Email)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		user.Name = "John Doe Updated"
		user.Email = email

		updCount, err := repo.Update(context.Background(), id, user)
		require.NoError(t, err)
		assert.Equal(t, int64(1), updCount)

		foundUser, err := repo.FindByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test UpdateMany
	t.Run("UpdateMany", func(t *testing.T) {
		user.Name = "John Doe UpdateMany"
		user.Email = email

		// Update all users with given email
		updCount, err := repo.UpdateMany(
			context.Background(),
			map[string]interface{}{"name": user.Name},
			mongorepository.Eq("email", user.Email),
		)
		require.NoError(t, err)
		assert.Equal(t, int64(1), updCount)

		foundUser, err := repo.FindByID(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		delCount, err := repo.Delete(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, int64(1), delCount)

		foundUser, err := repo.FindByID(context.Background(), id)
		require.ErrorIs(t, err, mongorepository.ErrNotFound)
		assert.Empty(t, foundUser)

		// Test delete non-existent user
		delCount, err = repo.Delete(context.Background(), id)
		require.ErrorIs(t, err, mongorepository.ErrNotFound)
		assert.Equal(t, int64(0), delCount)
	})

	// Test try to update non-existent user
	t.Run("UpdateNonExistent", func(t *testing.T) {
		updCount, err := repo.Update(context.Background(), id, user)
		require.ErrorIs(t, err, mongorepository.ErrNotFound)
		assert.Equal(t, int64(0), updCount)
	})
}
