-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_strategic_objectives  (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_id INT NOT NULL,
    audit_finding_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_audit_finding_strategic_objectives_strategic_objectives FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_strategic_objectives_audit_findings FOREIGN KEY (audit_finding_id) REFERENCES audit_findings(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_strategic_objectives_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_strategic_objectives;
-- +goose StatementEnd
