package storage

type Storage interface {
	Save(filename string, content []byte) error
	Get(filename string) ([]byte, error)
	Delete(filename string) error
}
