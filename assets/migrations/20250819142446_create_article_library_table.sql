-- +goose Up
-- +goose StatementBegin
 
CREATE TABLE article_library (
  id INT IDENTITY(1,1) PRIMARY KEY,
  article_category_id INT NOT NULL,
  image_path NVARCHAR(500) NOT NULL,
  title_en NVARCHAR(MAX) NOT NULL,
  title_ar NVARCHAR(MAX) NOT NULL,
  table_of_contents NVARCHAR(MAX) NULL, -- JSON array
  content NVARCHAR(MAX) NOT NULL, -- HTML content
  created_by NVARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  CONSTRAINT FK_article_library_article_categories
    FOREIGN KEY (article_category_id)
    REFERENCES article_categories(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
);
 
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
DROP TABLE article_library;
 
-- +goose StatementEnd