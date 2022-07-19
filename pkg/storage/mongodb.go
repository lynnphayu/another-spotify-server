package storage

import (
	"context"
	"os"
	"sync"
	"time"
	"utilserver/pkg/spotify"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var instance *mongo.Client = nil
var doOnce sync.Once
var instanceError error

type Storage struct {
	client   *mongo.Client
	database *mongo.Database
}

// New - initialize Storage instance
func New() (*Storage, error) {
	storage := new(Storage)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING"))
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	database := client.Database(os.Getenv("MONGODB_DATABASE"))
	storage.database = database
	storage.client = client
	return storage, nil
}

//GetDBClient - create instance and Return client instance to work with
func (storage *Storage) GetDBClient(CONNECTIONSTRING string) (*mongo.Client, error) {
	//Perform connection creation operation only once.
	doOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			instanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			instanceError = err
		}
		instance = client
	})
	return instance, instanceError
}

func (storage *Storage) GetProfileWithEmail(email string) (*spotify.Profile, error) {
	var profile spotify.Profile
	collection := storage.database.Collection(os.Getenv("MONGODB_PROFILE_COLLECTION"))
	findErr := collection.FindOne(context.TODO(), map[string]string{"email": email}).Decode(&profile)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, findErr
	}
	return &profile, nil
}

// CreateProfile - create profile func
func (storage *Storage) CreateOrUpdateProfile(profile spotify.Profile) (*spotify.Profile, error) {
	var profileContainer spotify.Profile
	collection := storage.database.Collection(os.Getenv("MONGODB_PROFILE_COLLECTION"))
	err := collection.FindOne(context.TODO(), map[string]string{"email": profile.Email}).Decode(&profileContainer)
	if err != mongo.ErrNoDocuments {
		profile.Credentials.UpdatedAt = time.Now()
		profile.Credentials.CreatedAt = profileContainer.Credentials.CreatedAt
		profile.UpdatedAt = time.Now()
		err := collection.FindOneAndUpdate(context.TODO(),
			map[string]string{"email": profile.Email}, map[string]interface{}{"$set": profile},
		).Decode(&profileContainer)
		return &profileContainer, err
	}
	profile.Credentials.CreatedAt = time.Now()
	profile.CreatedAt = time.Now()
	profile.Credentials.UpdatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	profile.ID = primitive.NewObjectID()
	_, createError := collection.InsertOne(context.TODO(), profile)
	if createError != nil {
		return nil, createError
	}
	return &profile, createError
}

func (storage *Storage) UpdateCredentials(email string, credentials *spotify.Credentials) (*spotify.Profile, error) {
	var profileContainer spotify.Profile
	collection := storage.database.Collection(os.Getenv("MONGODB_PROFILE_COLLECTION"))
	updateParams := map[string]interface{}{
		"updated_at":               time.Now(),
		"credentials.access_token": credentials.AccessToken,
		"credentials.scope":        credentials.Scope,
		"credentials.exprires_in":  credentials.ExpiresIn,
		"credentials.updated_at":   time.Now(),
	}
	if credentials.RefreshToken != "" {
		updateParams["credentials.refresh_token"] = credentials.RefreshToken
	}
	err := collection.FindOneAndUpdate(
		context.TODO(),
		map[string]string{"email": email},
		map[string]interface{}{"$set": updateParams},
	).Decode(&profileContainer)
	return &profileContainer, err
}
