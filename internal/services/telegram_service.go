package services

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	tb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	Bot                 *tb.BotAPI
	yandApiKey          string
	logger              *slog.Logger
	db                  *ScheduleStore
	waitingForSchedule  map[int64]bool
	waitingForTranslate map[int64]bool
}

func NewTgService(token, yandApiKey, WHAddr string, logger *slog.Logger, db *ScheduleStore) (*TelegramService, error) {
	bot, err := tb.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	wh, err := tb.NewWebhook(fmt.Sprintf("https://%s/tg/webhook", WHAddr))
	logger.Info(fmt.Sprintf("webhook address = %s", wh.URL))
	if err != nil {
		return nil, errors.New("error setting webhook " + err.Error())
	}

	_, err = bot.Request(wh)
	if err != nil {
		return nil, errors.New("error requesting webhook " + err.Error())
	}

	return &TelegramService{
		Bot:                 bot,
		yandApiKey:          yandApiKey,
		logger:              logger,
		db:                  db,
		waitingForSchedule:  make(map[int64]bool),
		waitingForTranslate: make(map[int64]bool),
	}, nil
}

func (s *TelegramService) SendHowAreYouPoll(chatId int64) error {
	questions := []string{
		"What made you smile today, and why?",
		"If you could change one small thing about your daily routine, what would it be?",
		"What’s something you’ve learned recently that surprised you?",
		"What’s a song, video, or quote that’s been on your mind lately? What do you think it means?",
		"What’s one decision you made recently? How do you feel about it now?",
		"What food do you like, and how often do you eat it?",
		"What is one place you want to visit, and why?",
		"What did you do last weekend? Tell us 2–3 things.",
		"What kind of music do you like, and when do you listen to it?",
		"What do you usually do after work or school?",
	}

	num := rand.Intn(10)

	poll := tb.SendPollConfig{
		BaseChat: tb.BaseChat{
			ChatID: chatId,
		},

		Question: questions[num],
		Options: []string{
			"I answered",
			"I'm gay",
			"Some love for developer",
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

	case "schedule":
		s.handleSchedule(update)

	case "translate":
		s.handleTranslate(update)
	}

}

func (s *TelegramService) HandleCallback(update *tb.Update) {
	if update.CallbackQuery == nil {
		return
	}

	cb := update.CallbackQuery
	s.logger.Info("callback received", "data", cb.Data)

	switch cb.Data {
	case "edit_schedule":
		s.handleEditSchedule(cb)
	case "ok":
		s.handleOk(cb)
	case "decline":
		s.handleDecline(cb)
	default:
		s.logger.Info(fmt.Sprintf("callback is: %v", cb.Data))
	}
}

func (s *TelegramService) HandleMessage(update *tb.Update) {
	if update.Message == nil || update.Message.IsCommand() {
		return
	}

	chatID := update.Message.Chat.ID

	if s.waitingForSchedule[chatID] {
		s.logger.Info("schedule update received", "chat", chatID)

		if err := SaveSchedule(s.db, update.Message.Text); err != nil {
			s.logger.Error("error saving new schedule: ", "err", err)
			s.Bot.Send(tb.NewMessage(chatID, "❌ Ошибка при сохранении расписания"))
			s.waitingForSchedule[chatID] = false
			return
		}

		s.waitingForSchedule[chatID] = false

		msg := tb.NewMessage(chatID, "✅ Расписание обновлено!")
		s.Bot.Send(msg)
	}

	if s.waitingForTranslate[chatID] {
		s.logger.Info("translate update received", "chat", chatID)
		text, err := ExecuteTranslateText(s.yandApiKey, update.Message.Text)
		if err != nil {
			s.logger.Error("error getting translation: ", "err", err)
			s.Bot.Send(tb.NewMessage(chatID, "❌ Error getting translation"))
			s.waitingForTranslate[chatID] = false
			return
		}

		s.waitingForTranslate[chatID] = false

		msg := tb.NewMessage(chatID, text)
		s.Bot.Send(msg)
	}

}

func (s *TelegramService) handleRandomWord(update *tb.Update) {
	s.logger.Info("handling randomword update...")

	word, err := ExecuteRandomWordCommand(s.yandApiKey)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error executing randomword command: " + err.Error()))
		msg := tb.NewMessage(update.Message.Chat.ID,
			"❌ Не могу получить случайное слово")
		s.Bot.Send(msg)
		return
	}

	s.logger.Info(fmt.Sprintf("formatted word is: %s", word))

	msg := tb.NewMessage(update.Message.Chat.ID, word)
	//msg.ParseMode = "Markdown"
	_, err = s.Bot.Send(msg)
	if err != nil {
		s.logger.Error("message not sent, error: " + err.Error())
	} else {
		s.logger.Info(fmt.Sprintf("message %+v sent successfully", msg))
	}
}

func (s *TelegramService) handleSchedule(update *tb.Update) {
	s.logger.Info("handling schedule command...")

	text, err := FormatScheduleForTelegram(s.db)
	if err != nil {
		s.logger.Error("error formatting schedule: ", "err", err)
	}

	btnSch := tb.NewInlineKeyboardButtonData("✏️ Изменить расписание", "edit_schedule")
	btnOk := tb.NewInlineKeyboardButtonData("✅ OK", "ok")
	keyboard := tb.NewInlineKeyboardMarkup(
		tb.NewInlineKeyboardRow(btnSch),
		tb.NewInlineKeyboardRow(btnOk),
	)

	msg := tb.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	s.Bot.Send(msg)

}

func (s *TelegramService) handleEditSchedule(cb *tb.CallbackQuery) {
	edit := tb.NewEditMessageText(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		"✏️ Введите новое расписание одним сообщением:",
	)

	s.waitingForSchedule[cb.Message.Chat.ID] = true

	if _, err := s.Bot.Send(edit); err != nil {
		s.logger.Error("error editing message ", "err", err)
	}

	answer := tb.NewCallback(cb.ID, "")
	if _, err := s.Bot.Request(answer); err != nil {
		s.logger.Error("error requesting answer: ", "err", err)
	}
}

func (s *TelegramService) handleOk(cb *tb.CallbackQuery) {

	text, err := FormatScheduleForTelegram(s.db)
	if err != nil {
		s.logger.Error("error formatting schedule: ", "err", err)
	}

	edit := tb.NewEditMessageText(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		text,
	)

	s.waitingForSchedule[cb.Message.Chat.ID] = false

	if _, err := s.Bot.Send(edit); err != nil {
		s.logger.Error("error editing message ", "err", err)
	}

}

func (s *TelegramService) handleTranslate(update *tb.Update) {
	s.logger.Info("handling translate command...")

	text := "✏️Enter the word you want to translate"

	btnDecl := tb.NewInlineKeyboardButtonData("🤞 Don't translate", "decline")
	btnTran := tb.NewInlineKeyboardButtonData("🌍 Translate", "translate cb")
	keyboard := tb.NewInlineKeyboardMarkup(
		tb.NewInlineKeyboardRow(btnDecl),
		tb.NewInlineKeyboardRow(btnTran),
	)

	msg := tb.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	if _, err := s.Bot.Send(msg); err != nil {
		s.logger.Error("error sending translate start message " + err.Error())
	}

}

func (s *TelegramService) handleTranslateCb(cb *tb.CallbackQuery) {
	edit := tb.NewEditMessageText(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		"Enter the word or phrase to translate:",
	)

	s.waitingForTranslate[cb.Message.Chat.ID] = true

	if _, err := s.Bot.Send(edit); err != nil {
		s.logger.Error("error editing message ", "err", err)
	}

	answer := tb.NewCallback(cb.ID, "")
	if _, err := s.Bot.Request(answer); err != nil {
		s.logger.Error("error requesting answer: ", "err", err)
	}
}

func (s *TelegramService) handleDecline(cb *tb.CallbackQuery) {
	text := "🧦Translation declined"

	edit := tb.NewEditMessageText(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		text,
	)

	s.waitingForTranslate[cb.Message.Chat.ID] = false

	if _, err := s.Bot.Send(edit); err != nil {
		s.logger.Error("error editing message ", "err", err)
	}
}
