-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  audit_finding_action_plan_comments (
    id INT PRIMARY KEY IDENTITY (1, 1),
    audit_finding_action_plan_id INT NOT NULL,
    comment VARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE ()
    -- Foreign Key Constraints
    CONSTRAINT FK_audit_finding_action_plan_comments_plan FOREIGN KEY (audit_finding_action_plan_id) REFERENCES audit_finding_action_plans (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_audit_finding_action_plan_comments_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
  );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_action_plan_comments;

-- +goose StatementEnd