package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"bytestream/internal/api"
	"bytestream/internal/cache"
	"bytestream/internal/clients"
	"bytestream/internal/config"
)

func main() {
	addr := config.Env("HTTP_ADDR", ":8080")
	redisAddr := config.Env("REDIS_ADDR", "localhost:6379")

	identityServiceBaseUrl := config.Env("IDENTITY_BASE_URL", "http://localhost:18081")
	availabilityServiceBaseUrl := config.Env("AVAILABILITY_BASE_URL", "http://localhost:18082")

	timeout := config.DurationFromEnv("UPSTREAM_TIMEOUT", 3*time.Second)
	identityCacheExpiry := config.DurationFromEnv("IDENTITY_CACHE_EXPIRY", 2*time.Minute)
	availabilityCACHEexpiry := config.DurationFromEnv("AVAILABILITY_CACHE_EXPIRY", 10*time.Minute)

	c := cache.NewRedisCache(redisAddr)
	if err := c.Ping(context.Background()); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	identityClient := clients.NewIdentityClient(identityServiceBaseUrl, timeout)
	availabilityClient := clients.NewAvailabilityClient(availabilityServiceBaseUrl, timeout)

	h := api.NewHandler(api.Deps{
		Cache:           c,
		Identity:        identityClient,
		Availability:    availabilityClient,
		IdentityTTL:     identityCacheExpiry,
		AvailabilityTTL: availabilityCACHEexpiry,
	})

	router := api.NewRouter(h)

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
