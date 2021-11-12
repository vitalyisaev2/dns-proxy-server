package v1

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-resty/resty"
    "github.com/mageddo/dns-proxy-server/cache/store"
    "github.com/stretchr/testify/assert"
)

func TestGetCache(t *testing.T) {

    s := httptest.NewServer(nil)
    defer s.Close()

    cache := store.GetInstance()
    cache.Put("ke1", "value1")

    r, err := resty.R().
        Get(s.URL + CACHE_V1)
    assert.Nil(t, err)
    assert.Equal(t, http.StatusOK, r.StatusCode())
    assert.Equal(t, `{"TEST_MODE":true,"ke1":"value1"}`, r.String())

}

func TestGetCacheSize(t *testing.T) {

    s := httptest.NewServer(nil)
    defer s.Close()

    cache := store.GetInstance()
    cache.Put("ke1", "value1")
    cache.Put("ke2", "value2")

    r, err := resty.R().
        Get(s.URL + CACHE_SIZE_V1)
    assert.Nil(t, err)
    assert.Equal(t, http.StatusOK, r.StatusCode())
    assert.Equal(t, `{"size":3}`, r.String())

}
