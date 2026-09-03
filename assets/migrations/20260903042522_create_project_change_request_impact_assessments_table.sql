-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_change_request_impact_assessments (
    id INT PRIMARY KEY IDENTITY(1, 1),
    project_change_request_id INT NOT NULL,
    project_impact_assessment_type_id INT NOT NULL,
    [current_date] DATE NULL,
    [revised_date] DATE NULL,
    current_budget DECIMAL(18, 2) NULL,
    revised_budget DECIMAL(18, 2) NULL,
    current_resource_count INT NULL,
    revised_resource_count INT NULL,
    technical_complexity NVARCHAR(MAX) NULL,

    CONSTRAINT FK_project_change_request_impact_assessments_change_requests FOREIGN KEY (project_change_request_id)
        REFERENCES project_change_requests(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_request_impact_assessments_types FOREIGN KEY (project_impact_assessment_type_id)
        REFERENCES project_impact_assessment_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_change_request_impact_assessments;
-- +goose StatementEnd
