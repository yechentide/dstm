package world

import (
	"encoding/json"
	"os"

	"github.com/yechentide/dstm/extractor"
)

func ReadWorldgenOverride(shardDirPath, tempJsonPath string) (*WorldConfig, error) {
	currentJson, err := extractor.GenerateWorldgenOverrideJson(shardDirPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(currentJson)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonObject worldgenOverride

	err = json.NewDecoder(file).Decode(&jsonObject)
	if err != nil {
		return nil, err
	}

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
	// TODO: reflect
	return nil
}
