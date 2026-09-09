-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_closures (
    id INT IDENTITY(1,1) PRIMARY KEY,
    penetration_test_id INT NOT NULL UNIQUE,
    closure_owner_id INT NOT NULL,
    closure_date DATE NOT NULL,
    closure_remarks NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_penetration_test_closures_test FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_closures_owner FOREIGN KEY (closure_owner_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_closures_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_closures_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_closures;
-- +goose StatementEnd
