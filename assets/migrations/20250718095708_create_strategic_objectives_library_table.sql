-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_objective_library (
    id INT PRIMARY KEY IDENTITY(1,1),
    title NVARCHAR(500) NOT NULL, 
    category NVARCHAR(255) NOT NULL
);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS strategic_objective_library;
-- +goose StatementEnd