package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

var (
	// Config for bot
	Config *BotConfig
)

// InitBaseConfig initialises the bot's config from file
func InitBaseConfig() (err error) {
	configFilePath := getConfigFilePathFromEnv()
	Config, err = getConfigFromFile(configFilePath)
	if err != nil {
		return err
	}
	err = checkConfig()
	return err
}

func checkConfig() error {
	if Config.Token == "" {
		return errors.New("No token specified in the bot's config")
	}
	return nil
}

// loads the config from a given file
func getConfigFromFile(fileName string) (config *BotConfig, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func getConfigFilePathFromEnv() (configFilePath string) {
	configFilePath = os.Getenv("KingsgateDiscordBotConfigPath")
	if configFilePath == "" {
		configFilePath = "./bot.json"
	}
	return
}
