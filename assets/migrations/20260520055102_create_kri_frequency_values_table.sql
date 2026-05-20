-- +goose Up
-- +goose StatementBegin

CREATE TABLE kri_frequency_values (
    id INT PRIMARY KEY IDENTITY(1, 1),
    kri_register_id INT NOT NULL,
    year INT NOT NULL,
    month_en VARCHAR(1000) NOT NULL,
    month_ar VARCHAR(1000) NOT NULL,
    type VARCHAR(1000) NOT NULL, 
    target_value DECIMAL(18,2) NOT NULL,
    actual_value DECIMAL(18,2) NULL,
    achievement_value DECIMAL(18,2) NULL,
    risk_tolerance_status_id INT NULL,
    is_achieved BIT NOT NULL DEFAULT 0,
    is_missed BIT NOT NULL DEFAULT 0,
    created_by_id INT,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_kri_frequency_values_kri_register FOREIGN KEY (kri_register_id) REFERENCES kri_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kri_frequency_values_kri_risk_tolerance_statuses FOREIGN KEY (risk_tolerance_status_id) REFERENCES kri_risk_tolerance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kri_frequency_values_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
    
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE kri_frequency_values;
-- +goose StatementEnd
