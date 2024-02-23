-- +goose Up
-- +goose StatementBegin
-- Inserting divisions
BEGIN TRANSACTION;

-- Inserting divisions
INSERT INTO organization_units (organization_level_id, title_en, title_ar, status)
VALUES (1, 'Information Technology Division', 'القسم التقنية المعلومات', 'active');

DECLARE @ITDivisionID INT = SCOPE_IDENTITY();

INSERT INTO organization_units (organization_level_id, title_en, title_ar, status)
VALUES (1, 'Human Resources Division', 'قسم الموارد البشرية', 'active');

DECLARE @HRDivisionID INT = SCOPE_IDENTITY();

-- Inserting departments under the Information Technology Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status)
VALUES 
(2, @ITDivisionID, 'Software Development Department', 'قسم تطوير البرمجيات', 'active'),
(2, @ITDivisionID, 'IT Support Department', 'قسم الدعم التقني', 'active');

-- Inserting departments under the Human Resources Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status)
VALUES 
(2, @HRDivisionID, 'Recruitment Department', 'قسم التوظيف', 'active'),
(2, @HRDivisionID, 'Employee Relations Department', 'قسم علاقات الموظفين', 'active');

COMMIT TRANSACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_levels;
-- +goose StatementEnd
