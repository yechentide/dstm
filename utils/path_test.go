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

func TestGetModIDFromPath_Normal(t *testing.T) {
	var modID string
	var err error

	modID, err = utils.GetModIDFromPath("./testdata/mods/1122334455")
	assert.Nil(t, err)
	assert.EqualValues(t, "1122334455", modID)

	modID, err = utils.GetModIDFromPath("./testdata/mods/workshop-1122334455")
	assert.Nil(t, err)
	assert.EqualValues(t, "1122334455", modID)
}

func TestGetModIDFromPath_Error(t *testing.T) {
	var modID string
	var err error

	modID, err = utils.GetModIDFromPath("./testdata/mods/")
	assert.NotNil(t, err)
	assert.EqualValues(t, "", modID)

	modID, err = utils.GetModIDFromPath("./testdata/mods/workshop1122334455")
	assert.NotNil(t, err)
	assert.EqualValues(t, "", modID)
}
