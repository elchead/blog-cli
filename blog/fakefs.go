package blog

import (
	"fmt"
	"path"

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
