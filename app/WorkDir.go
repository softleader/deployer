package app

import (
	"os"
	"path"
	"archive/zip"
	"io"
	"io/ioutil"
	"github.com/softleader/deployer/models"
)

var (
	yamlDir        = "yaml"
	compressOutput = "deployment.zip"
)

// 當前的 working directory, 通常是 workspace/${project}
type WorkDir struct {
	Path string
}

func NewWorkDir(cleanUp bool, p string) *WorkDir {
	wd := WorkDir{Path: p}

	if cleanUp {
		reMkdir(wd.Path)
	}
	d := path.Join(wd.Path, yamlDir)
	reMkdir(d)

	return &wd
}

func reMkdir(path string) {
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModeDir|os.ModePerm)
}

func (wd *WorkDir) CopyToYamlDir(files []models.Yaml) error {
	for _, f := range files {
		newpath := path.Join(wd.Path, yamlDir, path.Base(f.Path))
		err := copy(f.Path, newpath)
		if err != nil {
			return nil
		}
	}
	return nil
}

func copy(src string, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

func GetCompressPath(workDirPath string) string {
	return path.Join(workDirPath, yamlDir, compressOutput)
}

func (wd *WorkDir) CompressYamlDir() error {
	newfile, err := os.Create(GetCompressPath(wd.Path))
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()
	d, err := ioutil.ReadDir(path.Join(wd.Path, yamlDir))
	if err != nil {
		return err
	}

	for _, f := range d {
		if f.Name() == compressOutput {
			continue
		}

		zipfile, err := os.Open(path.Join(wd.Path, yamlDir, f.Name()))
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
