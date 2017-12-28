package datamodels

import (
	"archive/zip"
	"os"
	"io"
	"path"
)

type Yaml struct {
	Group string
	Path  string
}

type Yamls []Yaml

var output = "yamls.zip"

func ZipFile(pwd string) string {
	return path.Join(pwd, output)
}

func (y *Yamls) ZipTo(pwd string) error {
	newfile, err := os.Create(ZipFile(pwd))
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	for _, yaml := range *y {
		zipfile, err := os.Open(yaml.Path)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}
