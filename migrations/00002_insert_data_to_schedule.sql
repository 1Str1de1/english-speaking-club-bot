-- +goose Up
INSERT INTO schedule (text)
VALUES ('A2: tentatively Nth of February \nB2: tentatively Nth of February');

-- +goose Down
TRUNCATE TABLE schedule;
