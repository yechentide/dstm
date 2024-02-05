package cluster

import (
	"gopkg.in/ini.v1"
)

type ClusterConfig struct {
	Network clusterNetworkConfig  `json:"NETWORK"`
	Game    clusterGameplayConfig `json:"GAMEPLAY"`
	Misc    clusterMiscConfig     `json:"MISC"`
	Shard   clusterShardConfig    `json:"SHARD"`
	Steam   clusterSteamConfig    `json:"STEAM"`
}

type clusterNetworkConfig struct {
	Name           string `json:"cluster_name"`
	Desc           string `json:"cluster_description"`
	Password       string `json:"cluster_password"`
	Lang           string `json:"cluster_language"`
	LanOnly        bool   `json:"lan_only_cluster"`
	Offline        bool   `json:"offline_cluster"`
	TickRate       int    `json:"tick_rate"`
	WhitelistSlots int    `json:"whitelist_slots"`
	AutoSave       bool   `json:"autosaver_enabled"`
}

type clusterGameplayConfig struct {
	GameMode    string `json:"game_mode"`
	MaxPlayers  int    `json:"max_players"`
	PvP         bool   `json:"pvp"`
	AutoPause   bool   `json:"pause_when_empty"`
	VoteEnabled bool   `json:"vote_enabled"`
}

type clusterMiscConfig struct {
	ConsoleEnabled bool `json:"console_enabled"`
	MaxSnapshots   int  `json:"max_snapshots"`
}

type clusterShardConfig struct {
	ShardEnabled bool   `json:"shard_enabled"`
	BindIP       string `json:"bind_ip"`
	MasterIP     string `json:"master_ip"`
	MasterPort   int    `json:"master_port"`
	ClusterKey   string `json:"cluster_key"`
}

type clusterSteamConfig struct {
	SteamGroupOnly   bool `json:"steam_group_only"`
	SteamGroupID     int  `json:"steam_group_id"`
	SteamGroupAdmins bool `json:"steam_group_admins"`
}

func (c *ClusterConfig) SaveTo(clusterDirPath string) error {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, c)
	if err != nil {
		return err
	}
	return cfg.SaveTo(clusterDirPath + "/cluster.ini")
}
