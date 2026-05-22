package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheGetSet(t *testing.T) {
	cacheTTL = 100 * time.Millisecond

	cacheSet("test-key", "test-body", http.StatusOK, http.Header{"Content-Type": {"text/html"}})
	e := cacheGet("test-key")
	assert.NotNil(t, e)
	assert.Equal(t, "test-body", e.body)
	assert.Equal(t, http.StatusOK, e.status)

	time.Sleep(150 * time.Millisecond)
	e = cacheGet("test-key")
	assert.Nil(t, e)
}

func TestCacheDisabled(t *testing.T) {
	cacheTTL = 0

	cacheSet("key", "body", http.StatusOK, nil)
	e := cacheGet("key")
	assert.Nil(t, e)
}
