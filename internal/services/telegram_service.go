package services

import (
	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	Bot        *tb.BotAPI
	yandApiKey string
}

func NewTgService(token, yandApiKey string) (*TelegramService, error) {
	bot, err := tb.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramService{
		Bot:        bot,
		yandApiKey: yandApiKey,
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

	_, err := s.Bot.Send(poll)
	if err != nil {
		return err
	}

	return nil
}

func (s *TelegramService) HandleCommand(update tb.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	switch update.Message.Command() {
	case "start":

	case "randomword":
		s.handleRandomWord(update)
	}
}

func (s *TelegramService) handleRandomWord(update tb.Update) {
	word, err := ExecuteRandomWordCommand(s.yandApiKey)
	if err != nil {
		msg := tb.NewMessage(update.Message.Chat.ID,
			"❌ Не могу получить случайное слово")
		s.Bot.Send(msg)
		return
	}

	msg := tb.NewMessage(update.Message.Chat.ID, word)
	msg.ParseMode = "Markdown"
	s.Bot.Send(msg)
}
