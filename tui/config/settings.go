package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/tui/views"
	"github.com/jesseduffield/gocui"
	"gopkg.in/yaml.v3"
)

// KeybindingConfig stores the configuration for a single keybinding
type KeybindingConfig struct {
	Action string `yaml:"action"`
	Key    string `yaml:"key"`
	View   string `yaml:"view,omitempty"` // Optional, defaults to global "" (all views)
}

type AppConfig struct {
	Keybindings []KeybindingConfig `yaml:"keybindings"`
}

var specialKeyMap = map[string]gocui.Key{
	"<esc>":        gocui.KeyEsc,
	"<tab>":        gocui.KeyTab,
	"<enter>":      gocui.KeyEnter,
	"<space>":      gocui.KeySpace,
	"<up>":         gocui.KeyArrowUp,
	"<down>":       gocui.KeyArrowDown,
	"<left>":       gocui.KeyArrowLeft,
	"<right>":      gocui.KeyArrowRight,
	"<f1>":         gocui.KeyF1,
	"<f2>":         gocui.KeyF2,
	"<f3>":         gocui.KeyF3,
	"<f4>":         gocui.KeyF4,
	"<f5>":         gocui.KeyF5,
	"<ctrl+a>":     gocui.KeyCtrlA,
	"<ctrl+b>":     gocui.KeyCtrlB,
	"<ctrl+c>":     gocui.KeyCtrlC,
	"<ctrl+d>":     gocui.KeyCtrlD,
	"<ctrl+e>":     gocui.KeyCtrlE,
	"<ctrl+f>":     gocui.KeyCtrlF,
	"<ctrl+g>":     gocui.KeyCtrlG,
	"<ctrl+h>":     gocui.KeyCtrlH, // gocui.KeyBackspace
	"<ctrl+i>":     gocui.KeyCtrlI, // gocui.KeyTab
	"<ctrl+j>":     gocui.KeyCtrlJ,
	"<ctrl+k>":     gocui.KeyCtrlK,
	"<ctrl+l>":     gocui.KeyCtrlL,
	"<ctrl+m>":     gocui.KeyCtrlM,
	"<ctrl+n>":     gocui.KeyCtrlN,
	"<ctrl+o>":     gocui.KeyCtrlO,
	"<ctrl+p>":     gocui.KeyCtrlP,
	"<ctrl+q>":     gocui.KeyCtrlQ,
	"<ctrl+r>":     gocui.KeyCtrlR,
	"<ctrl+s>":     gocui.KeyCtrlS,
	"<ctrl+t>":     gocui.KeyCtrlT,
	"<ctrl+u>":     gocui.KeyCtrlU,
	"<ctrl+v>":     gocui.KeyCtrlV,
	"<ctrl+w>":     gocui.KeyCtrlW,
	"<ctrl+x>":     gocui.KeyCtrlX,
	"<ctrl+y>":     gocui.KeyCtrlY,
	"<ctrl+z>":     gocui.KeyCtrlZ,
	"<ctrl+space>": gocui.KeyCtrlSpace,
	"<ctrl+~>":     gocui.KeyCtrlTilde,
	"<ctrl+2>":     gocui.KeyCtrl2,
	"<ctrl+3>":     gocui.KeyCtrl3,
	"<ctrl+4>":     gocui.KeyCtrl4,
	"<ctrl+5>":     gocui.KeyCtrl5,
	"<ctrl+6>":     gocui.KeyCtrl6,
	"<ctrl+7>":     gocui.KeyCtrl7,
	"<ctrl+8>":     gocui.KeyCtrl8,     // gocui.KeyBackspace2
	"<backspace>":  gocui.KeyBackspace, // Note: gocui.KeyCtrlH is often the same
	"<delete>":     gocui.KeyDelete,
	"<insert>":     gocui.KeyInsert,
	"<home>":       gocui.KeyHome,
	"<end>":        gocui.KeyEnd,
	"<pgup>":       gocui.KeyPgup,
	"<pgdn>":       gocui.KeyPgdn,
}

// ParseKeyString converts a key string (e.g., "<ctrl+c>", "a", "<alt+b>")
// into a gocui.Key (or rune) and gocui.Modifier.
func ParseKeyString(keyStr string) (gocuiKey any, gocuiMod gocui.Modifier, err error) {
	gocuiMod = gocui.ModNone
	lKeyStr := strings.ToLower(keyStr)

	if strings.HasPrefix(lKeyStr, "<alt+") && strings.HasSuffix(lKeyStr, ">") {
		gocuiMod = gocui.ModAlt
		// Extract the key part from original keyStr to preserve case for runes
		innerKeyPart := keyStr[len("<alt+") : len(keyStr)-1]
		lInnerKeyPart := strings.ToLower(innerKeyPart)

		// Check if innerKeyPart is a special key itself, e.g., "<space>"
		// The specialKeyMap keys are all lowercase.
		// We need to re-wrap with <> for map lookup if it's a special key like <space>
		potentialSpecialKey := "<" + lInnerKeyPart + ">"
		if sk, ok := specialKeyMap[potentialSpecialKey]; ok && strings.HasPrefix(lInnerKeyPart, "<") && strings.HasSuffix(lInnerKeyPart, ">") {
			gocuiKey = sk
			return
		} else if len(innerKeyPart) == 1 { // e.g. "a" from "<alt+a>" or "A" from "<alt+A>"
			gocuiKey = rune(innerKeyPart[0]) // Use original case for rune
			return
		} else {
			err = fmt.Errorf("invalid alt key format: %s", keyStr)
			return
		}
	}

	// Handle other special keys (including <ctrl+c>, <tab>, etc.)
	if sk, ok := specialKeyMap[lKeyStr]; ok {
		gocuiKey = sk
		return
	}

	// Handle single characters (runes)
	if len(keyStr) == 1 {
		gocuiKey = rune(keyStr[0])
		return
	}

	err = fmt.Errorf("unknown or malformed key string: %s", keyStr)
	return
}

// LoadKeybindingConfig loads keybindings from a YAML file.
func LoadKeybindingConfig(filePath string) (*AppConfig, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s: %w", filePath, err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data from %s: %w", filePath, err)
	}

	return &config, nil
}

func GetDefaultKeybindings() *AppConfig {
	return &AppConfig{
		Keybindings: []KeybindingConfig{
			// Global
			{Action: "cycle_view", Key: "<tab>", View: ""},
			{Action: "refresh_view", Key: "<f5>", View: ""},
			{Action: "quit", Key: "<ctrl+c>", View: ""},
			{Action: "commandbar_activate", Key: ":", View: ""},

			// SidePanelView
			{Action: "sidepanel_enter", Key: "<enter>", View: views.SidePanelViewName},
			{Action: "sidepanel_up", Key: "<up>", View: views.SidePanelViewName},
			{Action: "sidepanel_down", Key: "<down>", View: views.SidePanelViewName},
			{Action: "sidepanel_up_alt", Key: "k", View: views.SidePanelViewName},
			{Action: "sidepanel_down_alt", Key: "j", View: views.SidePanelViewName},

			// DataTableView
			{Action: "dataview_up", Key: "<up>", View: views.DataTableViewName},
			{Action: "dataview_down", Key: "<down>", View: views.DataTableViewName},
			{Action: "dataview_left", Key: "<left>", View: views.DataTableViewName},
			{Action: "dataview_right", Key: "<right>", View: views.DataTableViewName},
			{Action: "dataview_up_alt", Key: "k", View: views.DataTableViewName},
			{Action: "dataview_down_alt", Key: "j", View: views.DataTableViewName},
			{Action: "dataview_left_alt", Key: "h", View: views.DataTableViewName},
			{Action: "dataview_right_alt", Key: "l", View: views.DataTableViewName},

			// TableCellView (popup)
			{Action: "tablecell_open", Key: "<enter>", View: views.DataTableViewName},
			{Action: "tablecell_close", Key: "<esc>", View: views.TableCellViewName},
			{Action: "tablecell_scroll_up", Key: "<up>", View: views.TableCellViewName},
			{Action: "tablecell_scroll_down", Key: "<down>", View: views.TableCellViewName},

			// CommandBarView
			{Action: "commandbar_close", Key: "<esc>", View: views.CommandBarViewName},
			{Action: "commandbar_enter", Key: "<enter>", View: views.CommandBarViewName},
		},
	}
}
