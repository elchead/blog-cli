package blog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elchead/blog-cli/blog"
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
	sut := &blog.FakeSymLinker{Fs:mockedFs}
	// files don't exist!
	assert.Error(t,sut.Symlink("/file/write.md","link.md"))
}

// Afero lib lacks proper symlink support
// func TestCopyLink(t *testing.T) {
// 	osFs := &afero.OsFs{}
// 	workDir, err := afero.TempDir(osFs, "", "afero-symlink")
// 	assert.NoError(t, err)
// 	osPath := filepath.Join(workDir, "afero.txt")
// 	afero.WriteFile(osFs, osPath, []byte("Hi, Afero!"), 0777)
// 	symPath := filepath.Join(workDir, "link.txt")
// 	err = osFs.SymlinkIfPossible(osPath, symPath)
// 	assert.NoError(t, err)
// 	res, err:= filepath.EvalSymlinks(symPath)
// 	assert.NoError(t, err)

// 	orgf,err:= osFs.Open(res)
// 	assert.NoError(t, err)
// 	// osFs.OpenFile()
// 	symf,err := osFs.OpenFile(res,os.O_CREATE|os.O_WRONLY,0755)
// 	assert.NoError(t, err)
// 	_,err = io.Copy(symf, orgf)
// 	assert.NoError(t, err)

// 	osResolved,_ := filepath.EvalSymlinks(osPath) 
// 	assert.Equal(t,osResolved,res)
// 	assert.Equal(t,true,false)

// 	// mockedFs := afero.NewMemMapFs()
// 	// sut := &blog.FakeSymLinker{Fs:mockedFs}
// }


func TestReplaceSymlinkWithHardlink(t *testing.T) {
	link := "link.md"
	os.Remove(link)

	os.Symlink("../README.md",link)
	assert.True(t,blog.IsSymlink(link))

	originalPath,err:=filepath.EvalSymlinks(link)
	assert.NoError(t, err)

	err = blog.MakeHardlink(link,originalPath)
	assert.NoError(t, err)
	assert.False(t,blog.IsSymlink(link))
}

