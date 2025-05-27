-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  id BIGSERIAL PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  capacity INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
