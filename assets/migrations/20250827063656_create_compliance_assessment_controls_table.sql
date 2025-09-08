-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_controls (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_assessment_id INT NOT NULL,
    compliance_control_id INT NOT NULL,
    completed_by_user_id INT NULL,
    completed_at DATETIME NULL,
    compliance_assessment_framework_value_id INT NULL,
    compliance_assessment_requirement_status_id INT NULL,
    comments NVARCHAR(MAX) NULL,
    is_completed TINYINT NOT NULL DEFAULT 0,
    score INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_compliance_assessment_controls_assessments FOREIGN KEY (compliance_assessment_id) 
        REFERENCES compliance_assessments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_controls_controls FOREIGN KEY (compliance_control_id) 
        REFERENCES compliance_controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_controls_completed_by_user FOREIGN KEY (completed_by_user_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_controls_framework_value FOREIGN KEY (compliance_assessment_framework_value_id) 
        REFERENCES compliance_assessment_framework_values(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_controls_req_status FOREIGN KEY (compliance_assessment_requirement_status_id) 
        REFERENCES compliance_assessment_requirement_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_controls_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_compliance_assessment_controls_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_assessment_controls;
-- +goose StatementEnd
