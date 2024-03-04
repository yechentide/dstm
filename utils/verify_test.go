package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestNotEmpty(t *testing.T) {
	var err error

	err = utils.NotEmpty("a")
	assert.Nil(t, err)

	err = utils.NotEmpty("")
	assert.EqualError(t, err, "cannot be empty")
}

func TestNotContainSpace(t *testing.T) {
	errMsg := "cannot contain space"
	var err error

	err = utils.NotContainSpace("")
	assert.Nil(t, err)

	err = utils.NotContainSpace("aaa")
	assert.Nil(t, err)

	err = utils.NotContainSpace(" ")
	assert.EqualError(t, err, errMsg)

	err = utils.NotContainSpace("a ")
	assert.EqualError(t, err, errMsg)

	err = utils.NotContainSpace(" a")
	assert.EqualError(t, err, errMsg)

	err = utils.NotContainSpace("a a")
	assert.EqualError(t, err, errMsg)
}

func TestIsClusterToken(t *testing.T) {
	errMsg := "not a valid cluster token"
	var err error

	err = utils.IsClusterToken("pds-g^KU_")
	assert.Nil(t, err)

	err = utils.IsClusterToken("pds-g^KU_aaa")
	assert.Nil(t, err)

	err = utils.IsClusterToken("")
	assert.EqualError(t, err, errMsg)

	err = utils.IsClusterToken("pds-g^KU")
	assert.EqualError(t, err, errMsg)
}

func TestUnique(t *testing.T) {
	var err error

	err = utils.Unique([]string{})("c")
	assert.Nil(t, err)

	err = utils.Unique([]string{"a", "b"})("c")
	assert.Nil(t, err)

	err = utils.Unique([]string{"a", "b"})("a")
	assert.EqualError(t, err, "already exists")
}

func TestIsPort(t *testing.T) {
	var err error

	err = utils.IsPort("0")
	assert.Nil(t, err)

	err = utils.IsPort("65535")
	assert.Nil(t, err)

	err = utils.IsPort("-1")
	assert.EqualError(t, err, "port out of range")

	err = utils.IsPort("65536")
	assert.EqualError(t, err, "port out of range")

	err = utils.IsPort("a80")
	assert.NotNil(t, err)

	err = utils.IsPort("80a")
	assert.NotNil(t, err)
}
