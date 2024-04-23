-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_register (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NULL,
    kri_category_id INT NOT NULL,
    owner_unit_id INT NOT NULL,
    owner_id INT NOT NULL,
    review_frequency_id INT NOT NULL,
    risk_level_id INT NOT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    calculation_type VARCHAR(10) NOT NULL DEFAULT 'automatic',
        CONSTRAINT CHK_kri_register_calculation_type CHECK (calculation_type IN ('automatic', 'manual')),
    low_threshold VARCHAR(455),
    medium_threshold VARCHAR(455),
    high_threshold VARCHAR(455),
    current_value VARCHAR(455),
    last_review_date DATE,
    next_review_date DATE NOT NULL,
    last_review_note NVARCHAR(MAX),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by INT NOT NULL,
    updated_by INT,
    CONSTRAINT FK_kri_register_kri_categories FOREIGN KEY (kri_category_id) REFERENCES kri_categories (id),
    CONSTRAINT FK_kri_register_organization_units FOREIGN KEY (owner_unit_id) REFERENCES organization_units (id),
    CONSTRAINT FK_kri_register_users_1 FOREIGN KEY (owner_id) REFERENCES users (id),
    CONSTRAINT FK_kri_register_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies (id),
    CONSTRAINT FK_kri_register_risk_levels FOREIGN KEY (risk_level_id) REFERENCES risk_levels (id),
    CONSTRAINT FK_kri_register_users_2 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_kri_register_users_3 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kri_register;
-- +goose StatementEnd
