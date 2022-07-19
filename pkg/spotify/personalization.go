package spotify

import (
	"errors"
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
