-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_policy_checklists (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_policy_id INT NOT NULL,
    document_requirement_group_id INT NOT NULL,
    question_group_title NVARCHAR(500) NOT NULL,
    question_title NVARCHAR(500) NOT NULL, 
    justification NVARCHAR(MAX) NULL,
    compliance_status_id INT NULL,
    answered_by_id INT NULL,
    answered_at DATETIME NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_asset_policy_checklists_asset_policies FOREIGN KEY (asset_policy_id) REFERENCES asset_policies(id) ON DELETE NO ACTION ON UPDATE NO ACTION ,
    CONSTRAINT FK_asset_policy_checklists_document_requirement_groups FOREIGN KEY (document_requirement_group_id) REFERENCES document_requirement_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policy_checklists_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policy_checklists_users_1 FOREIGN KEY (answered_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policy_checklists_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policy_checklists_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_policy_checklists;
-- +goose StatementEnd
