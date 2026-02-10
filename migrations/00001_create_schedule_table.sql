-- +goose Up
CREATE TABLE schedule (
  id SERIAL PRIMARY KEY,
  text TEXT,
  last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE schedule;
