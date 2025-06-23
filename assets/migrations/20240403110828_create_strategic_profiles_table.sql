-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_profiles (
  id INT PRIMARY KEY IDENTITY(1,1),
  title_en NVARCHAR(255) NOT NULL,
  title_ar NVARCHAR(255),
  address VARCHAR(MAX) NOT NULL,
  location VARCHAR(255),
  email VARCHAR(100),
  industry VARCHAR(100),
  phone VARCHAR(45),
  linkedin VARCHAR(100),
  x VARCHAR(100), 
  facebook VARCHAR(100),
  instagram VARCHAR(100),
  youtube VARCHAR(100),
  about_en VARCHAR(MAX),
  about_ar NVARCHAR(MAX),
  vision_en VARCHAR(500),
  vision_ar NVARCHAR(500),
  mission_en VARCHAR(500),
  mission_ar NVARCHAR(500),
  created_at DATETIME NOT NULL,
  updated_at DATETIME NULL,
  created_by_id INT NOT NULL,
  updated_by_id INT NULL,
  CONSTRAINT FK_strategic_profiles_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_strategic_profiles_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategic_profiles;
-- +goose StatementEnd
