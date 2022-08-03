package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func TopQueryValidator(e string, paramType string) (string, error) {
	topType := [...]string{"tracks", "artists"}
	timeRange := [...]string{"short_term", "medium_term", "long_term"}
	switch paramType {
	case "type":
		for _, v := range topType {
			if v == e {
				return e, nil
			}
		}
		return "tracks", nil
	case "time_range":
		for _, v := range timeRange {
			if v == e {
				return e, nil
			}
		}
		return "medium_term", nil
	}
	return "", errors.New("incorrect format parameter")
}

func (service *Service) GetRecentlyPlayed(
	email string, limit int,
	before string, after string,
) (*[]byte, error) {
	credentials, err := service.GetValidToken(email)
	if err != nil {
		return nil, err
	}
	URL := "https://api.spotify.com/v1/me/player/recently-played"
	if limit != 0 {
		URL = URL + "?limit=" + strconv.Itoa(limit)
	}
	if before != "-6795364578871" {
		URL = URL + "&before=" + before
	}
	if after != "-6795364578871" {
		URL = URL + "&after=" + after
	}
	recentlyPlayedResp, err := service.httpClient.Request(
		"GET",
		URL,
		nil,
		"application/json",
		"Bearer "+credentials.AccessToken,
	)
	if err != nil {
		return nil, err
	}
	defer recentlyPlayedResp.Body.Close()

	recentlyPlayed, err := ioutil.ReadAll(recentlyPlayedResp.Body)
	if err != nil {
		return nil, err
	}
	return &recentlyPlayed, nil
}

func (service *Service) GetTopArtistsOrTracks(email string,
	top string,
	timeRangeStr string,
	limit int,
	offset int) (*[]byte, error) {
	credentials, err := service.GetValidToken(email)
	if err != nil {
		return nil, err
	}
	toptype, err := TopQueryValidator(top, "type")
	if err != nil {
		return nil, err
	}
	timeRange, err := TopQueryValidator(timeRangeStr, "time_range")
	if err != nil {
		return nil, err
	}
	URL := os.Getenv("SPOTIFY_PERSONAL_TOP") + "/" + toptype +
		"?limit=" + strconv.Itoa(limit) +
		"&offset=" + strconv.Itoa(offset) +
		"&time_range=" + timeRange
	topResp, err := service.httpClient.Request(
		"GET",
		URL,
		nil,
		"application/json",
		"Bearer "+credentials.AccessToken,
	)
	if err != nil {
		return nil, err
	}
	defer topResp.Body.Close()

	topContainer, err := ioutil.ReadAll(topResp.Body)
	if err != nil {
		return nil, err
	}
	return &topContainer, nil
}

// get user's palylists
func (service *Service) GetUserPlaylists(email string, limit int, offset int) (*[]byte, error) {
	credentials, err := service.GetValidToken(email)
	if err != nil {
		return nil, err
	}
	URL := os.Getenv("SPOTIFY_PERSONAL_PLAYLISTS") +
		"?limit=" + strconv.Itoa(limit) +
		"&offset=" + strconv.Itoa(offset)
	fmt.Println(URL)
	playlistsResp, err := service.httpClient.Request(
		"GET",
		URL,
		nil,
		"application/json",
		"Bearer "+credentials.AccessToken,
	)
	if err != nil {
		return nil, err
	}
	defer playlistsResp.Body.Close()

	playlists, err := ioutil.ReadAll(playlistsResp.Body)
	if err != nil {
		return nil, err
	}
	return &playlists, nil
}

// get Top Tracks
func (service *Service) GetPersonalAudioFeatures(email string, timespan string) (*[]byte, error) {
	tracksByteArray, err := service.GetTopArtistsOrTracks(email, "tracks", timespan, 50, 0)
	if err != nil {
		return nil, err
	}

	// unmarshall tracksByteArray to array of Track
	type TopItemsResp struct {
		Items []Track `json:"items"`
	}
	var topItemsResp TopItemsResp
	err = json.Unmarshal(*tracksByteArray, &topItemsResp)
	if (err) != nil {
		return nil, err
	}
	// loop through tracks and concat id separated by ","
	var trackIds []string
	for _, track := range topItemsResp.Items {
		trackIds = append(trackIds, track.ID)
	}

	audioFeaturesByteArray, err := service.GetTracksAudioFeatures(email, trackIds)
	if err != nil {
		return nil, err
	}

	type AudioFeaturesResp struct {
		AudioFeatures []AudioFeatures `json:"audio_features"`
	}
	// unmarshall audioFeaturesByteArray to array of AudioFeatures
	var audioFeaturesResep AudioFeaturesResp
	err = json.Unmarshal(*audioFeaturesByteArray, &audioFeaturesResep)
	if (err) != nil {
		return nil, err
	}
	// loop through audioFeatures and sum integer and float and float fields and reduce to one audio feature
	var audioFeaturesSum AudioFeatures
	for _, audioFeature := range audioFeaturesResep.AudioFeatures {
		audioFeaturesSum.Danceability += audioFeature.Danceability
		audioFeaturesSum.Energy += audioFeature.Energy
		audioFeaturesSum.Loudness += audioFeature.Loudness
		audioFeaturesSum.Speechiness += audioFeature.Speechiness
		audioFeaturesSum.Acousticness += audioFeature.Acousticness
		audioFeaturesSum.Instrumentalness += audioFeature.Instrumentalness
		audioFeaturesSum.Liveness += audioFeature.Liveness
		audioFeaturesSum.Valence += audioFeature.Valence
		audioFeaturesSum.Tempo += audioFeature.Tempo
		audioFeaturesSum.TimeSignature += audioFeature.TimeSignature
	}
	// divide each attr of audioFeaturesSum by the length of audioFeaturesResep.AudioFeatures
	audioFeaturesSum.Danceability /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Energy /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Loudness /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Speechiness /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Acousticness /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Instrumentalness /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Liveness /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Valence /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.Tempo /= float64(len(audioFeaturesResep.AudioFeatures))
	audioFeaturesSum.TimeSignature /= int(len(audioFeaturesResep.AudioFeatures))

	// marshall audioFeaturesSum to byte array
	audioFeaturesSumByteArray, err := json.Marshal(audioFeaturesSum)
	return &audioFeaturesSumByteArray, err
}
