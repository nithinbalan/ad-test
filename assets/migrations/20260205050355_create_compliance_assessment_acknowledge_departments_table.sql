-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_acknowledge_departments 
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_assessment_id INT NOT NULL,
    user_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    department_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_assessment_acknowledge_departments_compliance FOREIGN KEY (compliance_assessment_id) REFERENCES compliance_assessments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_acknowledge_departments_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_acknowledge_departments_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_acknowledge_departments_departments FOREIGN KEY (department_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_acknowledge_departments;
-- +goose StatementEnd