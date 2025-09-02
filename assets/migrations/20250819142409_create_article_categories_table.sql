-- +goose Up
-- +goose StatementBegin
 
CREATE TABLE article_categories (
  id INT IDENTITY(1,1) PRIMARY KEY,
  title_en NVARCHAR(255) NOT NULL,
  title_ar NVARCHAR(255) NOT NULL,
  status VARCHAR(10) NOT NULL DEFAULT 'active',
  CONSTRAINT CHK_article_categories_status
      CHECK (status IN ('active', 'inactive'))
);
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
DROP TABLE article_categories;
 
-- +goose StatementEnd