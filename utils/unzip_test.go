package utils_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

/*
sample
├── a.txt
├── demo01
│   ├── b.txt
│   └── demo03
│       └── c.txt
└── demo02

4 directories, 3 files
*/

func TestUnzip(t *testing.T) {
	srcPath := "./testdata/sample.zip"
	destPath := "./test"
	const dirMode fs.FileMode = 0755
	const fileMode fs.FileMode = 0644
	var err error

	err = os.RemoveAll(destPath)
	assert.Nil(t, err)

	err = utils.Unzip(srcPath, destPath, dirMode)
	assert.Nil(t, err)

	checkExistenceAndPermission(t, destPath+"/sample/demo01/demo03", true, dirMode)
	checkExistenceAndPermission(t, destPath+"/sample/demo02", true, dirMode)

	checkExistenceAndPermission(t, destPath+"/sample/a.txt", false, fileMode)
	checkExistenceAndPermission(t, destPath+"/sample/demo01/b.txt", false, fileMode)
	checkExistenceAndPermission(t, destPath+"/sample/demo01/demo03/c.txt", false, fileMode)

	err = os.RemoveAll(destPath)
	assert.Nil(t, err)
}

func checkExistenceAndPermission(t *testing.T, filePath string, isDir bool, mode os.FileMode) {
	info, err := os.Stat(filePath)
	assert.Nil(t, err)
	assert.Equal(t, isDir, info.IsDir())
	assert.Equal(t, mode.Perm(), info.Mode().Perm())
}
