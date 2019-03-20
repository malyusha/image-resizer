package storage

import (
	"errors"
	"fmt"

	"github.com/mitchellh/goamz/aws"
	"github.com/ulule/gostorages"
)

const (
	s3StorageType        = "s3"
	fsStorageType        = "fs"
	openStackStorageType = "open_stack"
)

// New returns new instance of storage
func New(config *Config) (gostorages.Storage, error) {
	if config == nil {
		return new(DummyStorage), nil
	}

	return newStorage(config)
}

func newStorage(config *Config) (gostorages.Storage, error) {
	switch config.Type {
	case fsStorageType:
		if config.FS == nil {
			return nil, errors.New("fs storage config not provided")
		}
		return gostorages.NewFileSystemStorage(config.FS.Location, ""), nil
	case s3StorageType:
		s3Config := config.S3
		if s3Config == nil {
			return nil, errors.New("s3 storage config not provided")
		}

		acl, ok := gostorages.ACLs[s3Config.ACL]
		if !ok {
			return nil, fmt.Errorf("ACL 5s is not allowed", s3Config.ACL)
		}

		region, ok := aws.Regions[s3Config.Region]
		if !ok {
			return nil, fmt.Errorf("Region %s is not allowed", region)
		}

		return gostorages.NewS3Storage(
			s3Config.AccessKeyID,
			s3Config.SecretAccessKey,
			s3Config.BucketName,
			"",
			region,
			acl,
			s3Config.BaseURL,
		), nil

	case openStackStorageType:
		return NewOpenStackStorage(config.OpenStack)
	}

	return nil, fmt.Errorf("storage %s does not exist", config.Type)
}
