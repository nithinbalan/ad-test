-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_environment_preparations (
    id INT IDENTITY(1,1) PRIMARY KEY,
    penetration_test_id INT NOT NULL UNIQUE,
    environment_notify_date DATE NOT NULL,
    expected_ready_date DATE NOT NULL,
    environment_readiness_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_penetration_test_environment_preparations_test FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_environment_preparations_readiness FOREIGN KEY (environment_readiness_id) REFERENCES penetration_test_environment_readinesses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_environment_preparations_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_environment_preparations_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_environment_preparations;
-- +goose StatementEnd
