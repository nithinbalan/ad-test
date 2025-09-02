-- +goose Up
-- +goose StatementBegin

CREATE TABLE reminder_contexts (
  id INT IDENTITY(1,1) PRIMARY KEY,
  type VARCHAR(255) NOT NULL UNIQUE,
  title_en NVARCHAR(500) NOT NULL,
  title_ar NVARCHAR(500) NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE reminder_contexts;

-- +goose StatementEnd
