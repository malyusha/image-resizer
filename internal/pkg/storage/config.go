package storage

// S3StorageConfig represents struct for s3 storage
type S3StorageConfig struct {
	BaseURL         string `mapstructure:"base_url"`
	Region          string
	ACL             string
	AccessKeyID     string `mapstructure:"access_key_id"`
	BucketName      string `mapstructure:"bucket_name"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

// FSStorageConfig represents struct for file system storage
type FSStorageConfig struct {
	Location string
}

// OpenStackConfig represents struct for open stack storage
type OpenStackStorageConfig struct {
	// TODO write configuration structure
}

// Config represents configuration struct for storage
type Config struct {
	Type      string
	S3        *S3StorageConfig
	FS        *FSStorageConfig
	OpenStack *OpenStackStorageConfig `mapstructure:"open_stack"`
}
