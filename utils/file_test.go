package utils_test

import (
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestFileExists_Exists(t *testing.T) {
	filePath := "./dstm-test.txt"
	err := exec.Command("bash", "-c", "echo 'abc' > "+filePath).Run()
	assert.Nil(t, err)

	result, err := utils.FileExists(filePath)
	assert.Nil(t, err)
	assert.True(t, result)

	err = os.Remove(filePath)
	assert.Nil(t, err)
}

func TestFileExists_ExistsButIsNotFile(t *testing.T) {
	dirPath := "./dstm-test"
	err := os.MkdirAll(dirPath, 0755)
	assert.Nil(t, err)

	result, err := utils.FileExists(dirPath)
	assert.EqualError(t, err, "file is a directory: "+dirPath)
	assert.False(t, result)

	err = os.RemoveAll(dirPath)
	assert.Nil(t, err)
}

func TestFileExists_NotExists(t *testing.T) {
	filePath := "./dstm-test.txt"
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		assert.Fail(t, err.Error())
	}

	result, err := utils.FileExists(filePath)
	assert.Nil(t, err)
	assert.False(t, result)
}

func TestWriteToFile_NewFile(t *testing.T) {
	filePath := "./dstm-test.txt"
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		assert.Fail(t, err.Error())
	}

	err = utils.WriteToFile("abc", filePath)
	assert.Nil(t, err)

	file, err := os.Open(filePath)
	assert.Nil(t, err)
	defer file.Close()

	bytes, err := io.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, "abc", string(bytes))

	err = os.Remove(filePath)
	assert.Nil(t, err)
}

func TestWriteToFile_Overwrite(t *testing.T) {
	filePath := "./dstm-test.txt"
	err := exec.Command("bash", "-c", "echo 'abc' > "+filePath).Run()
	assert.Nil(t, err)

	err = utils.WriteToFile("def", filePath)
	assert.Nil(t, err)

	file, err := os.Open(filePath)
	assert.Nil(t, err)
	defer file.Close()

	bytes, err := io.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, "def", string(bytes))

	err = os.Remove(filePath)
	assert.Nil(t, err)
}

func TestCopyFile_SourceExists(t *testing.T) {
	srcPath := "./dstm-test01.txt"
	destPath := "./dstm-test02.txt"
	err := exec.Command("bash", "-c", "echo 'abc' > "+srcPath).Run()
	assert.Nil(t, err)
	err = os.RemoveAll(destPath)
	assert.Nil(t, err)

	err = utils.CopyFile(srcPath, destPath)
	assert.Nil(t, err)

	file, err := os.Open(destPath)
	assert.Nil(t, err)
	defer file.Close()

	bytes, err := io.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, "abc\n", string(bytes))

	err = os.Remove(srcPath)
	assert.Nil(t, err)
	err = os.Remove(destPath)
	assert.Nil(t, err)
}

func TestCopyFile_SourceDoesNotExists(t *testing.T) {
	srcPath := "./dstm-test01.txt"
	destPath := "./dstm-test02.txt"
	err := os.RemoveAll(srcPath)
	assert.Nil(t, err)
	err = os.RemoveAll(destPath)
	assert.Nil(t, err)

	err = utils.CopyFile(srcPath, destPath)
	assert.True(t, os.IsNotExist(err))
}

func TestCopyFile_OverwriteExistsDestFile(t *testing.T) {
	srcPath := "./dstm-test01.txt"
	destPath := "./dstm-test02.txt"
	err := exec.Command("bash", "-c", "echo 'abc' > "+srcPath).Run()
	assert.Nil(t, err)
	err = exec.Command("bash", "-c", "echo 'def' > "+destPath).Run()
	assert.Nil(t, err)

	err = utils.CopyFile(srcPath, destPath)
	assert.Nil(t, err)

	file, err := os.Open(destPath)
	assert.Nil(t, err)
	defer file.Close()

	bytes, err := io.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, "abc\n", string(bytes))

	err = os.Remove(srcPath)
	assert.Nil(t, err)
	err = os.Remove(destPath)
	assert.Nil(t, err)
}

func TestCopyEmbeddedFile_Normal(t *testing.T) {
	srcPath := "testdata/sample/demo01/demo03/c.txt"
	destPath := "./dstm-test.txt"

	err := utils.CopyEmbeddedFile(testDir, srcPath, destPath)
	assert.Nil(t, err)

	file, err := os.Stat(destPath)
	assert.Nil(t, err)
	assert.Equal(t, file.Name(), "dstm-test.txt")

	srcBytes, err := os.ReadFile(srcPath)
	assert.Nil(t, err)
	dstBytes, err := os.ReadFile(destPath)
	assert.Nil(t, err)
	assert.Equal(t, srcBytes, dstBytes)

	err = os.Remove(destPath)
	assert.Nil(t, err)
}

func TestCopyEmbeddedFile_SourceDoesNotExists(t *testing.T) {
	srcPath := "testdata/sample/demo01/demo03/z.txt"
	destPath := "./dstm-test.txt"

	err := utils.CopyEmbeddedFile(testDir, srcPath, destPath)
	assert.True(t, os.IsNotExist(err))
}
