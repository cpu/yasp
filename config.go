package yasp

import (
	"errors"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	errBadWidth  = errors.New("invalid WinWidth")
	errBadHeight = errors.New("invalid WinHeight")
)

// Config holds YASP configuration read to/from files on disk.
type Config struct {
	// WinWidth is the YASP window width.
	WinWidth int
	// WinHeight is the YASP window height.
	WinHeight int
}

// valid determines if the Config holds valid values.
func (c Config) Valid() error {
	if c.WinWidth <= 0 {
		return errBadWidth
	}
	if c.WinHeight <= 0 {
		return errBadHeight
	}
	return nil
}

// LoadConfigFile loads the YAML YASP configuration from the provided file path
// and returns it, or an error.
func LoadConfigFile(path string) (*Config, error) {
	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadConfig(confBytes)
}

// LoadConfig loads the YAML YASP configuration from the provided bytes and
// returns it, or an error.
func LoadConfig(yamlBytes []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(yamlBytes, &c); err != nil {
		return nil, err
	}
	if err := c.Valid(); err != nil {
		return nil, err
	}
	return &c, nil
}
