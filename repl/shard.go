package repl

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/yechentide/dstm/config/shard"
	"github.com/yechentide/dstm/config/world"
	"github.com/yechentide/dstm/extractor"
	"github.com/yechentide/dstm/utils"
	"golang.org/x/exp/slices"
)

func selectShard(clusterDirPath string) string {
	existShards, err := utils.ListShards(clusterDirPath)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	selected := Selector(existShards, "Please select a shard", false)[0]
	return selected
}

func updateShardConfig(config *shard.ShardConfig) {
	// TODO
	slog.Warn("Not implemented: UpdateShardConfig()")
}

func CreateShard() {
	worldsDirPath := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	clusterName := selectCluster(worldsDirPath)
	clusterDirPath := worldsDirPath + "/" + clusterName
	jsonExists, err := utils.IsClusterDir(clusterDirPath)
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
	if !jsonExists {
		printError("Not a valid cluster: " + clusterName)
		os.Exit(1)
	}

	fmt.Println()
	existShards, err := utils.ListShards(clusterDirPath)
	if err == nil {
		msg := "Exist shards:"
		if len(existShards) == 0 {
			msg += " <no shards>"
			printWarn("First shard will be generated as master shard by default.")
		} else {
			for _, shard := range existShards {
				msg += " " + shard
			}
		}
		printInfo(msg)
	} else {
		printError("can not list shards: " + err.Error())
		os.Exit(1)
	}
	shardType := Selector([]string{"forest", "cave"}, "Select shard type", false)[0]
	shardName := "Main"
	if len(existShards) > 0 {
		count := 1
		for _, shard := range existShards {
			name := strings.ToLower(shard)
			if strings.HasPrefix(name, shardType) {
				count++
			}
		}
		shardName = strings.ToUpper(string(shardType[0])) + shardType[1:] + fmt.Sprintf("%02d", count)
	}
	shardDir := clusterDirPath + "/" + shardName
	err = utils.MkDirIfNotExists(shardDir, 0755, false)
	if err != nil {
		printError("Failed to create shard: " + err.Error())
		os.Exit(1)
	}

	// server.ini
	config := shard.MakeDefaultConfig(shardType, len(existShards) == 0)
	updateShardConfig(config)
	err = config.SaveTo(shardDir)
	if err != nil {
		printError("Failed to save shard config: " + err.Error())
		os.Exit(1)
	}
	// worldgenoverride.lua
	outputDirPath := viper.GetString("cacheDirPath") + "/json"
	jsonPath := outputDirPath + "/en.forest.master.json"
	jsonExists, err = utils.FileExists(jsonPath)
	if err != nil {
		printError("Failed to check json file: " + err.Error())
		os.Exit(1)
	}
	if !jsonExists {
		err = extractor.ExtractSettings(viper.GetString("serverRootPath"), outputDirPath)
		if err != nil {
			printError("Failed to extract json file: " + err.Error())
			os.Exit(1)
		}
	}
	override, err := world.MakeDefaultConfig(jsonPath)
	updateWorldOverride(override, true)
	if err != nil {
		printError("Failed to generate worldgenoverride.lua: " + err.Error())
		os.Exit(1)
	}
	err = override.SaveTo(shardDir)
	if err != nil {
		printError("Failed to save worldgenoverride.lua: " + err.Error())
		os.Exit(1)
	}
}

func UpdateShard() {
	worldsDirPath := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	clusterName := selectCluster(worldsDirPath)
	shardName := selectShard(worldsDirPath + "/" + clusterName)
	shardDirPath := worldsDirPath + "/" + clusterName + "/" + shardName
	config, err := shard.ReadServerINI(shardDirPath, "")
	if err != nil {
		printError("Failed to read shard.ini in " + shardDirPath)
		os.Exit(1)
	}
	updateShardConfig(config)
	config.SaveTo(shardDirPath)
}

func UpdateWorld() {
	worldsDirPath := viper.GetString("dataRootPath") + "/" + viper.GetString("worldsDirName")
	clusterName := selectCluster(worldsDirPath)
	shardName := selectShard(worldsDirPath + "/" + clusterName)
	clusterDirPath := worldsDirPath + "/" + clusterName
	shardDirPath := clusterDirPath + "/" + shardName
	config, err := world.ReadWorldgenOverride(shardDirPath)
	if err != nil {
		printError("Failed to read worldgenoverride.lua in " + shardDirPath)
		os.Exit(1)
	}
	updateWorldOverride(config, false)
	config.SaveTo(shardDirPath)
}

func showGroupConfig(group *world.WorldConfigGroup, isGen bool, currentIdx, itemsCount int) {
	clearConsole()
	if isGen {
		printWarn("The settings in WORLDGEN can not be changed after creation.")
		printInfo(fmt.Sprintf("[%d/%d] WORLDGEN - %s", currentIdx+1, itemsCount, group.Label))
	} else {
		printInfo("The settings in WORLDSETTINGS can be updated anytime.")
		printInfo(fmt.Sprintf("[%d/%d] WORLDSETTINGS - %s", currentIdx+1, itemsCount, group.Label))
	}
	for _, item := range group.Items {
		currentValue := ""
		for _, opt := range item.Options {
			if opt.Data == item.Current {
				currentValue = opt.Text
				break
			}
		}
		fmt.Printf("%-30s: %s\n", item.Label, currentValue)
	}
}

func makeItemLabelList(group *world.WorldConfigGroup) []string {
	items := []string{"Cancel"}
	for _, item := range group.Items {
		items = append(items, item.Label)
	}
	return items
}

func makeOptionLabelList(item *world.WorldConfigItem) []string {
	options := make([]string, 0, len(item.Options))
	for _, opt := range item.Options {
		options = append(options, opt.Text)
	}
	return options
}

func updateGroupConfig(group *world.WorldConfigGroup, target *[]string) {
	labelToValue := func(item *world.WorldConfigItem, label string) string {
		for _, opt := range item.Options {
			if opt.Text == label {
				return opt.Data
			}
		}
		return ""
	}
	for itemIdx, item := range group.Items {
		if !slices.Contains(*target, item.Label) {
			continue
		}

		options := makeOptionLabelList(&item)
		msg := fmt.Sprintf("Select new value for [%s]", item.Label)
		newValue := Selector(options, msg, false)[0]
		newValue = labelToValue(&item, newValue)
		group.Items[itemIdx].Current = newValue
	}
}

func updateMasterGroupConfig(masterGroup *[]world.WorldConfigGroup, isGen bool) {
	groupIdx := 0
	for groupIdx < len(*masterGroup) {
		group := &(*masterGroup)[groupIdx]
		showGroupConfig(group, isGen, groupIdx, len(*masterGroup))

		items := makeItemLabelList(group)
		msg := fmt.Sprintf("Do you want to update values in [%s]?", group.Label)
		selected := Selector(items, msg, true)
		if selected[0] == "Cancel" {
			groupIdx += 1
			continue
		}
		updateGroupConfig(group, &selected)
	}
}

func updateWorldOverride(worldConfig *world.WorldConfig, updateGen bool) {
	if updateGen {
		updateMasterGroupConfig(&worldConfig.WorldGenGroup, true)
	}
	updateMasterGroupConfig(&worldConfig.WorldSettingsGroup, false)
}
