-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_report_kris (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_report_id INT NOT NULL,
    risk_report_identification_id INT NOT NULL,

    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    low_threshold INT NULL,
    medium_threshold INT NULL,
    current_value INT NULL,

    unit_title_en NVARCHAR(255) NULL,
    unit_title_ar NVARCHAR(255) NULL,

    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_risk_report_kris_risk_reports FOREIGN KEY (risk_report_id) REFERENCES risk_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_report_kris_risk_report_identifications FOREIGN KEY (risk_report_identification_id) REFERENCES risk_report_identifications (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_report_kris;
-- +goose StatementEnd
