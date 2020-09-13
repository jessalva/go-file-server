package saving

import "io"

type FileStore interface {
	Save( filename, post string, f io.Reader ) error
}

type Service interface {
	SaveFile( filename, post string , reader io.Reader ) error
}

type service struct {
	fileStore FileStore
}

func NewService(fileStore FileStore) Service {
	return &service{fileStore: fileStore}
}

func ( s *service ) SaveFile( filename, post string, file io.Reader ) error{

	return s.fileStore.Save( filename, post, file)
	
}
