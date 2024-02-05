package shard

import (
	"github.com/yechentide/dstm/utils"
	"gopkg.in/ini.v1"
)

type ShardConfig struct {
	Network shardNetworkConfig `ini:"NETWORK"`
	Shard   shardShardConfig   `ini:"SHARD"`
	Steam   shardSteanConfig   `ini:"STEAM"`
	Account shardAccountConfig `ini:"ACCOUNT"`
}

type shardNetworkConfig struct {
	ServerPort int `ini:"server_port"`
}

type shardShardConfig struct {
	IsMaster bool   `ini:"is_master"`
	Name     string `ini:"name"`
	ID       int    `ini:"id"`
}

type shardSteanConfig struct {
	MasterServerPort   int `ini:"master_server_port"`
	AuthenticationPort int `ini:"authentication_port"`
}

type shardAccountConfig struct {
	EncodeUserPath bool `ini:"encode_user_path"`
}

func (s *ShardConfig) SaveTo(shardDirPath string) error {
	shardDirPath = utils.ExpandPath(shardDirPath)
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, s)
	if err != nil {
		return err
	}
	return cfg.SaveTo(shardDirPath + "/server.ini")
}
