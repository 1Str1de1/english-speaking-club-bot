package services

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"log/slog"
)

type ScheduleStore struct {
	db *sql.DB
}

func NewScheduleStore(username, password, host, port, dbname string, logger *slog.Logger) (*ScheduleStore, error) {
	conn := fmt.Sprintf("postgres://postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, dbname,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("postgres successfully started on URL: ", conn)

	return &ScheduleStore{db: db}, nil
}

func FormatScheduleForTelegram(db *ScheduleStore) string {
	lastText, err := getScheduleText(db)
	if err != nil {
		return fmt.Sprintf("error getting last update %v", err)
	}
	text := "📅 Расписание занятий:\n\n" + lastText
	return text
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
	VALUES ($1, NOW())
`, text)
	return err
}
