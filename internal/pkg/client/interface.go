package client

// Client is the client retrieving content of source image file.
type Client interface {
	// GetImageContent returns data of given file. If file doesn't exist, or other error occurred
	// error must be returned.
	GetImageContent(path string) ([]byte, error)

	// Path path returns string, representing normalized path for client.
	// (e.g. it can be full path with host prefix for http client)
	Path(path string) string
}
