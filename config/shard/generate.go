package shard

import "gopkg.in/ini.v1"

func ReadServerINI(location, shardDir string) (*ShardConfig, error) {
	iniFile, err := ini.Load(shardDir + "/server.ini")
	if err != nil {
		return nil, err
	}

	isMaster := false
	if iniFile.Section("SHARD").Key("is_master").String() != "" {
		value, err := iniFile.Section("SHARD").Key("is_master").Bool()
		if err == nil && value {
			isMaster = true
		}
	}

	cfg := MakeDefaultConfig(location, isMaster)
	err = applyExistsServerINI(cfg, iniFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func MakeDefaultConfig(location string, isMaster bool) *ShardConfig {
	var name string
	var id, serverPort, masterServerPort, authenticationPort int
	if location == "forest" {
		name = "Forest"
		id = 6661000
		serverPort = 11000
		masterServerPort = 27018
		authenticationPort = 8768
	} else {
		name = "Cave"
		id = 6661001
		serverPort = 11001
		masterServerPort = 27019
		authenticationPort = 8769
	}
	if isMaster {
		id = 1
	}
	return &ShardConfig{
		Network: shardNetworkConfig{
			ServerPort: serverPort,
		},
		Shard: shardShardConfig{
			IsMaster: isMaster,
			Name:     name,
			ID:       id,
		},
		Account: shardAccountConfig{
			EncodeUserPath: true,
		},
		Steam: shardSteanConfig{
			MasterServerPort:   masterServerPort,
			AuthenticationPort: authenticationPort,
		},
	}
}

func applyExistsServerINI(cfg *ShardConfig, iniFile *ini.File) error {
	if iniFile.Section("NETWORK").Key("server_port").String() != "" {
		value, err := iniFile.Section("NETWORK").Key("server_port").Int()
		if err != nil {
			return err
		}
		cfg.Network.ServerPort = value
	}

	if iniFile.Section("SHARD").Key("is_master").String() != "" {
		value, err := iniFile.Section("SHARD").Key("is_master").Bool()
		if err != nil {
			return err
		}
		cfg.Shard.IsMaster = value
	}
	if iniFile.Section("SHARD").Key("name").String() != "" {
		cfg.Shard.Name = iniFile.Section("SHARD").Key("name").String()
	}
	if iniFile.Section("SHARD").Key("id").String() != "" {
		value, err := iniFile.Section("SHARD").Key("id").Int()
		if err != nil {
			return err
		}
		cfg.Shard.ID = value
	}

	if iniFile.Section("STEAM").Key("master_server_port").String() != "" {
		value, err := iniFile.Section("STEAM").Key("master_server_port").Int()
		if err != nil {
			return err
		}
		cfg.Steam.MasterServerPort = value
	}
	if iniFile.Section("STEAM").Key("authentication_port").String() != "" {
		value, err := iniFile.Section("STEAM").Key("authentication_port").Int()
		if err != nil {
			return err
		}
		cfg.Steam.AuthenticationPort = value
	}
	if iniFile.Section("ACCOUNT").Key("encode_user_path").String() != "" {
		value, err := iniFile.Section("ACCOUNT").Key("encode_user_path").Bool()
		if err != nil {
			return err
		}
		cfg.Account.EncodeUserPath = value
	}
	return nil
}
