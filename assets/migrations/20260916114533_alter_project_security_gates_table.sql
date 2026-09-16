-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_security_gates
ADD
    division_id INT NULL,
    department_id INT NULL,

    CONSTRAINT FK_project_security_gates_division
        FOREIGN KEY (division_id) REFERENCES organization_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gates_department
        FOREIGN KEY (department_id) REFERENCES organization_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_security_gates
DROP CONSTRAINT FK_project_security_gates_division;

ALTER TABLE project_security_gates
DROP CONSTRAINT FK_project_security_gates_department;

ALTER TABLE project_security_gates
DROP COLUMN division_id,
            department_id;
-- +goose StatementEnd
