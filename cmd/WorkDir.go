package cmd

import (
	"os"
	"path"
	"archive/zip"
	"io"
	"io/ioutil"
	"github.com/softleader/deployer/datamodels"
)

var (
	deployedDir    = "deployed"
	compressOutput = "deployed.zip"
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

	d := path.Join(wd.Path, deployedDir)
	reMkdir(d)
	return &wd
}

func reMkdir(path string) {
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModeDir|os.ModePerm)
}

func (wd *WorkDir) MoveToDeployedDir(files []datamodels.Yaml) error {
	for _, f := range files {
		newpath := path.Join(wd.Path, deployedDir, path.Base(f.Path))
		err := os.Rename(f.Path, newpath)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (wd *WorkDir) GetCompressPath() string {
	return path.Join(wd.Path, deployedDir, compressOutput)
}

func (wd *WorkDir) CompressDeployedDir() error {
	newfile, err := os.Create(wd.GetCompressPath())
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()
	d, err := ioutil.ReadDir(path.Join(wd.Path, deployedDir))
	if err != nil {
		return err
	}

	for _, f := range d {
		if f.Name() == compressOutput {
			continue
		}

		zipfile, err := os.Open(path.Join(wd.Path, deployedDir, f.Name()))
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
