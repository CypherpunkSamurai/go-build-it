package utils

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
)

type TarBallReader struct {
	Reader io.Reader
}

// TarDir creates a tar archive of a directory and returns it as a reader
// reader should be closed by the user
func TarDir(ctx context.Context, dir string) (io.Reader, error) {
	// create a temp file
	tempFile, err := os.CreateTemp("", "go-build-")
	if err != nil {
		return nil, err
	}
	// dont call defer or it will close when this function exits

	// we can use the CompressedArchive type to gzip a tarball
	format := archives.CompressedArchive{
		Archival: archives.Tar{},
	}

	// read all files recursively from the directory
	var fileInfos []archives.FileInfo
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Convert to relative path
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			// Create archives.FileInfo
			fileInfos = append(fileInfos, archives.FileInfo{
				FileInfo:      info,
				NameInArchive: relPath,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// create the archive
	err = format.Archive(ctx, tempFile, fileInfos)
	if err != nil {
		return nil, err
	}

	// Seek to beginning of file to return it
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}
