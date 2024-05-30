-- +goose Up
-- +goose StatementBegin
INSERT INTO organization_units (organization_level_id, title_en, title_ar, status, color_code)
VALUES (1, 'Information Technology Division', 'القسم التقنية المعلومات', 'active', '#1976D2');

DECLARE @ITDivisionID INT = SCOPE_IDENTITY();

INSERT INTO organization_units (organization_level_id, title_en, title_ar, status, color_code)
VALUES (1, 'Human Resources Division', 'قسم الموارد البشرية', 'active', '#388E3C');

DECLARE @HRDivisionID INT = SCOPE_IDENTITY();

-- Inserting departments under the Information Technology Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status, color_code)
VALUES 
(2, @ITDivisionID, 'Software Development Department', 'قسم تطوير البرمجيات', 'active', '#F57C00'),
(2, @ITDivisionID, 'IT Support Department', 'قسم الدعم التقني', 'active', '#7B1FA2');

-- Inserting departments under the Human Resources Division
INSERT INTO organization_units (organization_level_id, organization_unit_id, title_en, title_ar, status, color_code)
VALUES 
(2, @HRDivisionID, 'Recruitment Department', 'قسم التوظيف', 'active', '#E64A19'),
(2, @HRDivisionID, 'Employee Relations Department', 'قسم علاقات الموظفين', 'active', '#5D4037');

COMMIT TRANSACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE organization_levels;
-- +goose StatementEnd
