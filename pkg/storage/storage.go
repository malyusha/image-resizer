package storage

// Storage represents storage interface for image files manipulation (retrieving, storing, deleting)
type Storage interface {
	// Saves given content into a file with given name. Returns new path to file or error.
	Save(filename string, content []byte) (string, error)
	// Returns content of file by given name. Must return error if any error occurred while
	// opening/reading file content
	Get(filename string) ([]byte, error)
	// Deletes file by given name
	Delete(filename string) error
	// Clears all previously stored files
	Purge() error
}
