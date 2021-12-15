package fs

import (
	"os"

	"github.com/spf13/afero"
)

type Filesystem struct {}

func (f Filesystem) Symlink(target,link string) error {
	return os.Symlink(target,link)
}
func (f Filesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
func (f Filesystem) Create(path string) (afero.File,error) {
	return os.Create(path)
}

func (f Filesystem) Open(path string) (afero.File,error) {
	return os.Open(path)
}
