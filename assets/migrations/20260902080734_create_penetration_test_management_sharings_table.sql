-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_management_sharings (
    id INT IDENTITY(1,1) PRIMARY KEY,
    penetration_test_id INT NOT NULL UNIQUE,
    results_shared_with_management BIT NOT NULL DEFAULT 0,
    management_notify_date DATE NULL,
    management_response_date DATE NULL,
    remediation_notes NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_penetration_test_management_sharings_test FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_management_sharings_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_management_sharings_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_management_sharings;
-- +goose StatementEnd
