-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercise_test_result_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    test_and_exercise_test_result_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    
    CONSTRAINT FK_test_and_exercise_test_result_documents_test_and_exercise_test_results FOREIGN KEY (test_and_exercise_test_result_id) REFERENCES test_and_exercise_test_results(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_result_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_result_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.test_and_exercise_test_result_documents;
-- +goose StatementEnd