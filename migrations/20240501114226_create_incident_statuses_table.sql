-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'open',
        CONSTRAINT CHK_incident_statuses_type CHECK (type IN ('open', 'investigating', 'in_progress', 'resolved', 'closed', 'pending_review','escalated','re_opened','deferred','mitigated')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_incident_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_statuses;
-- +goose StatementEnd
