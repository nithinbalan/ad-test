-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

-- Inserting divisions
INSERT INTO organization_units (organization_level_id, title_en, title_ar, status, color_code)
VALUES (1, 'Information Technology', 'تقنية المعلومات', 'active', '#FF5733');

DECLARE @ITDivisionID INT = SCOPE_IDENTITY();

INSERT INTO organization_units (organization_level_id, title_en, title_ar, status, color_code)
VALUES (1, 'Human Resources', 'الموارد البشرية', 'active', '#33FF57');

DECLARE @HRDivisionID INT = SCOPE_IDENTITY();

-- Inserting departments under the Information Technology Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status, color_code)
VALUES 
(2, @ITDivisionID, 'Software Development', 'تطوير البرمجيات', 'active', '#FFD133'),
(2, @ITDivisionID, 'IT Support', 'الدعم التقني', 'active', '#FFD133');

-- Inserting departments under the Human Resources Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status, color_code)
VALUES 
(2, @HRDivisionID, 'Recruitment', 'التوظيف', 'active', '#337DFF'),
(2, @HRDivisionID, 'Employee Relations', 'علاقات الموظفين', 'active', '#337DFF');

COMMIT TRANSACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_units;
-- +goose StatementEnd
