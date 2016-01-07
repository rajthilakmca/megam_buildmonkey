package chefsolo

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

type Torr struct {
	Source string
}

func NewTorr(sourcefile string) *Torr {
	return &Torr{Source: sourcefile}
}

func (t *Torr) untar() error {
	file, err := os.Open(t.Source)
	if err != nil {
		return err
	}

	defer file.Close()
	var fileReader io.ReadCloser = file

	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(t.Source, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			return err
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)
	// Extracting tarred files
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Println("Creating directory :", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				return err
			}

		case tar.TypeReg:
			// handle normal file
			fmt.Println("Untarring :", filename)
			writer, err := os.Create(filename)

			if err != nil {
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				return err
			}

			writer.Close()
		default:
			fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}
	return nil
}

func (t *Torr) cleanup() error {
	if _, err := os.Stat(t.Source); err == nil {
		if err = os.Remove(t.Source); err != nil {
			return err
		}
	}
	return nil
}
