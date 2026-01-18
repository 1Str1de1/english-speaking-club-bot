package server

import (
	"english-speaking-club-bot/internal/config"
	"english-speaking-club-bot/internal/services"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	conf   *config.Config
	logger *slog.Logger
	cron   gocron.Scheduler
	tb     *services.TelegramService
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

	logger.Info("starting postgres...")
	scheduleDb, err := services.NewScheduleStore(
		conf.PostgresConf.Username,
		conf.PostgresConf.Password,
		conf.PostgresConf.Host,
		conf.PostgresConf.Port,
		conf.PostgresConf.DbName,
		logger,
	)

	tb, err := services.NewTgService(conf.Token, conf.YandApiKey, conf.WHAddr, logger, scheduleDb)
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
	//projDir, err := os.Getwd()
	//if err != nil {
	//	return nil, err
	//}
	//
	//logFile := filepath.Join(projDir, "logs", "logs.log")
	//
	//file, err := os.OpenFile(
	//	logFile,
	//	os.O_CREATE|os.O_WRONLY|os.O_APPEND,
	//	0666)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	//tb.SetLogger()

	return logger, nil
}

func (s *Server) Start() error {
	s.logger.Info("starting server...")

	http.HandleFunc("/tg/webhook", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("webhook received")

		update, err := s.tb.Bot.HandleUpdate(r)
		if err != nil {
			s.logger.Error("update error: " + err.Error())
		}

		s.logger.Info(fmt.Sprintf("update: %+v\n", update))

		s.tb.HandleCallback(update)
		s.tb.HandleCommand(update)
		s.tb.HandleMessage(update)
	})

	_, err := s.cron.NewJob(
		gocron.CronJob("0 18 * * *", false),
		gocron.NewTask(func() {
			err := s.tb.SendHowAreYouPoll(s.conf.ChatId)
			if err != nil {
				s.logger.Error("cron or poll error", "err", err)
			}
			s.logger.Info("successfully sent a poll")
		}))

	if err != nil {
		return err
	}

	s.cron.Start()

	s.logger.Info(fmt.Sprintf("server started successfully on port %s", s.conf.Port))
	return http.ListenAndServe(":"+s.conf.Port, nil)
}

func (s *Server) Stop() error {
	s.logger.Info("stopping server...")
	return s.cron.Shutdown()
}
