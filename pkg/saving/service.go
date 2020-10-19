package saving

import (
	"context"
	"io"
)

type FileStore interface {
	Save(ctx context.Context, filename, post string, f io.Reader) error
}

type Service interface {
	SaveFile(ctx context.Context, filename, post string, reader io.Reader) error
}

type service struct {
	fileStore FileStore
}

func NewService(fileStore FileStore) Service {
	return &service{fileStore: fileStore}
}

func (s *service) SaveFile(ctx context.Context, filename, post string, file io.Reader) error {

	return s.fileStore.Save(ctx, filename, post, file)

}
