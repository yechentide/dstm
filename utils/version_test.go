package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestCompareVersion_Normal(t *testing.T) {
	var result int
	var err error

	result, err = utils.CompareVersion("v1.2.3", "v1.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, 0, result)

	result, err = utils.CompareVersion("v1.2.9", "v1.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, 1, result)

	result, err = utils.CompareVersion("v1.9.3", "v1.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, 1, result)

	result, err = utils.CompareVersion("v9.2.3", "v1.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, 1, result)

	result, err = utils.CompareVersion("V1.2.3", "1.2.9")
	assert.Nil(t, err)
	assert.EqualValues(t, -1, result)

	result, err = utils.CompareVersion("V1.2.3", "1.9.3")
	assert.Nil(t, err)
	assert.EqualValues(t, -1, result)

	result, err = utils.CompareVersion("V1.2.3", "9.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, -1, result)
}

func TestCompareVersion_NotSemantic(t *testing.T) {
	var result int
	var err error

	result, err = utils.CompareVersion("v1.2.3.4", "1.2.3")
	assert.NotNil(t, err)
	assert.EqualValues(t, -999, result)

	result, err = utils.CompareVersion("v1.2.3", "1.2")
	assert.NotNil(t, err)
	assert.EqualValues(t, -999, result)
}

func TestGetVersionNumbers_Normal(t *testing.T) {
	var nums []string
	var err error

	nums, err = utils.GetVersionNumbers("v1.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"1", "2", "3"}, nums)

	nums, err = utils.GetVersionNumbers("V1.2.3")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"1", "2", "3"}, nums)

	nums, err = utils.GetVersionNumbers("11.22.33")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"11", "22", "33"}, nums)
}

func TestGetVersionNumbers_NotSemantic(t *testing.T) {
	var nums []string
	var err error

	nums, err = utils.GetVersionNumbers("v1.2")
	assert.EqualError(t, err, "invalid version number: v1.2")
	assert.Nil(t, nums)

	nums, err = utils.GetVersionNumbers("1.2.3.4")
	assert.EqualError(t, err, "invalid version number: 1.2.3.4")
	assert.Nil(t, nums)

	nums, err = utils.GetVersionNumbers("v1.2.3-beta")
	assert.EqualError(t, err, "invalid version number: v1.2.3-beta")
	assert.Nil(t, nums)
}
