package world

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

func (w *WorldConfig) SaveTo(shardDir string) error {
	return nil
}
