package util

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func ImageResponse(w http.ResponseWriter, img []byte) {

	contentType := http.DetectContentType(img)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))

	if _, err := w.Write(img); err != nil {
		log.Errorf("unable to write image: %v", err)
	}
}