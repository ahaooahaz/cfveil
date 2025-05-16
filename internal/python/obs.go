package python

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/ahaooahaz/cfveil/internal/code"
	"github.com/sirupsen/logrus"
)

func Process(in, out string, excludes []string) (err error) {
	root, err := filepath.Abs(in)
	if err != nil {
		return
	}

	excludesAbsPaths := make(map[string]bool)
	for _, v := range excludes {
		p := filepath.Join(root, v)
		excludesAbsPaths[p] = true
	}
	logrus.WithField("excludesAbsPaths", excludesAbsPaths).Debugf("excludes abs paths")

	files := []*code.File{}
	obsFiles := []*code.ObsFile{}
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, ie error) error {
		if ie != nil {
			return ie
		}

		if _, ok := excludesAbsPaths[path]; ok {
			return filepath.SkipDir
		}

		relPath, _ := filepath.Rel(root, path)
		f := &code.File{
			Language: code.PYTHON,
			IsDir:    d.IsDir(),
			AbsPath:  path,
			RelPath:  relPath,
		}

		if f.RelPath != "." && (f.IsDir || strings.HasSuffix(f.RelPath, ".py")) {
			files = append(files, f)
			obsFiles = append(obsFiles, f.ToObsFile(out))
		}
		return ie
	})
	logrus.WithField("obsfiles", obsFiles).Debugf("obsfiles")

	return
}
