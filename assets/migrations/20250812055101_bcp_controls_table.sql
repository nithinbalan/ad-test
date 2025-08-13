-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    control_id INT NOT NULL,
    business_continuity_plan_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_control_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_control_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_control_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_controls;
-- +goose StatementEnd