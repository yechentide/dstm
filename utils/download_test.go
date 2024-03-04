package utils_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestDownloadFile_Normal(t *testing.T) {
	destPath := "./dstm-download-test"
	err := utils.DownloadFile(destPath, "https://go.dev/help")
	assert.Nil(t, err)

	file, err := os.Open(destPath)
	assert.Nil(t, err)
	defer file.Close()

	reader := bufio.NewReader(file)
	firstLine, _, err := reader.ReadLine()
	assert.Nil(t, err)
	assert.Equal(t, "<!DOCTYPE html>", string(firstLine))

	err = os.Remove(destPath)
	assert.Nil(t, err)
}
