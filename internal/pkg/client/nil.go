package client

type NilClient struct {}

func (s *NilClient) GetImageContent(path string) ([]byte, error) {
	return []byte{}, nil
}

func (s *NilClient) Path(path string) string {
	return ""
}

