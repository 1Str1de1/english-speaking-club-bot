-- goose Up
CREATE TABLE schedule (
  id SERIAL PRIMARY KEY,
  text TEXT,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- goose Down
DROP TABLE schedule;