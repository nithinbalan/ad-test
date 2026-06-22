-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_assessment_themes
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,

    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_awareness_assessment_themes_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_assessment_themes;
-- +goose StatementEnd

