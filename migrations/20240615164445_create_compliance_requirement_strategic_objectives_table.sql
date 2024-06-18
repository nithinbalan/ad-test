-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_strategic_objectives (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_requirement_id INT NOT NULL,
    strategic_objective_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_requirement_strategic_objectives_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_strategic_objectives_strategic_objectives FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirement_strategic_objectives;
-- +goose StatementEnd
