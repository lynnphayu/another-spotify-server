package endpoint

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"utilserver/pkg/spotify"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid/v4"
)

func clearCookie(w *http.ResponseWriter) {
	c := &http.Cookie{
		Name:    "storage",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	http.SetCookie(*w, c)
}

func setCookie(w *http.ResponseWriter, name string, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(*w, c)
}

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
}
type Handler struct {
	cache    spotify.Cache
	services spotify.Services
}

// Handler - spotify authentication routes handler
func NewHandler(cache spotify.Cache, services spotify.Services) http.Handler {
	handler := new(Handler)
	handler.cache = cache
	handler.services = services
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.Handle("/spotify/login", handler.login()).Methods(http.MethodGet)
	api.Handle("/spotify/callback", handler.loginCallback()).Methods(http.MethodGet)
	api.Handle("/spotify/profile", attachMiddleware(handler.getProfile(), handler.authMiddleware)).Methods(http.MethodGet)
	api.Handle("/spotify/recently_played", attachMiddleware(handler.getRecentlyPlayed(), handler.authMiddleware)).Methods(http.MethodGet)
	// api.Handle("/spotify/audio_features", attachMiddleware(handler.getAudioFeatures(), handler.authMiddleware)).Methods(http.MethodGet)
	api.Handle("/spotify/top", attachMiddleware(handler.getTops(), handler.authMiddleware)).Methods(http.MethodGet)
	api.Handle("/spotify/playlists", attachMiddleware(handler.getPlaylists(), handler.authMiddleware)).Methods(http.MethodGet)
	api.Handle("/spotify/audio_features", attachMiddleware(handler.getPersonalAudioFeatures(), handler.authMiddleware)).Methods(http.MethodGet)
	return r
}

func attachMiddleware(h http.Handler, middlewares ...mux.MiddlewareFunc) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func (handler Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		claim, err := verifyToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if claim.Email == "" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		r.Header.Set("email", claim.Email)
		next.ServeHTTP(w, r)
	})
}

// get spotify login url from environment variables, parse url and redirect to that url
func (handler Handler) redirectToSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	parm := url.Values{}
	base, err := url.Parse(os.Getenv("SPOTIFY_LOGIN_ENDPOINT"))
	// get redirect from query params and set it to cache if it isn't empty
	redirect := r.URL.Query().Get("redirect")

	id := shortuuid.New()

	parm.Add("state", id)
	// set state to cookie
	setCookie(&w, os.Getenv("SPOTIFY_LOGIN_STATE_KEY"), id)
	if redirect != "" {
		handler.cache.Set(id, redirect, 0)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	parm.Add("client_id", os.Getenv("CLIENT_ID"))
	parm.Add("scope", os.Getenv("SCOPES"))
	parm.Add("response_type", "code")
	parm.Add("redirect_uri", os.Getenv("REDIRECT_URL"))
	base.RawQuery = parm.Encode()
	// enable cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// change base to string and print
	http.Redirect(w, r, base.String(), http.StatusTemporaryRedirect)
}

// verify jwt token and extract email from token
func verifyToken(tokenString string) (*spotify.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &spotify.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*spotify.CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

// login
func (handler Handler) login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.redirectToSpotifyLogin(w, r)
	})
}

func (handler *Handler) getProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("email")
		profile, profileErr := handler.services.Auth.Login(email)
		if profile == nil {
			http.Error(w, "no profile", http.StatusInternalServerError)
			return
		}
		if profileErr != nil {
			http.Error(w, profileErr.Error(), http.StatusInternalServerError)
			return
		}
		profileByteArr, marshallingErr := json.Marshal(profile)
		if marshallingErr != nil {
			http.Error(w, marshallingErr.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(profileByteArr)
	})
}

func (handler *Handler) loginCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		storedStateCookie, _ := r.Cookie(os.Getenv("SPOTIFY_LOGIN_STATE_KEY"))

		if state == "" || (storedStateCookie != nil && (state != storedStateCookie.Value)) {
			http.Error(w, "Invalid state", http.StatusForbidden)
		} else {
			clearCookie(&w)

			loginReponse, err := handler.services.Auth.AuthCallback(code)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// get redirect from cache with key from state
			redirect, _ := handler.cache.Get(state)
			if redirect == "" || redirect == nil {
				http.Redirect(w, r, "/api/v1/spotify/profile?token="+loginReponse.Token, http.StatusTemporaryRedirect)
			} else {
				redirect, ok := redirect.(string)
				if !ok {
					http.Error(w, "Invalid redirect", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, redirect+"?token="+loginReponse.Token, http.StatusTemporaryRedirect)
				err := handler.cache.Clear(state)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			return
		}
	})
}

func (handler *Handler) getRecentlyPlayed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		email := r.Header.Get("email")
		var query RecentlyPlayedQurey = RecentlyPlayedQurey{
			Email:  email,
			Before: r.URL.Query().Get("before"),
			After:  r.URL.Query().Get("after"),
		}

		if limit := r.URL.Query().Get("limit"); limit != "" {
			if i, err := strconv.Atoi(limit); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if err == nil {
				query.Limit = i
			}
		} else {
			query.Limit = 20
		}

		validate := validator.New()
		if errors := validate.Struct(query); errors != nil {
			http.Error(w, errors.Error(), http.StatusBadRequest)
			return
		}
		timeBefore, _ := time.Parse("2006-01-02", query.Before)
		timeAfter, _ := time.Parse("2006-01-02", query.After)

		recentlyPlayed, err := handler.services.PersonalInfo.GetRecentlyPlayed(
			query.Email, query.Limit,
			strconv.FormatInt(timeBefore.UnixNano()/1000000, 10),
			strconv.FormatInt(timeAfter.UnixNano()/1000000, 10),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(*recentlyPlayed)
	})
}

func (handler *Handler) getAudioFeatures() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("email")
		trackIDs := r.URL.Query().Get("ids")
		trackIDsArray := strings.Split(trackIDs, ",")
		resp, err := handler.services.General.GetTracksAudioFeatures(email, trackIDsArray)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(*resp)
	})
}

func (handler *Handler) getTops() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("email")
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}
		timeRange := r.URL.Query().Get("time_range")
		topType := r.URL.Query().Get("type")

		resp, err := handler.services.PersonalInfo.GetTopArtistsOrTracks(email, topType, timeRange, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(*resp)
	})
}

// get personal audio features
func (handler *Handler) getPersonalAudioFeatures() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("email")
		// get timespan from query
		timeSpan := r.URL.Query().Get("timespan")
		resp, err := handler.services.PersonalInfo.GetPersonalAudioFeatures(email, timeSpan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(*resp)
	})
}

// get playlists handler
func (handler *Handler) getPlaylists() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("email")
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}

		resp, err := handler.services.PersonalInfo.GetUserPlaylists(email, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(*resp)
	})
}
