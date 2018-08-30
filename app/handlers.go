package app

import (
	"net/http"
	"path"
	"encoding/json"

	"github.com/malyusha/image-resizer/app/client"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/malyusha/image-resizer/console"
	_ "github.com/nfnt/resize"
	"github.com/malyusha/image-resizer/app/util"
	"fmt"
	"github.com/nfnt/resize"
	"bytes"
	"image"
	"image/png"
	"image/gif"
	"image/jpeg"
	"strconv"
)

func (a *app) HandleImagesRequest() http.HandlerFunc {
	file := console.Args.PresetsFile
	presets, err := LoadPresets(file)

	if err != nil {
		log.Fatalf("Presets file not found. File - %v", file)
	}

	log.Infof("Successfully parsed presets file `%s`", file)

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			found                 bool
			preset                *Preset
			vars                  = mux.Vars(r)
			presetName, imagePath = vars["preset"], vars["path"]
		)

		// Looking for requested preset
		if preset, found = presets.Find(presetName); !found {
			notFound(w, r)
			return
		}

		log.Infof("Found requested preset: %s", preset.Name)

		storage := a.Storage()
		hash := strconv.Itoa(int(util.Hash(imagePath)))
		hashedFileName := fmt.Sprintf("%s/%s/%s", hash, preset.Name, path.Base(imagePath))

		// If we already have resized image in our storage, then we'll just return it
		if resizedImageBytes, err := storage.Get(hashedFileName); err == nil {
			w.Header().Set("Etag", strconv.Itoa(int(util.Hash(hashedFileName))))
			util.ImageResponse(w, resizedImageBytes)
			return
		}

		sourceImgBytes, err := a.ImageClient().GetImageContent(imagePath)
		if err != nil {
			log.Errorf("Client image error\n%v", err)
			notFound(w, r)
			return
		}

		log.Infof("Found image %s", imagePath)
		//util.ImageResponse(w, sourceImgBytes)

		contentType := http.DetectContentType(sourceImgBytes)
		// Otherwise we need to resize image and store it for future requests
		// First, let's create image instance
		newImageBytes, err := getDecodedImage(sourceImgBytes, contentType)
		if err != nil {
			log.Warn("Failed to decode image %s.", imagePath, contentType)
			// Just return source content
			util.ImageResponse(w, sourceImgBytes)
			return
		}

		img := resize.Resize(preset.Width, preset.Height, newImageBytes, resize.Lanczos3)
		buf := new(bytes.Buffer)
		if err := encodeImageToBytes(img, buf, contentType); err != nil {
			log.Errorf("Failed to encode resized image %s", imagePath)
			util.ImageResponse(w, sourceImgBytes)
			return
		}

		if err := storage.Save(hashedFileName, buf.Bytes()); err != nil {
			log.Errorf("Failed to save resized image `%s` to storage. Error: %v", imagePath, err)
			util.ImageResponse(w, buf.Bytes())
			return
		}

		util.ImageResponse(w, buf.Bytes())
	}
}

func (a *app) HandleHealthCheck() http.HandlerFunc {
	type response struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := &response{
			Status: "OK",
		}

		resp, err := json.Marshal(*response)
		if err != nil {
			http.Error(w, client.InternalErrorMessage, http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)

		log.Infof("Health check requested from %s", r.RemoteAddr)
	}
}

func (a *app) HandleNotFound() http.HandlerFunc {
	return notFound
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Infof("Not found called")
	http.Error(w, "404 - Not Found", http.StatusNotFound)
}

func getDecodedImage(content []byte, contentType string) (image.Image, error) {
	var (
		err error
		img image.Image
		r = bytes.NewReader(content)
	)

	switch contentType {
	case "image/png":
		img, err = png.Decode(r)
		break
	case "image/gif":
		img, err = gif.Decode(r)
		break
	case "image/jpeg":
		img, err = jpeg.Decode(r)
		break
	default:
		return nil, fmt.Errorf("can't resolve decoder for given content type: %s", contentType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to decode image :%v", err)
	}

	return img, nil
}

func encodeImageToBytes(img image.Image, buf *bytes.Buffer, contentType string) error {
	var err error

	switch contentType {
	case "image/png":
		err = png.Encode(buf, img)
	case "image/gif":
		err = gif.Encode(buf, img, nil)
	default:
		err = jpeg.Encode(buf, img, nil)
	}

	return err
}


