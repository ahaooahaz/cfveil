package python

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/fs"
	"path/filepath"
	"regexp"
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
			Language: p,
			IsDir:    d.IsDir(),
			AbsPath:  path,
			RelPath:  relPath,
		}

		if f.RelPath != "." && (f.IsDir || strings.HasSuffix(f.RelPath, ".py")) {
			files = append(files, f)
			obsFiles = append(obsFiles, ToObsFile(f, out))
		}
		return ie
	})
	logrus.WithField("obsfiles", obsFiles).Debugf("obsfiles")
	p.obsFiles = obsFiles
	for _, f := range obsFiles {
		err = f.Write()
		if err != nil {
			return
		}
	}

	return
}

func ToObsFile(f *code.File, dst string) (of *code.ObsFile) {
	of = &code.ObsFile{
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

type python struct {
	obsFiles []*code.ObsFile
}

var p = &python{}

func (p *python) Obs(scanner *bufio.Scanner) (b *bytes.Buffer, err error) {
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		re := regexp.MustCompile(`^\s*(from|import)\b.*`)

		if re.MatchString(strings.TrimSpace(line)) {
			for _, x := range p.obsFiles {
				line = strings.ReplaceAll(line, x.ImpName, x.ObsImpName)
			}
		}
		lines = append(lines, line)
	}
	b = bytes.NewBuffer([]byte(strings.Join(lines, "\n")))
	return
}
