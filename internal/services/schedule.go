package services

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"time"
)

type ScheduleStore struct {
	db *sql.DB
}

func NewScheduleStore(username, password, host, port, dbname string, logger *slog.Logger) (*ScheduleStore, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, dbname,
	)

	db, err := sql.Open("pgx", conn)
	if err != nil {
		logger.Error("error opening postgres ", err.Error(), conn)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logger.Error("error pinging postgres ", err.Error(), conn)
		return nil, err
	}

	logger.Info("postgres successfully started on URL: ", conn)

	return &ScheduleStore{db: db}, nil
}

func FormatScheduleForTelegram(db *ScheduleStore) (string, error) {
	lastText, err := getScheduleText(db)
	if err != nil {
		return "", errors.New("error getting last update" + err.Error())
	}
	text := "📅 Расписание занятий:\n\n" + lastText
	return text, nil
}

func getScheduleText(db *ScheduleStore) (string, error) {
	var text string
	err := db.db.QueryRow(`
	SELECT (text)
	FROM schedule
	ORDER BY last_update DESC
	LIMIT 1
`).Scan(&text)
	return text, err
}

func SaveSchedule(db *ScheduleStore, text string) error {
	_, err := db.db.Exec(`
	INSERT INTO schedule (text, last_update)
	VALUES ($1, $2)
`, text, time.Now())
	return err
}
