package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	ChatGpt ChatGptConfig `json:"chatgpt" yaml:"chatgpt" mapstructure:"chatgpt"`
}

type ChatGptConfig struct {
	// MattermostURL is the URL of your Mattermost Server instance.
	MattermostURL string `json:"mattermostURL,omitempty" yaml:"mattermostURL,omitempty" mapstructure:"mattermostURL,omitempty"`

	// MattermostToken is your Mattermost access token.
	MattermostToken string `json:"mattermostToken,omitempty" yaml:"mattermostToken,omitempty" mapstructure:"mattermostToken,omitempty"`

	// ChatGPTAPIKey is your ChatGPT API key.
	ChatGPTAPIKey string `json:"chatGPTAPIKey,omitempty" yaml:"chatGPTAPIKey,omitempty" mapstructure:"chatGPTAPIKey,omitempty"`

	// BotUsername is your Mattermost Bot.
	BotUsername string `json:"botUsername,omitempty" yaml:"botUsername,omitempty" mapstructure:"botUsername,omitempty"`
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = os.Getenv(strings.ToUpper(key))
	}

	if len(value) > 0 {
		return value
	}

	if config == nil {
		return ""
	}

	if len(value) > 0 {
		return value
	} else if config.ChatGpt.BotUsername != "" {
		value = config.ChatGpt.BotUsername
	}

	return ""
}
func GetMattermostURL() string {
	mmURL := getEnv("MattermostURL")

	if mmURL != "" {
		return mmURL
	}
	if config == nil {
		return ""
	}
	if mmURL == "" {
		mmURL = config.ChatGpt.MattermostURL
	}
	return mmURL
}

func GetOpenAIAPIKey() string {
	APIKey := getEnv("ChatGPTAPIKey")

	if APIKey != "" {
		return APIKey
	}

	if config == nil {
		return ""
	}

	if APIKey == "" {
		APIKey = config.ChatGpt.ChatGPTAPIKey
	}
	return APIKey
}

func GetMattermostToken() string {
	mmToken := getEnv("MattermostToken")

	if mmToken != "" {
		return mmToken
	}

	if config == nil {
		return ""
	}

	if mmToken == "" {
		mmToken = config.ChatGpt.MattermostToken
	}
	return mmToken
}

func GetBotUsername() string {
	botName := getEnv("BotUsername")

	if botName != "" {
		return botName
	}

	if config == nil {
		return ""
	}

	if botName == "" {
		botName = config.ChatGpt.BotUsername
	}
	return botName
}
