-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_tests (
    id INT IDENTITY(1,1) PRIMARY KEY,
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    system_application_name NVARCHAR(255) NOT NULL,
    penetration_test_type_id INT NOT NULL,
    penetration_test_plan_id INT NOT NULL,
    year INT NOT NULL,
    period VARCHAR(2) NOT NULL,
    penetration_test_classification_id INT NOT NULL,
    penetration_test_executor_id INT NOT NULL,
    cs_responsible_employee_id INT NOT NULL,
    scope NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT CHK_penetration_tests_period CHECK (period IN ('Q1', 'Q2', 'Q3', 'Q4')),
    CONSTRAINT FK_penetration_tests_types FOREIGN KEY (penetration_test_type_id) REFERENCES penetration_test_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_tests_plans FOREIGN KEY (penetration_test_plan_id) REFERENCES penetration_test_plans (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_tests_classifications FOREIGN KEY (penetration_test_classification_id) REFERENCES penetration_test_classifications (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_tests_executors FOREIGN KEY (penetration_test_executor_id) REFERENCES penetration_test_executors (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_tests_responsible_employee FOREIGN KEY (cs_responsible_employee_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_tests_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_tests_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_tests;
-- +goose StatementEnd
