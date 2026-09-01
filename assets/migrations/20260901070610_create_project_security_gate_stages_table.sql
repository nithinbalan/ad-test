-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_stages (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    status VARCHAR(10) NOT NULL,

    CONSTRAINT UQ_project_security_gate_stages_type UNIQUE (type),
    CONSTRAINT CHK_project_security_gate_stages_type
        CHECK (type IN ('pre_implementation', 'during_implementation', 'post_implementation')),
    CONSTRAINT CHK_project_security_gate_stages_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_stages;
-- +goose StatementEnd
