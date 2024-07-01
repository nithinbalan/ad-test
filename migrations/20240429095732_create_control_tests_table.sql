-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_tests (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_register_id INT NULL,
    control_id INT NOT NULL,
    control_implementation_status_id INT NOT NULL,
    control_efficiency_status_id INT NOT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    score INT,
    comments NVARCHAR(MAX),
    tested_source VARCHAR(45) NOT NULL,
    tested_at DATE,
    tested_by INT NOT NULL,
    updated_at DATETIME,
    updated_by_id INT,
    CONSTRAINT FK_control_tests_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_tests_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_tests_control_implementation_statuses FOREIGN KEY (control_implementation_status_id) REFERENCES control_implementation_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_tests_control_efficiency_statuses FOREIGN KEY (control_efficiency_status_id) REFERENCES control_efficiency_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_tests_users_1 FOREIGN KEY (tested_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_control_tests_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS control_tests;
-- +goose StatementEnd
