package fs_test

import (
	"os"
	"testing"

	"github.com/elchead/blog-cli/blog"
	"github.com/elchead/blog-cli/fs"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	t.Run("save file",func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		title := "title"
		path := "blog-cli"
		appFS.Create(blog.GetFilepath(title,path))
		exists, _ := afero.Exists(appFS,path+"/title.md")
		assert.Equal(t,true,exists)
	})
	t.Run("non existent file is not saved", func(t *testing.T) {
		title := "title"
		path := "blog-cli"
		appFS := afero.NewMemMapFs()
		appFS.Create(blog.GetFilepath(title,path))
		exists, _ := afero.Exists(appFS,path+"nofile.md")
		assert.Equal(t,false,exists)
	})
}

func TestFakeSymLink(t *testing.T){
	mockedFs := afero.NewMemMapFs()
	sut := &fs.FakeSymLinker{Fs:mockedFs}
	// files don't exist!
	assert.Error(t,sut.Symlink("/file/write.md","link.md"))
}


func TestReplaceSymlinkWithHardlink(t *testing.T) {
	link := "link.md"
	os.Symlink("../README.md",link)
	assert.True(t,fs.IsSymlink(link))
	assert.NoError(t, fs.MakeHardlink(link))
	assert.False(t,fs.IsSymlink(link))
	os.Remove(link)
}

