package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is the configuration data used by ETL
type Config struct {
	Path     string    `toml:"path"`
	Sheet    string    `toml:"sheet"`
	Columns  ColumnMap `toml:"columns"`
	FirstRow int       `toml:"first_row"`
	NomsPath string    `toml:"noms_path"`
}

// ColumnMap is the map of column names to indices
//
// This permits us not to worry about the headers, and instead simply
// fetch the desired columns directly.
type ColumnMap map[string]int

// each column needs a const string identifier
const (
	AddressS           = "address"
	QtyPurchasedS      = "qty_purchased"
	PurchaseDateS      = "purchase_date"
	UnlockDateS        = "unlock_date"
	NotifyImmediatelyS = "notify_immediately"
)

// DefaultConfig creates a default config struct
func DefaultConfig() (*Config, error) {
	conf := new(Config)
	conf.Columns = make(ColumnMap)
	for _, col := range conf.missingColumns() {
		conf.Columns[col] = 0
	}
	return conf, nil
}

// DefaultConfigPath returns the default path at which a config file is expected
func DefaultConfigPath(ndauhome string) string {
	return "./config.toml"
}

// LoadConfig returns a config object loaded from its file
func LoadConfig(configPath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, os.ErrNotExist
	}
	config := new(Config)
	err = toml.Unmarshal(bytes, config)
	return config, err
}

// LoadDefaultConfig returns a config object loaded from its file
//
// If the file does not exist, a default is transparently created
func LoadDefaultConfig(configPath string) (*Config, error) {
	config, err := LoadConfig(configPath)
	if err != nil && os.IsNotExist(err) {
		config, err = DefaultConfig()
	}
	return config, err
}

// Dump writes the given config object to the specified file
func (conf *Config) Dump(configPath string) error {
	fp, err := os.Create(configPath)
	defer fp.Close()
	if err != nil {
		return err
	}
	return toml.NewEncoder(fp).Encode(conf)
}

// WithConfig wraps configuration editing.
//
// It reads the config file from the specified path, creating a default if
// necessary. It then calls the provided function with that configuration.
// If the provided function returns without error, it writes the new configuration
// to the same path.
func WithConfig(configPath string, lambda func(*Config) error) error {
	config, err := LoadDefaultConfig(configPath)
	if err != nil {
		return err
	}
	err = lambda(config)
	if err != nil {
		return err
	}
	return config.Dump(configPath)
}

func (conf *Config) missingColumns() []string {
	missing := make([]string, 0)
	for _, header := range []string{AddressS, QtyPurchasedS, PurchaseDateS, UnlockDateS, NotifyImmediatelyS} {
		_, ok := conf.Columns[header]
		if !ok {
			missing = append(missing, header)
		}
	}
	return missing
}

// CheckColumns verifies that all expected columns are present
func (conf *Config) CheckColumns() error {
	missing := conf.missingColumns()
	if len(missing) > 0 {
		missingS := strings.Join(missing, ", ")
		return fmt.Errorf("columns missing from config.columns: %s", missingS)
	}
	return nil
}
