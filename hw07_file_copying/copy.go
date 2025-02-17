package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("negative offset")
	ErrorNegativeLimit       = errors.New("negative limit")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	log.Println("Copying from", fromPath, "to", toPath, "offset", offset, "limit", limit)

	// check for negative offset or limit
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrorNegativeLimit
	}

	// open input file
	input, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer closeFile(fromPath, input)

	// get file info
	fileInfo, err := input.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	size := fileInfo.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	// can't copy device files or directories etc
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// open output file
	output, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer closeFile(toPath, output)

	_, err = input.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking: %w", err)
	}

	sizeWithoutOffset := size - offset
	copyAmount := limit
	// if limit is 0 or greater than the size of the file without the offset, copy the whole file
	if limit == 0 || limit > sizeWithoutOffset {
		copyAmount = sizeWithoutOffset
	}

	bar := pb.Simple.Start64(copyAmount)
	barReader := bar.NewProxyReader(input)

	copied, err := io.CopyN(output, barReader, copyAmount)
	if err != nil {
		return fmt.Errorf("error copying: %w", err)
	}
	bar.Finish()

	log.Printf("Copied %d bytes", copied)

	return nil
}

func closeFile(path string, file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing file %s: %v", path, err)
	}
}
