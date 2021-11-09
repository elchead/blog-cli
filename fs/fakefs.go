package fs

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

type FakeSymLinker struct {
	afero.Fs
	TargetFile string
	CreatedSymlink string
}

func assertDirExists(fs afero.Fs, path string) error {
	exists, err := afero.Exists(fs, path)
	if !exists || err != nil {
		return fmt.Errorf("directory (%s) does not exist: %w",path,err)
	}
	return nil
}

func (f* FakeSymLinker) Symlink(target, link string) error {
	if err := assertDirExists(f,path.Dir(link)); err != nil {
		return err
	}
	
	if err:= assertDirExists(f,path.Dir(target)); err != nil {
		return err
	}
	f.TargetFile = target
	f.CreatedSymlink = link
	f.Create(link)
	return nil
}

func IsSymlink(link string) bool {
	info, _ := os.Lstat(link)
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func MakeHardlink(link string) error {
	target,err:=filepath.EvalSymlinks(link)
	if err != nil {
		return err
	}
	os.Remove(link)
	linkFile, err := os.OpenFile(link,os.O_CREATE|os.O_WRONLY,0777)
	if err != nil {
		return err
	}
	targetFile, err := os.Open(target)
	if err != nil {
		return err
	}
	_,err = io.Copy(linkFile, targetFile)
	return err
}
