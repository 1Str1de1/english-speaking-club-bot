package server

import (
	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot *tb.BotAPI
}

func NewTgService(token string) (*TelegramService, error) {
	bot, err := tb.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramService{
		bot: bot,
	}, nil
}

func (s *TelegramService) SendHowAreYouPoll(chatId int64) error {
	poll := tb.SendPollConfig{
		BaseChat: tb.BaseChat{
			ChatID: chatId,
		},

		Question: "How are you today?",
		Options: []string{
			"😄 Great",
			"🙂 OK",
			"😕 Not great",
		},

		IsAnonymous:           false,
		AllowsMultipleAnswers: false,
	}

	_, err := s.bot.Send(poll)
	if err != nil {
		return err
	}

	return nil
}
