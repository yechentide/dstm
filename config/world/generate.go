package world

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	"github.com/yechentide/dstm/config/cluster"
	"github.com/yechentide/dstm/config/shard"
	"github.com/yechentide/dstm/extractor"
)

func getClusterLanguage(shardDirPath string) string {
	clusterDirPath := path.Dir(shardDirPath)
	config, err := cluster.ReadClusterINI(clusterDirPath)
	if err == nil && config.Network.Lang != "" {
		return config.Network.Lang
	}
	return "en"
}

func getLocation(override *worldgenOverride) string {
	containCave := func(text string) bool {
		return strings.Contains(strings.ToLower(text), "cave")
	}
	if containCave(override.ID) || containCave(override.Name) || containCave(override.Preset) || containCave(override.Location) {
		return "cave"
	}
	return "forest"
}

func getTemplateJsonPath(shardDirPath string, override *worldgenOverride) string {
	lang := getClusterLanguage(shardDirPath)
	location := getLocation(override)
	jsonFileName := lang + "." + location
	shardConfig, err := shard.ReadServerINI(shardDirPath, location)
	if err == nil && shardConfig.Shard.IsMaster {
		jsonFileName += ".master"
	}
	return viper.GetString("cacheDirPath") + "/json/" + jsonFileName + ".json"
}

func ReadWorldgenOverride(shardDirPath string) (*WorldConfig, error) {
	err := extractor.ExtractWorldgenOverride(shardDirPath, shardDirPath)
	if err != nil {
		return nil, err
	}
	currentJsonPath := shardDirPath + "/worldgenoverride.json"
	defer os.Remove(currentJsonPath)

	file, err := os.Open(currentJsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonObject worldgenOverride

	err = json.NewDecoder(file).Decode(&jsonObject)
	if err != nil {
		return nil, err
	}

	tempJsonPath := getTemplateJsonPath(shardDirPath, &jsonObject)
	cfg, err := MakeDefaultConfig(tempJsonPath)
	if err != nil {
		return nil, err
	}
	err = applyExistsWorldgenOverride(cfg, &jsonObject)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func MakeDefaultConfig(tempJsonPath string) (*WorldConfig, error) {
	var cfg WorldConfig

	file, err := os.Open(tempJsonPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	cfg.setAllCurrentDefault()
	return &cfg, nil
}

func applyExistsWorldgenOverride(cfg *WorldConfig, exists *worldgenOverride) error {
	overrides := map[string]interface{}{}
	for k, v := range exists.Overrides {
		overrides[k] = v
	}
	slog.Warn("Not implemented: applyExistsWorldgenOverride()")
	return nil
}
