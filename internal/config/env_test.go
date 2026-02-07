package config

import (
	"os"
	"testing"
	"time"
)

func TestEnv(t *testing.T) {
	const key = "TEST_ENV_KEY"

	t.Run("returns fallback when missing", func(t *testing.T) {
		t.Setenv(key, "")
		got := Env(key, "fallback")
		if got != "fallback" {
			t.Fatalf("expected fallback, got %q", got)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv(key, "value")
		got := Env(key, "fallback")
		if got != "value" {
			t.Fatalf("expected value, got %q", got)
		}
	})
}

func TestDurationFromEnv(t *testing.T) {
	const key = "TEST_DURATION_KEY"

	t.Run("returns fallback when missing", func(t *testing.T) {
		os.Unsetenv(key)
		got := DurationFromEnv(key, 3*time.Second)
		if got != 3*time.Second {
			t.Fatalf("expected 3s, got %v", got)
		}
	})

	t.Run("parses duration when valid", func(t *testing.T) {
		t.Setenv(key, "2m")
		got := DurationFromEnv(key, 3*time.Second)
		if got != 2*time.Minute {
			t.Fatalf("expected 2m, got %v", got)
		}
	})

	t.Run("returns fallback when invalid", func(t *testing.T) {
		t.Setenv(key, "not-a-duration")
		got := DurationFromEnv(key, 3*time.Second)
		if got != 3*time.Second {
			t.Fatalf("expected fallback 3s, got %v", got)
		}
	})
}

func TestBoolFromEnv(t *testing.T) {
	const key = "TEST_BOOL_KEY"

	t.Run("returns fallback when missing", func(t *testing.T) {
		os.Unsetenv(key)
		got := BoolFromEnv(key, false)
		if got != false {
			t.Fatalf("expected false, got %v", got)
		}
	})

	t.Run("parses true", func(t *testing.T) {
		t.Setenv(key, "true")
		got := BoolFromEnv(key, false)
		if got != true {
			t.Fatalf("expected true, got %v", got)
		}
	})

	t.Run("parses 1 as true", func(t *testing.T) {
		t.Setenv(key, "1")
		got := BoolFromEnv(key, false)
		if got != true {
			t.Fatalf("expected true, got %v", got)
		}
	})

	t.Run("returns fallback when invalid", func(t *testing.T) {
		t.Setenv(key, "maybe")
		got := BoolFromEnv(key, true)
		if got != true {
			t.Fatalf("expected fallback true, got %v", got)
		}
	})
}
