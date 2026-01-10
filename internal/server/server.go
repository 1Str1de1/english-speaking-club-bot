package server

import (
	"english-speaking-club-bot/internal/config"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"os"
	"path/filepath"
)

type Server struct {
	conf   *config.Config
	logger *slog.Logger
	cron   gocron.Scheduler
	tb     *TelegramService
}

func NewServer(conf *config.Config) *Server {
	logger, err := setupLogger()
	if err != nil {
		panic("logger error" + err.Error())
	}

	cron, err := gocron.NewScheduler()
	if err != nil {
		panic("error starting cron " + err.Error())
	}

	tb, err := NewTgService(conf.Token)
	if err != nil {
		panic("error starting tg service" + err.Error())
	}

	return &Server{
		conf:   conf,
		logger: logger,
		cron:   cron,
		tb:     tb,
	}
}

func setupLogger() (*slog.Logger, error) {
	projDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	logFile := filepath.Join(projDir, "logs", "logs.log")

	file, err := os.OpenFile(
		logFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	//tb.SetLogger()

	return logger, nil
}

func (s *Server) Start() error {
	s.logger.Info("starting server...")

	_, err := s.cron.NewJob(
		gocron.CronJob("0 18 * * *", false),
		gocron.NewTask(func() {
			err := s.tb.SendHowAreYouPoll(s.conf.ChatId)
			if err != nil {
				s.logger.Error("cron or poll error", "err", err)
			}
		}))

	if err != nil {
		return err
	}

	s.cron.Start()

	s.logger.Info("server started successfully")
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping server...")
	return s.cron.Shutdown()
}
