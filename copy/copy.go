package copy

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Copy(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	return copy(src, dst, info)
}

func copy(src, dst string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dst, info)
	}

	if info.IsDir() {
		return dcopy(src, dst, info)
	}

	return fcopy(src, dst, info)
}

func fcopy(src, dst string, info os.FileInfo) error {
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)

	return err
}

func dcopy(sdir, ddir string, info os.FileInfo) error {
	if err := os.MkdirAll(ddir, info.Mode()); err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(sdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(sdir, content.Name()), filepath.Join(ddir, content.Name())
		if err := copy(cs, cd, content); err != nil {
			return err
		}
	}

	return nil
}

func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(src, dest)
}
