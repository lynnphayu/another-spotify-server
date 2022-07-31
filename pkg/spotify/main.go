package spotify

import (
	"net/http"
	"time"
)

type Storage interface {
	CreateOrUpdateProfile(profile Profile) (*Profile, error)
	GetProfileWithEmail(email string) (*Profile, error)
	UpdateCredentials(email string, credentials *Credentials) (*Profile, error)
}

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Clear(key string) error
}

type HTTPClient interface {
	Request(methodType string, URL string, body map[string]interface{}, contentType string, auth string) (*http.Response, error)
}

type LoginResponse struct {
	Token   string   `json:"token"`
	Profile *Profile `json:"profile"`
}

type Services struct {
	Auth         AuthService
	PersonalInfo PersonalInfoService
}

// AuthService - functions implemented
type AuthService interface {
	Login(email string) (*Profile, error)
	AuthCallback(authorizationCode string) (*LoginResponse, error)
	GetCredentials(authorizationCode string) (*Credentials, error)
	GetValidToken(email string) (*Credentials, error)
}

type PersonalInfoService interface {
	GetProfile(accessToken string) (*Profile, error)
	GetRecentlyPlayed(email string, limit int, before string, after string) (*[]byte, error)
	GetTracksAudioFeatures(email string, trackIDs []string) (*[]byte, error)
	GetTopArtistsOrTracks(email string, top string, timeRange string, limit int, offset int) (*[]byte, error)
}

type Service struct {
	storage    Storage
	httpClient HTTPClient
	cache      Cache
}

// New - return map of both serivces
func NewServices(storage Storage, httpClient HTTPClient, cache Cache) Services {
	return Services{
		Auth:         &Service{storage, httpClient, cache},
		PersonalInfo: &Service{storage, httpClient, cache},
	}
}
