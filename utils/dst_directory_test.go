package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/utils"
)

func TestIsClusterDir_Normal(t *testing.T) {
	isCluster, err := utils.IsClusterDir("./testdata/worlds/cluster01")
	assert.Nil(t, err)
	assert.True(t, isCluster)
}

func TestIsClusterDir_Invalid(t *testing.T) {
	isCluster, err := utils.IsClusterDir("./testdata/worlds/cluster-no-config")
	assert.Nil(t, err)
	assert.False(t, isCluster)
}

func TestIsShardDir_Normal(t *testing.T) {
	isShard, err := utils.IsShardDir("./testdata/worlds/cluster01/shard-ok")
	assert.Nil(t, err)
	assert.True(t, isShard)
}

func TestIsShardDir_Invalid(t *testing.T) {
	isShard, err := utils.IsShardDir("./testdata/worlds/cluster01/shard-no-config")
	assert.Nil(t, err)
	assert.False(t, isShard)
}

func TestListAllClusters_Normal(t *testing.T) {
	expected := []string{"cluster01", "cluster02"}
	actual, err := utils.ListAllClusters("./testdata/worlds")
	assert.Nil(t, err)
	assert.EqualValues(t, expected, actual)
}

func TestListShards_Normal(t *testing.T) {
	actual, err := utils.ListShards("./testdata/worlds/cluster01")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"shard-ok"}, actual)

	actual, err = utils.ListShards("./testdata/worlds/cluster02")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"shard01", "shard02", "shard03"}, actual)
}
