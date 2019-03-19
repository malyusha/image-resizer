package client

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPImageClient(t *testing.T) {
	instance, err := NewHTTPImageClient(newConfig())
	assert.Implements(t, (*Client)(nil), instance)
	assert.NoError(t, err)
}

func TestHTTPImageClient_FullPath(t *testing.T) {
	host := "http://lorempixel.com"
	instance, _ := NewHTTPImageClient(newConfig(map[string]interface{}{"host": host}))

	assert.Equal(t, host + "/400/200", instance.FullPath("400/200"))
	assert.Equal(t, host + "/400/200", instance.FullPath("/400/200"))
	assert.Equal(t, host + "/400/200", instance.FullPath("////400/200"))
}

func TestHTTPImageClient_GetImageContent(t *testing.T) {
	server := setupHttpServer()
	instance, _ := NewHTTPImageClient(newConfig(map[string]interface{}{"host": server.URL}))
	defer server.Close()

	content, err := instance.GetImageContent("any.png")
	assert.NoError(t, err)
	assert.Equal(t, "image/png", http.DetectContentType(content))
	assert.True(t, len(content) > 0)
}

func setupHttpServer() *httptest.Server {
	testImage := image.NewRGBA(image.Rect(0, 0, 100, 100))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buff bytes.Buffer
		if err := png.Encode(&buff, testImage); err != nil {
			http.Error(w, "Failed to create image", http.StatusInternalServerError)
			return
		}

		w.Write(buff.Bytes())

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(buff.Len()))
		w.Header().Set("Cache-Control", "public, max-age=86400")
	})
	return httptest.NewServer(handler)
}
