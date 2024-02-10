package world

import (
	"fmt"
	"os"
	"strings"
)

type worldgenOverride struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Desc      string                 `json:"desc"`
	Location  string                 `json:"location"`
	PlayStyle string                 `json:"playstyle"`
	Version   int                    `json:"version"`
	Overrides map[string]interface{} `json:"overrides"`
}

type WorldConfig struct {
	Version            string             `json:"version"`
	Language           string             `json:"language"`
	Location           string             `json:"location"`
	IsMaster           bool               `json:"is_master"`
	WorldGenGroup      []worldConfigGroup `json:"worldgen_group"`
	WorldSettingsGroup []worldConfigGroup `json:"worldsettings_group"`
}

type worldConfigGroup struct {
	Name  string            `json:"name"`
	Label string            `json:"label"`
	Items []worldConfigItem `json:"items"`
}

type worldConfigItem struct {
	Name       string                     `json:"name"`
	Label      string                     `json:"label"`
	Default    string                     `json:"default"`
	Current    string                     `json:"-"`
	Options    []worldConfigItemOption    `json:"options"`
	WidgetType string                     `json:"widget_type"`
	Image      string                     `json:"image"`
	Atlas      string                     `json:"atlas"`
	OptsRemap  worldConfigItemOptionRemap `json:"options_remap"`
}

type worldConfigItemOption struct {
	Data string `json:"data"`
	Text string `json:"text"`
}

type worldConfigItemOptionRemap struct {
	Image string `json:"img"`
	Alias string `json:"alias"`
}

func (w *WorldConfig) setAllCurrentDefault() {
	for i := 0; i < len(w.WorldGenGroup); i++ {
		for j := 0; j < len(w.WorldGenGroup[i].Items); j++ {
			w.WorldGenGroup[i].Items[j].Current = w.WorldGenGroup[i].Items[j].Default
		}
	}
	for i := 0; i < len(w.WorldSettingsGroup); i++ {
		for j := 0; j < len(w.WorldSettingsGroup[i].Items); j++ {
			w.WorldSettingsGroup[i].Items[j].Current = w.WorldSettingsGroup[i].Items[j].Default
		}
	}
}

func (w *WorldConfig) SaveTo(shardDirPath string) error {
	overrideFilePath := shardDirPath + "/worldgenoverride.lua"

	file, err := os.Create(overrideFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	builder := strings.Builder{}
	builder.WriteString("return {\n")
	if w.Location == "forest" {
		builder.WriteString("  preset = \"SURVIVAL_TOGETHER\",\n")
	} else {
		builder.WriteString("  preset = \"DST_CAVE\",\n")
	}
	builder.WriteString("  override_enabled=true,\n")
	builder.WriteString("  overrides={\n")
	for _, group := range w.WorldGenGroup {
		for _, item := range group.Items {
			builder.WriteString(fmt.Sprintf("    %s=%s,\n", item.Name, item.Current))
		}
	}
	for _, group := range w.WorldSettingsGroup {
		for _, item := range group.Items {
			builder.WriteString(fmt.Sprintf("    %s=%s,\n", item.Name, item.Current))
		}
	}
	builder.WriteString("  }\n")
	builder.WriteString("}\n")

	_, err = file.WriteString(builder.String())
	return err
}
