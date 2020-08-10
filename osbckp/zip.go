package osbckp

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joivo/osbckp/util"
	"github.com/nuveo/log"
)

func Zip(source, target string) error {
	log.Printf("Zipping %s into %s\n", source, target)

	zipfile, err := os.Create(target)
	util.HandleErr(err)
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	util.HandleErr(err)

	info, err := os.Stat(source)
	util.HandleErr(err)

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)

		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		log.Errorln(err.Error())
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			log.Errorln(err.Error())
			return err
		}
	}

	return nil
}
