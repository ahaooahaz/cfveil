package code

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
)

type Language interface {
	Obs(*bufio.Scanner) (*bytes.Buffer, error)
}

type File struct {
	Language Language
	IsDir    bool
	AbsPath  string
	RelPath  string
}

type ObsFile struct {
	Src        *File
	ImpName    string
	ObsImpName string
	ObsAbsPath string
}

func (f *ObsFile) Write() (err error) {
	if f.Src.IsDir {
		return
	}

	file, err := os.Open(f.Src.AbsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return err
	}
	buf, err := f.Src.Language.Obs(scanner)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(f.ObsAbsPath), os.ModePerm)
	if err != nil {
		return
	}
	return os.WriteFile(f.ObsAbsPath, buf.Bytes(), fileInfo.Mode())
}
