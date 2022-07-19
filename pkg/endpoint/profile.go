package endpoint

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Profile Interface
type Profile struct {
	ID              primitive.ObjectID `bson:"_id"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	Country         string             `bson:"country" json:"country"`
	DisplayName     string             `bson:"display_name" json:"display_name"`
	Email           string             `bson:"email" json:"email"`
	Product         string             `bson:"product" json:"product"`
	Type            string             `bson:"type" json:"type"`
	URI             string             `bson:"uri" json:"uri"`
	Credentials     Credentials        `bson:"credentials" json:"credentials"`
	ExplicitContent struct {
		FilerEnabled bool `bson:"filter_enabled" json:"filter_enabled"`
		FilterLocked bool `bson:"filter_locked" json:"filter_blocked"`
	} `bson:"explicit_content" json:"explicit_content"`
	ExternalUrls struct {
		Spotify string `bson:"spotify" json:"spotify"`
	} `bson:"external_urls" json:"external_urls"`
	Followers struct {
		Href  string `bson:"href" json:"href"`
		Total int    `bson:"total" json:"total"`
	} `bson:"followers" json:"followers"`
	Href      string `bson:"href" json:"href"`
	ProfileID string `bson:"profile_id" json:"id"`
	Images    []struct {
		URL string `bson:"url" json:"url"`
	} `bson:"images" json:"images"`
}

// Credentials Struct
type Credentials struct {
	AccessToken  string    `bson:"access_token" json:"access_token"`
	ExpiresIn    int       `bson:"expires_in" json:"expires_in"`
	RefreshToken string    `bson:"refresh_token" json:"refresh_token"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
	TokenType    string    `bson:"token_type" json:"token_type"`
	Scope        string    `bson:"scope" json:"scope"`
}
