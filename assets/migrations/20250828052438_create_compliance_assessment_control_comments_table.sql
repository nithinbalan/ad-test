-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_control_comments (
    id INT IDENTITY(1,1) PRIMARY KEY,
    compliance_assessment_control_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    compliance_assessment_requirement_status_id INT NULL,

    CONSTRAINT FK_compliance_assessment_control_comments_control
        FOREIGN KEY (compliance_assessment_control_id)
        REFERENCES compliance_assessment_controls(id)
        ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT FK_compliance_assessment_control_comments_status
        FOREIGN KEY (compliance_assessment_requirement_status_id)
        REFERENCES compliance_assessment_requirement_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_control_comments_user
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_control_comments;
-- +goose StatementEnd
