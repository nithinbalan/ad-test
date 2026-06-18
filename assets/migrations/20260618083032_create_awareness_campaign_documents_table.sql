-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_campaign_documents (
    id INT PRIMARY KEY IDENTITY(1,1),

    awareness_campaign_type_id INT NOT NULL,
    awareness_campaign_detail_id INT NULL,

    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,

    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_awareness_campaign_documents_type FOREIGN KEY (awareness_campaign_type_id)
        REFERENCES awareness_campaign_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_documents_detail FOREIGN KEY (awareness_campaign_detail_id)
        REFERENCES awareness_campaign_details(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_documents_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_documents_users_2 FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_campaign_documents;
-- +goose StatementEnd
