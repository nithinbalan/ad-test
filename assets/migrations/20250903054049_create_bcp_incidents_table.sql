-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_incidents
(
    id INT PRIMARY KEY IDENTITY(1, 1),
	business_continuity_plan_id INT NOT NULL,
    incident_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_bcp_incidents_bcps FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_incidents_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_bcp_incidents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
   
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_incidents;

-- +goose StatementEnd