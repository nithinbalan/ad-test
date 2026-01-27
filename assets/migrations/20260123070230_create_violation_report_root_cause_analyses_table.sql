-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_root_cause_analyses (
    id INT PRIMARY KEY IDENTITY(1,1),
    violation_report_id INT NOT NULL,
    why INT NOT NULL,
    question NVARCHAR(MAX) NOT NULL,
    answer NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_violation_report_root_cause_analyses FOREIGN KEY (violation_report_id) REFERENCES violation_reportS(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_root_cause_analyses_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_root_cause_analyses_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.violation_report_root_cause_analyses;
-- +goose StatementEnd
