package code

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

type LANGUAGE string

const (
	PYTHON LANGUAGE = "python"
)

type File struct {
	Language LANGUAGE
	IsDir    bool
	AbsPath  string
	RelPath  string
}

func (f *File) ToObsFile(dst string) (of *ObsFile) {
	of = &ObsFile{
		Src: f,
	}

	s := f.RelPath
	if !f.IsDir {
		s = strings.TrimSuffix(f.RelPath, ".py")
	}
	parts := strings.Split(s, "/")
	obsParts := make([]string, len(parts))
	for i := range parts {
		hash := md5.Sum([]byte(parts[i]))
		obsParts[i] = "obs" + hex.EncodeToString(hash[:])
	}
	of.ImpName = strings.Join(parts, ".")
	of.ObsImpName = strings.Join(obsParts, ".")
	s = strings.Join(obsParts, "/")
	if !f.IsDir {
		s += ".py"
	}
	of.ObsAbsPath, _ = filepath.Abs(filepath.Join(dst, s))
	return
}

type ObsFile struct {
	Src        *File
	ImpName    string
	ObsImpName string
	ObsAbsPath string
}

func (f *ObsFile) Write() (err error) {
	err = os.MkdirAll(filepath.Dir(f.ObsAbsPath), os.ModePerm)
	if err != nil {
		return
	}
	return
}
