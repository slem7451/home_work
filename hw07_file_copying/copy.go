package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileFromStats, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	fileToStats, err := os.Stat(toPath)

	if err == nil && os.SameFile(fileFromStats, fileToStats) {
		tmpPath := toPath + "_copy"
		err := os.Rename(fromPath, tmpPath)
		if err != nil {
			return err
		}
		fromPath = tmpPath
		defer os.Remove(fromPath)
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	size := fileFromStats.Size()

	if !fileFromStats.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if size < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > size-offset {
		limit = size - offset
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	readSection := io.NewSectionReader(fileFrom, offset, limit)

	bar := pb.New(int(limit))
	bar.Start()
	reader := bar.NewProxyReader(readSection)

	_, err = io.CopyN(fileTo, reader, limit)
	bar.Finish()
	if err != nil && errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
