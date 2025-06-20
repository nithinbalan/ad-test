-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_comments (
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_register_id INT NOT NULL,
    risk_timeline_type_id INT NOT NULL,
    comments NVARCHAR(MAX) NOT NULL,
    risk_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME ,
    updated_at DATETIME NULL,

    CONSTRAINT FK_risk_comments_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_comments_risk_timeline_type FOREIGN KEY (risk_timeline_type_id) REFERENCES risk_timeline_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_comments_risk_status FOREIGN KEY (risk_status_id) REFERENCES risk_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_comments_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_comments;
-- +goose StatementEnd
