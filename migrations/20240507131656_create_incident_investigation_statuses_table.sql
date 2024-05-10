-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_investigation_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'open',
        CONSTRAINT CHK_incident_investigation_statuses_type CHECK (type IN ('open', 'under-investigation', 'closed', 'resolved','on-hold')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_incident_investigation_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_investigation_statuses;
-- +goose StatementEnd
