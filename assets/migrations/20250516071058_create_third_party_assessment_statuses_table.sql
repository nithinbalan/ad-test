-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_assessment_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(15) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_third_party_assessment_statuses_type CHECK (type IN ('draft','new','in-progress', 'submitted')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_third_party_assessment_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.third_party_assessment_statuses;
-- +goose StatementEnd
