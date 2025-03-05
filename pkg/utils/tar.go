package utils

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
	"github.com/spf13/afero"
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

// TarFile creates a tar archive of a file and returns it as a reader
// reader should be closed by the user
func TarFile(ctx context.Context, tarfilename string, filename string, filecontents []byte) (io.Reader, error) {
	var archiveFiles []archives.FileInfo
	// create a in memory input folder
	aferoFs := afero.NewMemMapFs()
	err := afero.WriteFile(aferoFs, filename, filecontents, 0644)
	if err != nil {
		return nil, err
	}

	err = afero.Walk(aferoFs, "", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Println("Adding file to archive:", path)
			archiveFiles = append(archiveFiles, archives.FileInfo{
				FileInfo:      info,
				NameInArchive: path,
				// !! NOTE: VERY IMPORTANT: Open function to read the file from the aferoFs
				Open: func() (fs.File, error) {
					return aferoFs.Open(path)
				},
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// create a in memory output file
	aferofsFile, err := afero.NewMemMapFs().Create(tarfilename)
	if err != nil {
		return nil, err
	}

	format := archives.CompressedArchive{
		Compression: archives.Gz{},
		Archival:    archives.Tar{},
	}

	// create the archive
	err = format.Archive(ctx, aferofsFile, archiveFiles)
	if err != nil {
		return nil, err
	}

	// Seek to beginning of file to return it
	_, err = aferofsFile.Seek(0, 0)

	return aferofsFile, nil
}
