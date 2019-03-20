package client

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/malyusha/image-resizer/pkg/dot"
)

type LocalStorageClientTestSuite struct {
	suite.Suite
	images []string
	dir string
}

func TestLocalStorageClientTestSuite(t *testing.T) {
	suite.Run(t, new(LocalStorageClientTestSuite))
}

func (s *LocalStorageClientTestSuite) SetupTest() {
	s.dir = ".testdata/local-storage"
	s.images = append(s.images, "one.png", "tow.png", "three.png")
	os.MkdirAll(s.dir, os.ModePerm)

	for _, name := range s.images {
		img := image.NewRGBA(image.Rect(0,0, 100, 50))
		img.Set(1, 2, color.RGBA{R: 255})

		f, err := os.OpenFile(filepath.Join(s.dir, name), os.O_WRONLY | os.O_CREATE, 0600)

		if err != nil {
			s.FailNow(fmt.Sprintf("failed to create test images: %s", err))
		}
		defer f.Close()

		if err := png.Encode(f, img); err != nil {
			s.FailNow(fmt.Sprintf("failed to create test images: %s", err))
		}
	}
}

func (s *LocalStorageClientTestSuite) TearDownTest() {
	if err := os.RemoveAll(".testdata"); err != nil {
		s.FailNow("failed to remove tmp directory with images")
	}
}

func (s *LocalStorageClientTestSuite) TestNewInstance() {
	emptyConfig := newConfig()

	_, err := NewLocalStorageClient(emptyConfig)
	assert.Error(s.T(), err)

	cl, err := NewLocalStorageClient(s.correctConfig())

	s.NoError(err, "No error should be returned with correct configuration")
	s.Implements((*Client)(nil), cl)
}

func (s *LocalStorageClientTestSuite) TestCorrectFullPath() {
	instance, _ := NewLocalStorageClient(s.correctConfig())

	s.Equal(instance.FullPath("file"), filepath.Join(s.dir, "file"))
}

func (s *LocalStorageClientTestSuite) TestReturnsImageContent() {
	instance, _ := NewLocalStorageClient(s.correctConfig())

	for _, name := range s.images {
		b, err := instance.GetImageContent(name)

		s.True(len(b) > 0, "Returned image content is empty")
		s.NoError(err, "No error should be returned from `GetImageContent` in `FileSystemClient`")
	}
}

func (s *LocalStorageClientTestSuite) correctConfig() *dot.Map {
	return newConfig(map[string]interface{}{"dir": s.dir})
}

func newConfig(values ...map[string]interface{}) *dot.Map {
	var v map[string]interface{}
	if len(values) == 0 {
		return dot.NewMap(v)
	}

	return dot.NewMap(values[0])
}
