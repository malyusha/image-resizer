package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"

	"github.com/malyusha/image-resizer/internal/pkg/server/errors"
	"github.com/malyusha/image-resizer/pkg/client"
	"github.com/malyusha/image-resizer/pkg/preset"
	"github.com/malyusha/image-resizer/pkg/storage"
	"github.com/malyusha/image-resizer/pkg/util"
)

func (s *Instance) HandleImagesRequest(st storage.Storage, cl client.Client) http.HandlerFunc {
	logger := s.app.Logger()
	// TODO: move presets file parsing to startup and enable watcher for changes
	file := s.app.Config().Get("presets_file").String()
	if file == "" {
		logger.Fatal("Presets file is required")
	}

	presets, err := preset.LoadPresets(file)

	if err != nil {
		logger.Fatalf("Error reading presets file. File - %s. %s", file, err)
	}

	logger.Infof("Successfully parsed presets file `%s`", file)

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			found                 bool
			preset                *preset.Preset
			vars                  = mux.Vars(r)
			presetName, imagePath = vars["preset"], vars["path"]
		)

		// Looking for requested preset
		if preset, found = presets.Find(presetName); !found {
			notFound(w, r)
			return
		}

		logger.WithField("preset", preset).Info("Found requested preset")

		hash := strconv.Itoa(int(util.Hash(imagePath)))
		hashedFileName := fmt.Sprintf("%s/%s/%s", hash, preset.Name, path.Base(imagePath))

		// If we already have resized image in our storage, then we'll just return it
		if resizedImageBytes, err := st.Get(hashedFileName); err == nil {
			w.Header().Set("Etag", strconv.Itoa(int(util.Hash(hashedFileName))))
			s.RespondWithFile(w, resizedImageBytes)
			return
		}

		sourceImgBytes, err := cl.GetImageContent(imagePath)
		if err != nil {
			logger.WithError(err).Error("Client image error")
			notFound(w, r)
			return
		}

		fullPath := cl.FullPath(imagePath)
		logger.Infof("Found file %s", fullPath)
		//util.ImageResponse(w, sourceImgBytes)

		contentType := http.DetectContentType(sourceImgBytes)
		// Otherwise we need to resize image and store it for future requests
		// First, let's create image instance
		newImageBytes, err := getDecodedImage(sourceImgBytes, contentType)
		if err != nil {
			logger.WithField("file", fullPath).Warnf("Failed to decode image %s.", contentType)
			// Just return source content
			s.RespondWithFile(w, sourceImgBytes)
			return
		}

		img := resize.Resize(preset.Width, preset.Height, newImageBytes, resize.Lanczos3)
		buf := new(bytes.Buffer)
		if err := encodeImageToBytes(img, buf, contentType); err != nil {
			logger.Errorf("Failed to encode resized image %s", imagePath)
			s.RespondWithFile(w, sourceImgBytes)
			return
		}

		savedFile, err := st.Save(hashedFileName, buf.Bytes())
		if err != nil {
			logger.Errorf("Failed to save resized image `%s` to storage. Error: %v", imagePath, err)
			s.RespondWithFile(w, buf.Bytes())
			return
		}

		logger.Infof(`Successfully saved resized file "%s"`, savedFile)

		s.RespondWithFile(w, buf.Bytes())
	}
}

func (s *Instance) HandleHealthCheck() http.HandlerFunc {
	type response struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := &response{
			Status: "OK",
		}

		resp, err := json.Marshal(*response)
		if err != nil {
			http.Error(w, errors.InternalErrorMessage, http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// HandleNotFound serves not found error to client
func (s *Instance) HandleNotFound() http.HandlerFunc {
	return notFound
}

// RespondWithFile writes image meta info into headers
func (s *Instance) RespondWithFile(w http.ResponseWriter, img []byte) {
	contentType := http.DetectContentType(img)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))

	if _, err := w.Write(img); err != nil {
		s.app.Logger().Errorf("Unable to write file: %s", err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 - Not Found", http.StatusNotFound)
}

func getDecodedImage(content []byte, contentType string) (image.Image, error) {
	var (
		err error
		img image.Image
		r   = bytes.NewReader(content)
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
