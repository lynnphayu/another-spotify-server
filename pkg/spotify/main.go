package spotify

import (
	"net/http"
)

type Repository interface {
	CreateOrUpdateProfile(profile Profile) (*Profile, error)
	GetProfileWithEmail(email string) (*Profile, error)
	UpdateCredentials(email string, credentials *Credentials) (*Profile, error)
}

type HTTPClient interface {
	Request(methodType string, URL string, body map[string]interface{}, contentType string, auth string) (*http.Response, error)
}

// AuthService - functions implemented
type AuthService interface {
	Login(email string) (*Profile, error)
	AuthCallback(authorizationCode string) (*Profile, error)
	GetCredentials(authorizationCode string) (*Credentials, error)
	GetProfile(accessToken string) (*Profile, error)
	GetRecentlyPlayed(email string, limit int, before string, after string) (*[]byte, error)
	GetTracksAudioFeatures(email string, trackIDs []string) (*[]byte, error)
	GetValidToken(email string) (*Credentials, error)
	GetTopArtistsOrTracks(email string, top string, timeRange string, limit int, offset int) (*[]byte, error)
}

type Service struct {
	repository Repository
	httpClient HTTPClient
}

// New - instantiate auth instance
func New(repository Repository, httpClient HTTPClient) AuthService {
	return &Service{repository, httpClient}
}
