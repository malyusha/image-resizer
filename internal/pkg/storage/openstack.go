package storage

import (
	"time"

	"github.com/ulule/gostorages"
)

type OpenStackStorage struct {

}

func NewOpenStackStorage(config *OpenStackStorageConfig) (*OpenStackStorage, error) {
	panic("not implemented")
}

func (s *OpenStackStorage) Save(filepath string, file gostorages.File) error {
	panic("implement me")
}

func (s *OpenStackStorage) Path(filepath string) string {
	panic("implement me")
}

func (s *OpenStackStorage) Exists(filepath string) bool {
	panic("implement me")
}

func (s *OpenStackStorage) Delete(filepath string) error {
	panic("implement me")
}

func (s *OpenStackStorage) Open(filepath string) (gostorages.File, error) {
	panic("implement me")
}

func (s *OpenStackStorage) ModifiedTime(filepath string) (time.Time, error) {
	panic("implement me")
}

func (s *OpenStackStorage) Size(filepath string) int64 {
	panic("implement me")
}

func (s *OpenStackStorage) URL(filename string) string {
	panic("implement me")
}

func (s *OpenStackStorage) HasBaseURL() bool {
	panic("implement me")
}

func (s *OpenStackStorage) IsNotExist(err error) bool {
	panic("implement me")
}
