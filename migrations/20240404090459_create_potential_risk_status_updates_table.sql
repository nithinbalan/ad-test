-- +goose Up
-- +goose StatementBegin
CREATE TABLE potential_risk_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    potential_risk_id INT NOT NULL,
    potential_risk_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL ,
    created_by INT NOT NULL,
    CONSTRAINT FK_potential_risk_status_updates_potential_risk FOREIGN KEY (potential_risk_id) REFERENCES potential_risks(id),
    CONSTRAINT FK_potential_risk_status_updates_potential_risk_statuses FOREIGN KEY (potential_risk_status_id) REFERENCES potential_risk_statuses(id),
    CONSTRAINT FK_potential_risk_status_updates_users FOREIGN KEY (created_by) REFERENCES users(id),
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE potential_risk_status_updates;
-- +goose StatementEnd
