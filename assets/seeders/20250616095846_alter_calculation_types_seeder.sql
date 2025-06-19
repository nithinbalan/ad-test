-- +goose Up
-- +goose StatementBegin
DELETE FROM calculation_types
WHERE type = 'max';

DELETE FROM calculation_types
WHERE type = 'min';

DELETE FROM calculation_types
WHERE type = 'count';

DELETE FROM calculation_types
WHERE type = 'custom-formula';

UPDATE calculation_types
SET type = 'formula', title_en='Formula Based', title_ar='على أساس الصيغة'
WHERE type = 'average';

UPDATE calculation_types
SET type = 'direct', title_en='Direct', title_ar='مباشر'
WHERE type = 'sum';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE calculation_types
SET type = 'sum', title_en='Sum', title_ar='مجموع'
WHERE type = 'direct';

UPDATE calculation_types
SET type = 'average', title_en='Average', title_ar='متوسط'
WHERE type = 'formula';

INSERT INTO calculation_types (title_en, title_ar,type, status)
VALUES
    ('Max', 'الأعلى', 'max', 'active');

    INSERT INTO calculation_types (title_en, title_ar,type, status)
VALUES
    ('Min', 'دقيقة', 'min', 'active');

    INSERT INTO calculation_types (title_en, title_ar,type, status)
VALUES
    ('Count', 'عدد', 'count', 'active');

    INSERT INTO calculation_types (title_en, title_ar,type, status)
VALUES
    ('Custom Formula', 'صيغة مخصصة', 'custom-formula', 'active');
-- +goose StatementEnd
