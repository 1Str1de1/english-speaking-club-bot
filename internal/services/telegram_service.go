package services

import (
	"errors"
	"fmt"
	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type TelegramService struct {
	Bot        *tb.BotAPI
	yandApiKey string
	logger     *slog.Logger
}

func NewTgService(token, yandApiKey string, logger *slog.Logger) (*TelegramService, error) {
	bot, err := tb.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	wh, err := tb.NewWebhook("https://poetic-warmth-production-7bfc.up.railway.app/tg/webhook")
	if err != nil {
		return nil, errors.New("error setting webhook " + err.Error())
	}

	_, err = bot.Request(wh)
	if err != nil {
		return nil, errors.New("error requesting webhook " + err.Error())
	}

	return &TelegramService{
		Bot:        bot,
		yandApiKey: yandApiKey,
		logger:     logger,
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

func (s *TelegramService) HandleCommand(update *tb.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	switch update.Message.Command() {
	case "start":

	case "randomword":
		s.handleRandomWord(update)
	}
}

func (s *TelegramService) handleRandomWord(update *tb.Update) {
	s.logger.Info("handling randomword update...")

	word, err := ExecuteRandomWordCommand(s.yandApiKey)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error executing randomword commmand: " + err.Error()))
		msg := tb.NewMessage(update.Message.Chat.ID,
			"❌ Не могу получить случайное слово")
		s.Bot.Send(msg)
		return
	}

	s.logger.Info(fmt.Sprintf("formatted word is: %s", word))

	msg := tb.NewMessage(update.Message.Chat.ID, word)
	msg.ParseMode = "Markdown"
	_, err = s.Bot.Send(msg)
	if err != nil {
		s.logger.Error("message not sent, error: " + err.Error())
	} else {
		s.logger.Info(fmt.Sprintf("message %+v sent successfully", msg))
	}
}
