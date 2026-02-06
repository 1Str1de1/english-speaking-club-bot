package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Token        string
	ChatId       int64
	ThreadPoolId int
	YandApiKey   string
	Port         string
	WHAddr       string
	PostgresConf PostgresConf
	SendPollFlag bool
}

type PostgresConf struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

func NewConfig() (*Config, error) {
	token := os.Getenv("BOT_TOKEN")
	chatIdStr := os.Getenv("CHAT_ID")
	threadPoolIdStr := os.Getenv("MESSAGE_THREAD_ID")
	yandApiKey := os.Getenv("YANDEX_DICT_API_KEY")
	port := os.Getenv("PORT")
	whAddr := os.Getenv("WEBHOOK_ADDRESS")
	sendPollFlagStr := os.Getenv("SEND_POLL")
	pgUsername := os.Getenv("PGUSER")
	pgPassword := os.Getenv("PGPASSWORD")
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	pgDb := os.Getenv("PGDATABASE")

	pgConf := PostgresConf{
		Username: pgUsername,
		Password: pgPassword,
		Host:     pgHost,
		Port:     pgPort,
		DbName:   pgDb,
	}
	if len(token) == 0 {
		return nil, errors.New("error getting bot_token")
	}

	if len(yandApiKey) == 0 {
		return nil, errors.New("error getting yandex_dict_api_key")
	}

	if len(whAddr) == 0 {
		return nil, errors.New("error getting webhook_address")
	}

	sendPollFlag, err := strconv.ParseBool(sendPollFlagStr)
	if err != nil {
		sendPollFlag = false
	}

	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if chatId == 0 || err != nil {
		return nil, errors.New("error getting chat_id")
	}

	threadPoolId, err := strconv.Atoi(threadPoolIdStr)
	if chatId == 0 || err != nil {
		return nil, errors.New(" error getting thread_pool_id: " + err.Error())
	}

	return &Config{
		Token:        token,
		ChatId:       chatId,
		ThreadPoolId: threadPoolId,
		YandApiKey:   yandApiKey,
		Port:         port,
		WHAddr:       whAddr,
		SendPollFlag: sendPollFlag,
		PostgresConf: pgConf,
	}, nil
}
