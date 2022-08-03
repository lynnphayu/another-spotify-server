package spotify

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Profile Interface
type Profile struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt       time.Time          `bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `bson:"updated_at,omitempty"`
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
	CreatedAt    time.Time `bson:"created_at,omitempty"`
	UpdatedAt    time.Time `bson:"updated_at,omitempty"`
	TokenType    string    `bson:"token_type" json:"token_type"`
	Scope        string    `bson:"scope" json:"scope"`
}

// ALBUM
// {
// 	"album_type": "ALBUM",
// 	"artists": [
// 		{
// 			"external_urls": {
// 				"spotify": "https://open.spotify.com/artist/7m0BsF0t3K9WQFgKoPejfk"
// 			},
// 			"href": "https://api.spotify.com/v1/artists/7m0BsF0t3K9WQFgKoPejfk",
// 			"id": "7m0BsF0t3K9WQFgKoPejfk",
// 			"name": "ArrDee",
// 			"type": "artist",
// 			"uri": "spotify:artist:7m0BsF0t3K9WQFgKoPejfk"
// 		}
// 	],
// 	"external_urls": {
// 		"spotify": "https://open.spotify.com/album/2acy6L0ZXAGSHoW6TIVtyW"
// 	},
// 	"href": "https://api.spotify.com/v1/albums/2acy6L0ZXAGSHoW6TIVtyW",
// 	"id": "2acy6L0ZXAGSHoW6TIVtyW",
// 	"images": [
// 		{
// 			"height": 640,
// 			"url": "https://i.scdn.co/image/ab67616d0000b27380e060e9c13d966d13607a01",
// 			"width": 640
// 		},
// 		{
// 			"height": 300,
// 			"url": "https://i.scdn.co/image/ab67616d00001e0280e060e9c13d966d13607a01",
// 			"width": 300
// 		},
// 		{
// 			"height": 64,
// 			"url": "https://i.scdn.co/image/ab67616d0000485180e060e9c13d966d13607a01",
// 			"width": 64
// 		}
// 	],
// 	"name": "Pier Pressure",
// 	"release_date": "2022-03-18",
// 	"release_date_precision": "day",
// 	"total_tracks": 14,
// 	"type": "album",
// 	"uri": "spotify:album:2acy6L0ZXAGSHoW6TIVtyW"
// }

// ARTIST
// "artists": [
// 	{
// 		"external_urls": {
// 			"spotify": "https://open.spotify.com/artist/7m0BsF0t3K9WQFgKoPejfk"
// 		},
// 		"href": "https://api.spotify.com/v1/artists/7m0BsF0t3K9WQFgKoPejfk",
// 		"id": "7m0BsF0t3K9WQFgKoPejfk",
// 		"name": "ArrDee",
// 		"type": "artist",
// 		"uri": "spotify:artist:7m0BsF0t3K9WQFgKoPejfk"
// 	}
// ],

// external url struct
type ExternalUrls struct {
	Spotify string `bson:"spotify" json:"spotify"`
}
type Track struct {
	Explicit     bool         `bson:"explicit" json:"explicit"`
	ExternalUrls ExternalUrls `bson:"external_urls" json:"external_urls"`
	Href         string       `bson:"href" json:"href"`
	ID           string       `bson:"id" json:"id"`
	Name         string       `bson:"name" json:"name"`
	PreviewURL   string       `bson:"preview_url" json:"preview_url"`
	TrackNumber  int          `bson:"track_number" json:"track_number"`
	Type         string       `bson:"type" json:"type"`
	URI          string       `bson:"uri" json:"uri"`
	IsLocal      bool         `bson:"is_local" json:"is_local"`
	Popularity   int          `bson:"popularity" json:"popularity"`
}

type AudioFeatures struct {
	Danceability     float64 `bson:"danceability" json:"danceability"`
	Energy           float64 `bson:"energy" json:"energy"`
	Loudness         float64 `bson:"loudness" json:"loudness"`
	Speechiness      float64 `bson:"speechiness" json:"speechiness"`
	Acousticness     float64 `bson:"acousticness" json:"acousticness"`
	Instrumentalness float64 `bson:"instrumentalness" json:"instrumentalness"`
	Liveness         float64 `bson:"liveness" json:"liveness"`
	Valence          float64 `bson:"valence" json:"valence"`
	Tempo            float64 `bson:"tempo" json:"tempo"`
	ID               string  `bson:"id" json:"id"`
	DurationMS       int     `bson:"duration_ms" json:"duration_ms"`
	TimeSignature    int     `bson:"time_signature" json:"time_signature"`
	Key              int     `bson:"key" json:"key"`
	Mode             int     `bson:"mode" json:"mode"`
	Type             string  `bson:"type" json:"type"`
	URI              string  `bson:"uri" json:"uri"`
	AnalysisURL      string  `bson:"analysis_url" json:"analysis_url"`
	TrackHref        string  `bson:"track_href" json:"track_href"`
}
