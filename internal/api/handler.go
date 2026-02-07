package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"bytestream/internal/cache"
	"bytestream/internal/clients"
	"bytestream/internal/config"
	"bytestream/internal/domain"

	"github.com/go-chi/chi/v5"
)

type Deps struct {
	Cache           *cache.RedisCache
	Identity        *clients.IdentityClient
	Availability    *clients.AvailabilityClient
	IdentityTTL     time.Duration
	AvailabilityTTL time.Duration
}

type Handler struct {
	cache           *cache.RedisCache
	identity        *clients.IdentityClient
	availability    *clients.AvailabilityClient
	identityTTL     time.Duration
	availabilityTTL time.Duration
}

func NewHandler(d Deps) *Handler {
	return &Handler{
		cache:           d.Cache,
		identity:        d.Identity,
		availability:    d.Availability,
		identityTTL:     d.IdentityTTL,
		availabilityTTL: d.AvailabilityTTL,
	}
}

func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	videoIDStr := chi.URLParam(r, "videoID")
	videoID, err := strconv.Atoi(videoIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errBody("bad_request", "videoID must be an integer"))
		return
	}

	meta, ok := domain.LookupVideo(videoID)
	if !ok {
		writeJSON(w, http.StatusNotFound, errBody("not_found", "unknown video_id"))
		return
	}

	token, _ := extractBearer(r)
	skipCache := config.BoolFromEnv("SKIP_CACHE", false)

	identityKey := "identity:" + token
	var identity domain.Identity
	foundIdentity := false

	if !skipCache {
		found, cacheErr := h.cache.GetJSON(r.Context(), identityKey, &identity)
		if cacheErr != nil {
			log.Printf("cache read failed (identity): %v", cacheErr)
			found = false
		}
		foundIdentity = found
	}

	if !foundIdentity {
		identity, err = h.identity.GetUserInfo(token)
		if err != nil {
			if upErr, isUpstreamErr := err.(clients.UpstreamError); isUpstreamErr &&
				(upErr.Status == http.StatusUnauthorized || upErr.Status == http.StatusForbidden) {
				writeJSON(w, http.StatusUnauthorized, errBody("unauthorized", "token rejected by identity service"))
				return
			}
			writeJSON(w, http.StatusBadGateway, errBody("bad_gateway", "identity service failed"))
			return
		}

		if !skipCache {
			if err := h.cache.SetJSON(r.Context(), identityKey, identity, h.identityTTL); err != nil {
				log.Printf("cache write failed (identity): %v", err)
			}
		}
	}

	isPremium := hasRole(identity.Roles, "premium")

	availKey := "availability:" + strconv.Itoa(videoID)
	var avail domain.AvailabilityInfo
	foundAvailability := false

	if !skipCache {
		found, cacheErr := h.cache.GetJSON(r.Context(), availKey, &avail)
		if cacheErr != nil {
			log.Printf("cache read failed (availability): %v", cacheErr)
			found = false
		}
		foundAvailability = found
	}

	if !foundAvailability {
		avail, err = h.availability.GetAvailability(token, videoID)
		if err != nil {
			if upErr, isUpstreamErr := err.(clients.UpstreamError); isUpstreamErr && upErr.Status == http.StatusNotFound {
				writeJSON(w, http.StatusNotFound, errBody("not_found", "video_id not found in availability service"))
				return
			}
			writeJSON(w, http.StatusBadGateway, errBody("bad_gateway", "availability service failed"))
			return
		}

		if !skipCache {
			if err := h.cache.SetJSON(r.Context(), availKey, avail, h.availabilityTTL); err != nil {
				log.Printf("cache write failed (availability): %v", err)
			}
		}
	}

	// Availibilty Window check
	if !isInWindow(time.Now(), avail.AvailabilityWindow.From, avail.AvailabilityWindow.To) {
		writeJSON(w, http.StatusForbidden, errBody("not_available", "video is not available in the current window"))
		return
	}

	resp := domain.BuildVideoResponse(videoID, meta, isPremium)
	writeJSON(w, http.StatusOK, resp)
}

func hasRole(roles []string, want string) bool {
	for _, r := range roles {
		if r == want {
			return true
		}
	}
	return false
}

func isInWindow(now time.Time, fromYYYYMMDD, toYYYYMMDD string) bool {
	const format = "2006-01-02"
	from, err1 := time.Parse(format, fromYYYYMMDD)
	to, err2 := time.Parse(format, toYYYYMMDD)
	if err1 != nil || err2 != nil {
		return false
	}
	toEnd := to.Add(24*time.Hour - time.Second)
	n := now.UTC()
	return !n.Before(from) && !n.After(toEnd)
}

func errBody(code, msg string) map[string]string {
	return map[string]string{"error": code, "message": msg}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
