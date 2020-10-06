package storage

import (
	"golang.org/x/xerrors"
	"io"
	"log"
	"os"
	"path/filepath"
)

type LocalFileStore struct {
	basePath    string
	maxFileSize int
}

func NewLocalFileStore(basePath string, maxFileSize int) *LocalFileStore {

	fileDirectory, err := filepath.Abs(basePath)

	if err != nil {
		log.Printf("[WARN]: base path doesn't exist")
		return &LocalFileStore{basePath: "./resources/LFS/", maxFileSize: maxFileSize}
	}

	return &LocalFileStore{basePath: fileDirectory}
}

func (LFS *LocalFileStore) Save(filename, postId string, file io.Reader) error {

	//if err != nil {
	//	return xerrors.Errorf("Got error reading file: %s",err.Error())
	//}

	pathToPost := filepath.Join(LFS.basePath, postId)
	log.Print(pathToPost)

	if _, err := os.Stat(pathToPost); os.IsNotExist(err) {

		if err = os.MkdirAll(pathToPost, os.ModePerm); err != nil {
			return err
		}

	}

	pathToFile := filepath.Join(pathToPost, filename)
	if _, err := os.Stat(pathToFile); err == nil {
		return err
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Got weird error: %s", err.Error())
	}

	savedFile, err := os.Create(pathToFile)
	if err != nil {
		return xerrors.Errorf("Got error creating file: %s", err.Error())
	}
	defer func() {

		err = savedFile.Close()
		if err != nil {
			err = os.Remove(pathToFile)
		}

	}()

	written, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}

	log.Print(written)
	return nil
}

func (LFS *LocalFileStore) get() {

}
