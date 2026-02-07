# Bytestream Video Service

This service exposes a single endpoint that returns video playback
information based on:

1.  The user identity (premium or standard)
2.  The video availability window

It integrates with: - An Identity Service (mocked with WireMock) - An
Availability Service (mocked with WireMock) - Redis (for caching
identity and availability responses)

------------------------------------------------------------------------

## Endpoint

GET /videos/{videoID}

Example:

curl -sS -i -H "Authorization: Bearer anytoken"
http://localhost:8080/videos/46325

The service:

1.  Validates the Authorization: Bearer `<token>`{=html} header
2.  Retrieves user identity (from cache or Identity service)
3.  Retrieves availability window (from cache or Availability service)
4.  Checks if the video is currently within its availability window
5.  Returns playback details:
    -   Standard quality for regular users
    -   Premium quality for premium users

------------------------------------------------------------------------

## Architecture Overview

Redis: - Caches Identity responses (keyed by bearer token) - Caches
Availability responses (keyed by videoID) - Reduces repeated upstream
calls

WireMock: - Simulates Identity Service - Simulates Availability
Service - Allows easy testing of different scenarios (premium,
unauthorized, unavailable, etc.)

Video Metadata: - Hardcoded in internal/domain/video.go - This is
intentional per assessment requirements (no video datastore provided)

------------------------------------------------------------------------

## Running Locally

Start everything:

docker compose up --build

Services started: - App → http://localhost:8080 - Identity Mock →
http://localhost:18081 - Availability Mock → http://localhost:18082 -
Redis → localhost:6379

------------------------------------------------------------------------

## Example Requests

Premium (default mock behavior):

curl -sS -i -H "Authorization: Bearer anytoken"
http://localhost:8080/videos/46325

By default, the Identity mock returns a premium role.

------------------------------------------------------------------------

Make Video Unavailable:

Edit: docker/wiremock/availability/mappings/availability.json

Change the date range so the current date is outside the window.

Then recreate the availability mock container:

docker compose restart availability-mock

If changes are not reflected:

docker compose up -d --force-recreate availability-mock

------------------------------------------------------------------------

## Important: Updating Mock Files

WireMock loads mappings at container startup.

If you modify any file under:

docker/wiremock/\*\*/mappings/

You must restart the corresponding mock container:

docker compose restart identity-mock docker compose restart
availability-mock

The application container does not need to be restarted when mocks
change.

------------------------------------------------------------------------

## Caching Behavior

-   Identity is cached using: identity:`<bearer-token>`{=html}

-   Availability is cached using: availability:`<videoID>`{=html}

Cache TTL is configurable via environment variables.

You can disable cache completely:

SKIP_CACHE: "true"

When enabled: - The service bypasses Redis - Calls upstream mocks
directly

------------------------------------------------------------------------

## Error Handling

Upstream 401 / 403 (Identity) → 401 Unauthorized\
Upstream 404 (Availability) → 404 Not Found\
Other upstream errors → 502 Bad Gateway\
Outside availability window → 404 Not Available

------------------------------------------------------------------------

## Design Notes

-   Clean separation of HTTP layer, clients, cache, and domain logic
-   Context-aware upstream calls
-   Graceful cache failure handling (cache errors do not fail request)
-   Explicit upstream error mapping
-   Structured project layout suitable for production-style services
