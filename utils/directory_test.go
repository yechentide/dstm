package utils_test

import (
	"embed"
	"io/fs"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestDirExists_Exists(t *testing.T) {
	dirPath := "./dstm-test"
	err := os.RemoveAll(dirPath)
	assert.Nil(t, err)

	err = os.MkdirAll(dirPath, 0755)
	assert.Nil(t, err)

	exists, err := utils.DirExists(dirPath)
	assert.Nil(t, err)
	assert.True(t, exists)

	err = os.RemoveAll(dirPath)
	assert.Nil(t, err)
}

func TestDirExists_ExistsButIsNotDir(t *testing.T) {
	filePath := "./dstm-test"
	err := exec.Command("bash", "-c", "echo 'abc' > "+filePath).Run()
	assert.Nil(t, err)

	exists, err := utils.DirExists(filePath)
	assert.EqualError(t, err, "file is not a directory: "+filePath)
	assert.False(t, exists)

	err = os.Remove(filePath)
	assert.Nil(t, err)
}

func TestDirExists_NotExists(t *testing.T) {
	dirPath := "./dstm-test"
	err := os.RemoveAll(dirPath)
	assert.Nil(t, err)

	exists, err := utils.DirExists(dirPath)
	assert.Nil(t, err)
	assert.False(t, exists)

	err = os.RemoveAll(dirPath)
	assert.Nil(t, err)
}

func TestRemakeDir_Exists(t *testing.T) {
	const mode fs.FileMode = 0755
	dirPath := "./dstm-test"
	err := os.RemoveAll(dirPath)
	assert.Nil(t, err)

	err = os.MkdirAll(dirPath+"/child", mode)
	assert.Nil(t, err)

	err = utils.RemakeDir(dirPath, mode, false)
	assert.Nil(t, err)

	dir, err := os.Stat(dirPath)
	assert.Nil(t, err)
	assert.Equal(t, dir.Mode().Perm(), mode.Perm())
	_, err = os.Stat(dirPath + "/child")
	assert.True(t, os.IsNotExist(err))

	err = os.RemoveAll(dirPath)
	assert.Nil(t, err)
}

func TestRemakeDir_NotExists(t *testing.T) {
	const mode fs.FileMode = 0755
	dirPath := "./dstm-test"
	err := os.RemoveAll(dirPath)
	assert.Nil(t, err)

	err = utils.RemakeDir(dirPath, mode, false)
	assert.Nil(t, err)

	dir, err := os.Stat(dirPath)
	assert.Nil(t, err)
	assert.Equal(t, dir.Mode().Perm(), mode.Perm())

	err = os.RemoveAll(dirPath)
	assert.Nil(t, err)
}

func TestRemakeDir_Recursive(t *testing.T) {
	const mode fs.FileMode = 0755
	parentPath := "./dstm-test"
	targetPath := parentPath + "/child/abc"

	err := os.RemoveAll(parentPath)
	assert.Nil(t, err)

	err = utils.RemakeDir(targetPath, mode, true)
	assert.Nil(t, err)

	info, err := os.Stat(targetPath)
	assert.Nil(t, err)
	assert.Equal(t, "abc", info.Name())
	assert.Equal(t, info.Mode().Perm(), mode.Perm())

	err = os.RemoveAll(parentPath)
	assert.Nil(t, err)
}

func TestCopyDir_Exists(t *testing.T) {
	const mode fs.FileMode = 0755
	srcPath := "./dstm-test01"
	destPath := "./dstm-test02"

	err := os.MkdirAll(srcPath+"/child", mode)
	assert.Nil(t, err)
	err = os.RemoveAll(destPath)
	assert.Nil(t, err)

	err = utils.CopyDir(srcPath, destPath)
	assert.Nil(t, err)

	dir, err := os.Stat(srcPath + "/child")
	assert.Nil(t, err)
	assert.Equal(t, dir.Mode().Perm(), mode.Perm())
	dir, err = os.Stat(destPath)
	assert.Nil(t, err)
	assert.Equal(t, dir.Mode().Perm(), mode.Perm())

	err = os.RemoveAll(srcPath)
	assert.Nil(t, err)
	err = os.RemoveAll(destPath)
	assert.Nil(t, err)
}

func TestCopyDir_NotExists(t *testing.T) {
	srcPath := "./dstm-test01"
	destPath := "./dstm-test02"

	err := os.RemoveAll(srcPath)
	assert.Nil(t, err)
	err = os.RemoveAll(destPath)
	assert.Nil(t, err)

	err = utils.CopyDir(srcPath, destPath)
	assert.True(t, os.IsNotExist(err))

	err = os.RemoveAll(srcPath)
	assert.Nil(t, err)
	err = os.RemoveAll(destPath)
	assert.Nil(t, err)
}

func TestListChildDirs_Exists(t *testing.T) {
	dirPath := "./dstm-test"
	err := os.RemoveAll(dirPath)
	assert.Nil(t, err)

	childDirs := []string{"dir01", "dir02", "dir03"}
	for _, childDir := range childDirs {
		err = os.MkdirAll(dirPath+"/"+childDir, 0755)
		assert.Nil(t, err)
	}

	children, err := utils.ListChildDirs(dirPath)
	assert.Nil(t, err)
	assert.EqualValues(t, childDirs, children)

	err = os.RemoveAll(dirPath)
	assert.Nil(t, err)
}

func TestListChildDirs_NotExists(t *testing.T) {
	dirPath := "./dstm-test"
	err := os.RemoveAll(dirPath)
	assert.Nil(t, err)

	children, err := utils.ListChildDirs(dirPath)
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, children)
}

//go:embed testdata/sample
var testDir embed.FS

func TestCopyEmbeddedDir_WholeDir(t *testing.T) {
	err := utils.CopyEmbeddedDir(testDir, "testdata/sample", "./sample-test")
	assert.Nil(t, err)

	checkExistenceAndPermission(t, "./sample-test/a.txt", false, 0644)
	checkExistenceAndPermission(t, "./sample-test/demo01/b.txt", false, 0644)
	checkExistenceAndPermission(t, "./sample-test/demo01/demo03/c.txt", false, 0644)

	emptyDir, err := os.Stat("./sample-test/demo02")
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, emptyDir)

	err = os.RemoveAll("./sample-test")
	assert.Nil(t, err)
}

func TestCopyEmbeddedDir_SubDir(t *testing.T) {
	err := utils.CopyEmbeddedDir(testDir, "testdata/sample/demo01/demo03", "./sample-test")
	assert.Nil(t, err)

	checkExistenceAndPermission(t, "./sample-test/c.txt", false, 0644)

	err = os.RemoveAll("./sample-test")
	assert.Nil(t, err)
}
