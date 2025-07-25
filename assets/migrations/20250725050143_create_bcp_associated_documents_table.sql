-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_associated_documents (
    id INT PRIMARY KEY IDENTITY(1, 1),
    business_continuity_plan_id INT NOT NULL,
    document_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_associated_documents_ business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES  business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_associated_documents_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_bcp_associated_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_associated_documents;
-- +goose StatementEnd