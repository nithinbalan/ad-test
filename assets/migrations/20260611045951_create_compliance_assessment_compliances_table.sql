-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_compliances (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_register_id INT NOT NULL,
    compliance_assessment_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,

    CONSTRAINT FK_compliance_assessment_compliances_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE CASCADE,
    CONSTRAINT FK_compliance_assessment_compliances_compliance_assessment FOREIGN KEY (compliance_assessment_id) REFERENCES compliance_assessments(id) ON DELETE NO ACTION ON UPDATE CASCADE,
    CONSTRAINT FK_compliance_assessment_compliances_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_compliances;
-- +goose StatementEnd
