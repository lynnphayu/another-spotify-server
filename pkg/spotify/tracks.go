package spotify

import (
	"io/ioutil"
	"os"
)

func (service *Service) GetTracksAudioFeatures(email string, trackIDs []string) (*[]byte, error) {
	profile, err := service.repository.GetProfileWithEmail(email)
	if err != nil {
		return nil, err
	}
	URL := os.Getenv("SPOTIFY_AUDIO_FEATURES") + "?ids="
	for i, trackID := range trackIDs {
		URL = URL + trackID
		if i+1 != len(trackIDs) {
			URL = URL + ","
		}
	}
	container, err := service.httpClient.Request(
		"GET",
		URL,
		nil,
		"application/json",
		"Bearer "+profile.Credentials.AccessToken,
	)
	if err != nil {
		return nil, err
	}
	defer container.Body.Close()

	features, err := ioutil.ReadAll(container.Body)
	if err != nil {
		return nil, err
	}
	return &features, nil
}
