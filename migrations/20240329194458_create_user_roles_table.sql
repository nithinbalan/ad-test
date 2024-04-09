-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_roles (
  id INT PRIMARY KEY IDENTITY(1,1),
  user_id INT NOT NULL,
  role_id INT NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_roles;
-- +goose StatementEnd
