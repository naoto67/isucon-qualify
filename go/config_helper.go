package main

import "errors"

//// config cache /////
var configCache map[string]Config

func initializeConfig() error {
	configCache = make(map[string]Config)
	configs := []Config{}
	err := dbx.Select(&configs, "SELECT * FROM `configs`")
	if err != nil {
		return err
	}

	for _, config := range configs {
		configCache[config.Name] = config
	}

	return nil
}

func getConfigByName(name string) (string, error) {
	if _, ok := configCache[name]; !ok {
		return "", errors.New("config not found")
	}
	return configCache[name].Val, nil
}
