-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_requirement_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_assessment_requirement_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_compliance_assessment_requirement_comments_compliance_assessment_requirements FOREIGN KEY (compliance_assessment_requirement_id) REFERENCES compliance_assessment_requirements(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_compliance_assessment_requirement_comments_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_requirement_comments;
-- +goose StatementEnd
