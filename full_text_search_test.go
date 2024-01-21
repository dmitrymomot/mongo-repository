package mongorepository_test

import (
	"context"
	"testing"

	mongorepository "github.com/dmitrymomot/mongo-repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFullTextSearch(t *testing.T) {
	type User struct {
		ID   primitive.ObjectID `bson:"_id,omitempty"`
		Name string             `bson:"name"`
		Bio  string             `bson:"bio,omitempty"`
		Tags []string           `bson:"tags,omitempty"`
	}

	db := setupMongoDB(t)
	repo := mongorepository.NewMongoRepository[User](db, "users")

	// Create unique index for email field
	require.NoError(t, repo.CreateFullTextIndex(
		context.Background(),
		map[string]int32{
			"name": 10,
			"bio":  5,
			"tags": 1,
		},
		"english",
	))

	// Create test users
	users := []User{
		{
			ID:   primitive.NewObjectID(),
			Name: "Test John Doe",
			Bio:  "Software Engineer",
			Tags: []string{"go", "mongodb", "developer", "test"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Jane Smith",
			Bio:  "Data Scientist",
			Tags: []string{"python", "machine learning", "data analysis"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Kayla TestJohnson",
			Bio:  "Frontend Developer",
			Tags: []string{"javascript", "react", "web development"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Alex Brown",
			Bio:  "Backend Test Developer",
			Tags: []string{"golang", "spring", "web development"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Emily Davis",
			Bio:  "UI/UX Designer Test",
			Tags: []string{"design", "user experience", "prototyping"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Michael Wilson",
			Bio:  "Golang Engineer",
			Tags: []string{"test", "docker", "kubernetes", "cloud", "microservices", "go-test"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Clark Thompson",
			Bio:  "Product Manager",
			Tags: []string{"agile", "scrum", "product development", "testify"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "David Lee",
			Bio:  "Web Developer",
			Tags: []string{"javascript", "node.js", "react", "web development"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Jessica Martinez",
			Bio:  "Mobile App Developer",
			Tags: []string{"android", "java", "kotlin"},
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Ryan Clark",
			Bio:  "Data Engineer",
			Tags: []string{"big data", "hadoop", "test"},
		},
	}
	for _, user := range users {
		// Test Create
		id, err := repo.Create(context.Background(), user)
		require.NoError(t, err)
		require.NotEmpty(t, id)
	}

	// Test full text search
	t.Run("Search", func(t *testing.T) {
		users, err := repo.Search(context.Background(), 0, 10, "test")
		require.NoError(t, err)
		assert.Len(t, users, 5)
		assert.Equal(t, "Test John Doe", users[0].Name)
		assert.Equal(t, "Alex Brown", users[1].Name)
		assert.Equal(t, "Emily Davis", users[2].Name)
		assert.Equal(t, "Michael Wilson", users[3].Name)
		assert.Equal(t, "Ryan Clark", users[4].Name)
	})

	// Test full text search with exclusion
	t.Run("SearchExclude", func(t *testing.T) {
		users, err := repo.Search(context.Background(), 0, 10, "web -test")
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "David Lee", users[0].Name)
		assert.Equal(t, "Kayla TestJohnson", users[1].Name)
	})
}
