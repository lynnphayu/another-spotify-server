package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Login - login logic
func (service *Service) Login(email string) (*Profile, error) {
	profile, profileErr := service.storage.GetProfileWithEmail(email)
	return profile, profileErr
}

func (service *Service) GetCredentials(authorizationCode string) (*Credentials, error) {
	secretToken := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))
	resp, err := service.httpClient.Request(
		"POST",
		os.Getenv("SPOTIFY_TOKEN_GENERATOR_ENTPOINT"),
		map[string]interface{}{
			"code":         authorizationCode,
			"redirect_uri": os.Getenv("REDIRECT_URL"),
			"grant_type":   "authorization_code",
		},
		"application/x-www-form-urlencoded",
		"Basic "+secretToken,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	returnBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var credentials Credentials
	json.Unmarshal(returnBody, &credentials)
	return &credentials, nil
}

func (service *Service) GetProfile(accessToken string) (*Profile, error) {
	profileResp, err := service.httpClient.Request(
		"GET",
		os.Getenv("SPOTIFY_PROFILE_URL"), nil,
		"application/json",
		"Bearer "+accessToken,
	)

	if err != nil {
		return nil, err
	}
	defer profileResp.Body.Close()

	profileResponse, err := ioutil.ReadAll(profileResp.Body)
	if err != nil {
		return nil, err
	}

	var profile Profile
	json.Unmarshal(profileResponse, &profile)
	return &profile, err
}

// AuthCallback - callback function when spotify hit the  authorization endpoint
func (service *Service) AuthCallback(authorizationCode string) (*LoginResponse, error) {
	credentials, err := service.GetCredentials(authorizationCode)

	if err != nil {
		return nil, err
	}

	profile, err := service.GetProfile(credentials.AccessToken)
	if err != nil {
		return nil, err
	}

	profile.Credentials = *credentials
	profileArtifact, createError := service.storage.CreateOrUpdateProfile(*profile)
	if createError != nil {
		return nil, createError
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": profile.Email,
	})
	tokenString, signedErr := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return nil, signedErr
	}
	respose := LoginResponse{
		Profile: profileArtifact,
		Token:   tokenString,
	}
	return &respose, nil
}

// GetValidToken - return credentials with valid token meaning if token is expred, token will be refreshed
func (service *Service) GetValidToken(email string) (*Credentials, error) {
	if email == "" {
		return nil, errors.New("email expected")
	}
	profile, err := service.storage.GetProfileWithEmail(email)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("no profile found for current user")
	}
	remaingTokenTime := 3600 - time.Since(profile.Credentials.UpdatedAt).Seconds()
	if remaingTokenTime <= 10 {
		refreshCredentials, err := service.RefreshToken(profile.Credentials.RefreshToken)
		if err != nil {
			return nil, err
		}
		_, updateErr := service.storage.UpdateCredentials(email, refreshCredentials)
		if updateErr != nil {
			return nil, updateErr
		}
		return refreshCredentials, nil
	}
	return &profile.Credentials, nil
}

func (service *Service) RefreshToken(refreshToken string) (*Credentials, error) {
	secretToken := base64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))
	profileResp, err := service.httpClient.Request(
		"POST",
		os.Getenv("SPOTIFY_TOKEN_GENERATOR_ENTPOINT"),
		map[string]interface{}{
			"refresh_token": refreshToken,
			"grant_type":    "refresh_token",
		},
		"application/x-www-form-urlencoded",
		"Basic "+secretToken,
	)

	if err != nil {
		return nil, err
	}
	defer profileResp.Body.Close()

	refreshTokenResp, err := ioutil.ReadAll(profileResp.Body)
	if err != nil {
		return nil, err
	}

	var refreshTokenPayload Credentials
	json.Unmarshal(refreshTokenResp, &refreshTokenPayload)
	return &refreshTokenPayload, nil
}
