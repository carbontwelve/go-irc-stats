package ircstats

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Location          string
	SaveLocation      string
	DatabaseLocation  string
	PageTitle         string
	PageDescription   string
	HeatMapInterval   uint
	ActivityPeriod    uint
	Ignore            []string
	NickNameMapping   map[string][]string
	NickNameHashTable map[string]string
	Profiles          map[string]map[string]string
}

func (c *Config) Load(path string) (err error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	if c.Location == "" {
		return errors.New("The full path to your log file must be provided")
	}

	c.Location, _ = filepath.Abs(c.Location)

	if _, err := os.Stat(c.Location); os.IsNotExist(err) {
		return errors.New("The path [" + c.Location + "] could not be read, does it exist?")
	}

	if c.SaveLocation == "" {
		return errors.New("The full path to your save location must be provided")
	}

	c.SaveLocation, _ = filepath.Abs(c.SaveLocation)

	if _, err := os.Stat(c.SaveLocation); os.IsNotExist(err) {
		return errors.New("The path [" + c.SaveLocation + "] could not be read, does it exist?")
	}

	if c.DatabaseLocation == "" {
		c.DatabaseLocation = "./db.bin"
	}

	c.DatabaseLocation, _ = filepath.Abs(c.DatabaseLocation)

	if c.HeatMapInterval == 0 {
		c.HeatMapInterval = 100
	}

	// Create a hash table for username corrections
	c.NickNameHashTable = make(map[string]string)
	for u, mappings := range c.NickNameMapping {
		for _, mappedNick := range mappings {
			c.NickNameHashTable[mappedNick] = u
		}
	}

	return
}
