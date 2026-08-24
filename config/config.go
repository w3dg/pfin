package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type Config struct {
	configFilePath   string // for internal bookkeeping - Config by default will be read from ~/.config/pfin/config.toml
	DataDir          string `toml:"data_dir"`           // data_dir = "~/.local/share/pfin"
	Currency         string `toml:"currency"`           // default INR
	DateFormat       string `toml:"date_format"`        // YYYY-MM-DD
	DefaultEntryType string `toml:"default_entry_type"` // See EntryType for available types
}

func (c *Config) GetConfigFilePath() string {
	return c.configFilePath
}

func (c *Config) GetDataFilePath() string {
	return path.Join(c.DataDir, "pfin.jsonl")
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	cfg.configFilePath = path
	_, err := toml.DecodeFile(cfg.configFilePath, &cfg)
	return cfg, err
}

func ParseConfigOrWriteDefault() (Config, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not figure out user's home directory, exiting.")
	}

	cfgDir := path.Join(userHome, ".config/pfin")
	cfgFile := "config.toml"
	_, err = os.Stat(path.Join(cfgDir, cfgFile))
	if err != nil { // could not stat the file, so it doesnt exist, write a default out
		err := writeDefaultConfig(cfgDir, cfgFile)
		if err != nil {
			log.Fatal("Attempted to write default config, but failed - ", err)
		}
	}

	cfg, err := LoadConfig(path.Join(cfgDir, cfgFile))
	return cfg, err
}

func writeDefaultConfig(destDir, destFile string) error {
	defaultCfg := `# Where data files live
data_dir = "~/.local/share/pfin"

# Display
currency = "INR"
date_format = "2006-01-02"   # Go's reference layout

# Default entry type when none is specified via --type
default_entry_type = "expense"`

	err := os.MkdirAll(destDir, 0775)
	if err != nil {
		return err
	}
	destPath := path.Join(destDir, destFile)
	err = os.WriteFile(destPath, []byte(defaultCfg), 0644)
	if err != nil {
		return err
	}

	fmt.Println("A default config has been written out to", destPath)
	return nil
}
