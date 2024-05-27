-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_test_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    control_test_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_control_test_documents_control_tests FOREIGN KEY (control_test_id) REFERENCES control_tests(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS control_test_documents;
-- +goose StatementEnd
