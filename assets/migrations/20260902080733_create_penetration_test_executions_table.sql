-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_executions (
    id INT IDENTITY(1,1) PRIMARY KEY,
    penetration_test_id INT NOT NULL UNIQUE,
    test_start_date DATE NOT NULL,
    test_end_date DATE NOT NULL,
    execution_status_id INT NOT NULL,
    security_result_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT CHK_penetration_test_executions_dates CHECK (test_end_date >= test_start_date),
    CONSTRAINT FK_penetration_test_executions_test FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_executions_status FOREIGN KEY (execution_status_id) REFERENCES penetration_test_execution_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_executions_result FOREIGN KEY (security_result_id) REFERENCES penetration_test_security_results (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_executions_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_executions_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_executions;
-- +goose StatementEnd
