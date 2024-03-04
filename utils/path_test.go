package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestExpandPath(t *testing.T) {
	homeDir := os.Getenv("HOME")
	var path string

	path = utils.ExpandPath("~")
	assert.EqualValues(t, homeDir, path)

	path = utils.ExpandPath("~/abc")
	assert.EqualValues(t, homeDir+"/abc", path)

	path = utils.ExpandPath("$HOME")
	assert.EqualValues(t, homeDir, path)

	path = utils.ExpandPath("$HOME/abc")
	assert.EqualValues(t, homeDir+"/abc", path)

	err := os.Setenv("DSTM_TEST", "/tmp/dstm-test")
	assert.Nil(t, err)

	path = utils.ExpandPath("$DSTM_TEST")
	assert.EqualValues(t, "/tmp/dstm-test", path)

	path = utils.ExpandPath("/tmp$DSTM_TEST/abc")
	assert.EqualValues(t, "/tmp/tmp/dstm-test/abc", path)
}

func TestBuildAbsPathString_Normal(t *testing.T) {
	var absPath string
	var err error

	absPath, err = utils.BuildAbsPathString("/tmp")
	assert.Nil(t, err)
	assert.Equal(t, "\"/tmp\"", absPath)

	absPath, err = utils.BuildAbsPathString("/tmp/")
	assert.Nil(t, err)
	assert.Equal(t, "\"/tmp\"", absPath)

	absPath, err = utils.BuildAbsPathString("/tmp/a'b\"c d")
	assert.Nil(t, err)
	assert.Equal(t, "\"/tmp/a'b\\\"c d\"", absPath)
}

func TestBuildAbsPathString_Error(t *testing.T) {
	errMsg := "path must be absolute"
	var absPath string
	var err error

	absPath, err = utils.BuildAbsPathString("")
	assert.EqualError(t, err, errMsg)
	assert.Equal(t, "", absPath)

	absPath, err = utils.BuildAbsPathString("./aaa")
	assert.EqualError(t, err, errMsg)
	assert.Equal(t, "", absPath)

	absPath, err = utils.BuildAbsPathString("../aaa")
	assert.EqualError(t, err, errMsg)
	assert.Equal(t, "", absPath)
}
