-- +goose Up
-- +goose StatementBegin
CREATE TABLE third_party_question_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(100) NOT NULL,
    title_ar VARCHAR(100) NOT NULL,
    type VARCHAR(45) NOT NULL,
    CONSTRAINT CHK_third_party_question_types_type CHECK (type IN (
        'yes-no', 
        'dropdown-single', 
        'dropdown-multiple', 
        'rating-scale', 
        'single-text', 
        'text-area', 
        'file-upload', 
        'date-picker'
    )),
    status VARCHAR(45) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_third_party_question_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.third_party_question_types;
-- +goose StatementEnd
