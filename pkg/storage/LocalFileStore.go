package storage

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/xerrors"
	"io"
	"log"
	"os"
	"path/filepath"
)

type LocalFileStore struct {
	basePath    string
	maxFileSize int
	tracer      opentracing.Tracer
	logger      *log.Logger
}

func NewLocalFileStore(basePath string, maxFileSize int, tracer opentracing.Tracer, logger *log.Logger) *LocalFileStore {

	fileDirectory, err := filepath.Abs(basePath)

	if err != nil {
		logger.Printf("[WARN]: base path doesn't exist")
		return &LocalFileStore{basePath: "./resources/LFS/", maxFileSize: maxFileSize, tracer: tracer, logger: logger}
	}

	return &LocalFileStore{basePath: fileDirectory, tracer:tracer, logger: logger}
}

func (LFS *LocalFileStore) Save(ctx context.Context, filename, postId string, file io.Reader) error {

	//if err != nil {
	//	return xerrors.Errorf("Got error reading file: %s",err.Error())
	//}

	saveFileSpan, ctx := opentracing.StartSpanFromContextWithTracer( ctx, LFS.tracer,"LocalFileStore::Save", opentracing.ChildOf(
		opentracing.SpanFromContext(ctx).Context()))
	defer saveFileSpan.Finish()

	pathToPost := filepath.Join(LFS.basePath, postId)
	LFS.logger.Print(pathToPost)

	if _, err := os.Stat(pathToPost); os.IsNotExist(err) {

		if err = os.MkdirAll(pathToPost, os.ModePerm); err != nil {
			return err
		}

	}

	pathToFile := filepath.Join(pathToPost, filename)

	saveFileSpan.LogFields(opentracingLog.String("Path To File", pathToFile))
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

	LFS.logger.Print(written)
	saveFileSpan.LogFields(opentracingLog.String("Bytes Written", fmt.Sprintf("%v", written)))
	return nil
}
