package cluster

import "gopkg.in/ini.v1"

func ReadClusterINI(clusterDirPath string) (*ClusterConfig, error) {
	iniFile, err := ini.Load(clusterDirPath + "/cluster.ini")
	if err != nil {
		return nil, err
	}
	cfg := MakeDefaultConfig()
	err = applyExistsClusterINI(cfg, iniFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func MakeDefaultConfig() *ClusterConfig {
	return &ClusterConfig{
		Network: clusterNetworkConfig{
			Name:           "Fun playground",
			Desc:           "Managed by dstm",
			Password:       "",
			Lang:           "en",
			LanOnly:        false,
			Offline:        false,
			TickRate:       15,
			WhitelistSlots: 0,
			AutoSave:       true,
		},
		Game: clusterGameplayConfig{
			GameMode:    "survival",
			MaxPlayers:  6,
			PvP:         false,
			AutoPause:   true,
			VoteEnabled: false,
		},
		Misc: clusterMiscConfig{
			ConsoleEnabled: true,
			MaxSnapshots:   6,
		},
		Shard: clusterShardConfig{
			ShardEnabled: true,
			BindIP:       "127.0.0.1",
			MasterIP:     "127.0.0.1",
			MasterPort:   10888,
			ClusterKey:   "defaultPass",
		},
		Steam: clusterSteamConfig{
			SteamGroupOnly:   false,
			SteamGroupID:     0,
			SteamGroupAdmins: false,
		},
	}
}

func applyExistsClusterINI(cfg *ClusterConfig, iniFile *ini.File) error {
	if iniFile.Section("NETWORK").Key("cluster_name").String() != "" {
		cfg.Network.Name = iniFile.Section("NETWORK").Key("cluster_name").String()
	}
	if iniFile.Section("NETWORK").Key("cluster_description").String() != "" {
		cfg.Network.Desc = iniFile.Section("NETWORK").Key("cluster_description").String()
	}
	if iniFile.Section("NETWORK").Key("cluster_password").String() != "" {
		cfg.Network.Password = iniFile.Section("NETWORK").Key("cluster_password").String()
	}
	if iniFile.Section("NETWORK").Key("cluster_language").String() != "" {
		cfg.Network.Lang = iniFile.Section("NETWORK").Key("cluster_language").String()
	}
	if iniFile.Section("NETWORK").Key("lan_only_cluster").String() != "" {
		value, err := iniFile.Section("NETWORK").Key("lan_only_cluster").Bool()
		if err != nil {
			return err
		}
		cfg.Network.LanOnly = value
	}
	if iniFile.Section("NETWORK").Key("offline_cluster").String() != "" {
		value, err := iniFile.Section("NETWORK").Key("offline_cluster").Bool()
		if err != nil {
			return err
		}
		cfg.Network.Offline = value
	}
	if iniFile.Section("NETWORK").Key("tick_rate").String() != "" {
		value, err := iniFile.Section("NETWORK").Key("tick_rate").Int()
		if err != nil {
			return err
		}
		cfg.Network.TickRate = value
	}
	if iniFile.Section("NETWORK").Key("whitelist_slots").String() != "" {
		value, err := iniFile.Section("NETWORK").Key("whitelist_slots").Int()
		if err != nil {
			return err
		}
		cfg.Network.WhitelistSlots = value
	}
	if iniFile.Section("NETWORK").Key("autosaver_enabled").String() != "" {
		value, err := iniFile.Section("NETWORK").Key("autosaver_enabled").Bool()
		if err != nil {
			return err
		}
		cfg.Network.AutoSave = value
	}

	if iniFile.Section("GAMEPLAY").Key("game_mode").String() != "" {
		cfg.Game.GameMode = iniFile.Section("GAMEPLAY").Key("game_mode").String()
	}
	if iniFile.Section("GAMEPLAY").Key("max_players").String() != "" {
		value, err := iniFile.Section("GAMEPLAY").Key("max_players").Int()
		if err != nil {
			return err
		}
		cfg.Game.MaxPlayers = value
	}
	if iniFile.Section("GAMEPLAY").Key("pvp").String() != "" {
		value, err := iniFile.Section("GAMEPLAY").Key("pvp").Bool()
		if err != nil {
			return err
		}
		cfg.Game.PvP = value
	}
	if iniFile.Section("GAMEPLAY").Key("pause_when_empty").String() != "" {
		value, err := iniFile.Section("GAMEPLAY").Key("pause_when_empty").Bool()
		if err != nil {
			return err
		}
		cfg.Game.AutoPause = value
	}
	if iniFile.Section("GAMEPLAY").Key("vote_enabled").String() != "" {
		value, err := iniFile.Section("GAMEPLAY").Key("vote_enabled").Bool()
		if err != nil {
			return err
		}
		cfg.Game.VoteEnabled = value
	}

	if iniFile.Section("MISC").Key("console_enabled").String() != "" {
		value, err := iniFile.Section("MISC").Key("console_enabled").Bool()
		if err != nil {
			return err
		}
		cfg.Misc.ConsoleEnabled = value
	}
	if iniFile.Section("MISC").Key("max_snapshots").String() != "" {
		value, err := iniFile.Section("MISC").Key("max_snapshots").Int()
		if err != nil {
			return err
		}
		cfg.Misc.MaxSnapshots = value
	}

	if iniFile.Section("SHARD").Key("shard_enabled").String() != "" {
		value, err := iniFile.Section("SHARD").Key("shard_enabled").Bool()
		if err != nil {
			return err
		}
		cfg.Shard.ShardEnabled = value
	}
	if iniFile.Section("SHARD").Key("bind_ip").String() != "" {
		cfg.Shard.BindIP = iniFile.Section("SHARD").Key("bind_ip").String()
	}
	if iniFile.Section("SHARD").Key("master_ip").String() != "" {
		cfg.Shard.MasterIP = iniFile.Section("SHARD").Key("master_ip").String()
	}
	if iniFile.Section("SHARD").Key("master_port").String() != "" {
		value, err := iniFile.Section("SHARD").Key("master_port").Int()
		if err != nil {
			return err
		}
		cfg.Shard.MasterPort = value
	}
	if iniFile.Section("SHARD").Key("cluster_key").String() != "" {
		cfg.Shard.ClusterKey = iniFile.Section("SHARD").Key("cluster_key").String()
	}

	if iniFile.Section("STEAM").Key("steam_group_only").String() != "" {
		value, err := iniFile.Section("STEAM").Key("steam_group_only").Bool()
		if err != nil {
			return err
		}
		cfg.Steam.SteamGroupOnly = value
	}
	if iniFile.Section("STEAM").Key("steam_group_id").String() != "" {
		value, err := iniFile.Section("STEAM").Key("steam_group_id").Int()
		if err != nil {
			return err
		}
		cfg.Steam.SteamGroupID = value
	}
	if iniFile.Section("STEAM").Key("steam_group_admins").String() != "" {
		value, err := iniFile.Section("STEAM").Key("steam_group_admins").Bool()
		if err != nil {
			return err
		}
		cfg.Steam.SteamGroupAdmins = value
	}
	return nil
}
