-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_retests (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL,
    penetration_test_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    penetration_test_executor_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT UQ_penetration_test_retests_reference_number UNIQUE (reference_number),
    CONSTRAINT CHK_penetration_test_retests_status CHECK (status IN ('in-progress', 'completed')),
    CONSTRAINT FK_penetration_test_retests_penetration_tests FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retests_penetration_test_executors FOREIGN KEY (penetration_test_executor_id) REFERENCES penetration_test_executors (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retests_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retests_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_retests;
-- +goose StatementEnd
