package spotify

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Create the Claims

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
	General      GeneralService
}

// AuthService - functions implemented
type AuthService interface {
	Login(email string) (*Profile, error)
	AuthCallback(authorizationCode string) (*LoginResponse, error)
	GetCredentials(authorizationCode string) (*Credentials, error)
	GetValidToken(email string) (*Credentials, error)
	GetProfileFromSpotify(accessToken string) (*Profile, error)
}

type PersonalInfoService interface {
	GetRecentlyPlayed(email string, limit int, before string, after string) (*[]byte, error)
	GetPersonalAudioFeatures(email string, timespan string) (*[]byte, error)
	GetTopArtistsOrTracks(email string, top string, timeRange string, limit int, offset int) (*[]byte, error)
	GetUserPlaylists(email string, limit int, offset int) (*[]byte, error)
}

type GeneralService interface {
	GetTracksAudioFeatures(email string, trackIDs []string) (*[]byte, error)
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
		General:      &Service{storage, httpClient, cache},
	}
}
