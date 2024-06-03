-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

-- Inserting divisions
INSERT INTO organization_units (organization_level_id, title_en, title_ar, status, color_code)
VALUES (1, 'Information Technology Division', 'القسم التقنية المعلومات', 'active', '#FF5733'); -- Example color for Information Technology Division

DECLARE @ITDivisionID INT = SCOPE_IDENTITY();

INSERT INTO organization_units (organization_level_id, title_en, title_ar, status, color_code)
VALUES (1, 'Human Resources Division', 'قسم الموارد البشرية', 'active', '#33FF57'); -- Example color for Human Resources Division

DECLARE @HRDivisionID INT = SCOPE_IDENTITY();

-- Inserting departments under the Information Technology Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status, color_code)
VALUES 
(2, @ITDivisionID, 'Software Development Department', 'قسم تطوير البرمجيات', 'active', '#FFD133'), -- Example color for Software Development Department
(2, @ITDivisionID, 'IT Support Department', 'قسم الدعم التقني', 'active', '#FFD133'); -- Example color for IT Support Department

-- Inserting departments under the Human Resources Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status, color_code)
VALUES 
(2, @HRDivisionID, 'Recruitment Department', 'قسم التوظيف', 'active', '#337DFF'), -- Example color for Recruitment Department
(2, @HRDivisionID, 'Employee Relations Department', 'قسم علاقات الموظفين', 'active', '#337DFF'); -- Example color for Employee Relations Department

COMMIT TRANSACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_units;
-- +goose StatementEnd
