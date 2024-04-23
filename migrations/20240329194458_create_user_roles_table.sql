-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_roles (
  id INT PRIMARY KEY IDENTITY(1,1),
  user_id INT NOT NULL,
  role_id INT NOT NULL,
  CONSTRAINT FK_user_roles_users FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT FK_user_roles_roles FOREIGN KEY (role_id) REFERENCES roles(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.user_roles;
-- +goose StatementEnd
