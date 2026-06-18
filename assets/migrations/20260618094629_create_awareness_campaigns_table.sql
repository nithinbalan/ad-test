
-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_campaigns(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(50) NOT NULL UNIQUE,
    awareness_campaign_type_id INT NOT NULL,

    campaign_title NVARCHAR(255) NOT NULL,
    objective NVARCHAR(MAX) NOT NULL,

    awareness_priority_id INT NULL,

    target_audience NVARCHAR(MAX) NOT NULL,
    owner_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME  NULL,
    awareness_campaign_status_id INT NOT NULL,
    linked_sercurity_control NVARCHAR(255),

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_awareness_campaigns_type FOREIGN KEY (awareness_campaign_type_id) REFERENCES awareness_campaign_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_organization_units_1 FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_2 FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_priority FOREIGN KEY (awareness_priority_id) REFERENCES awareness_priorities(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_status FOREIGN KEY (awareness_campaign_status_id) REFERENCES awareness_campaign_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaigns_users_4 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_campaigns;
-- +goose StatementEnd
