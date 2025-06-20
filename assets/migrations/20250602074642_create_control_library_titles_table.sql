-- +goose Up
-- +goose StatementBegin

CREATE TABLE control_library_titles (
    id INT PRIMARY KEY IDENTITY(1,1),
    control_library_id INT NOT NULL,
    compliance_library_id INT NOT NULL,
    control_code NVARCHAR(255),
    control_title_en NVARCHAR(MAX),
    control_title_ar NVARCHAR(MAX),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    FOREIGN KEY (control_library_id) REFERENCES control_library(id),
    FOREIGN KEY (compliance_library_id) REFERENCES compliance_libraries(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE control_library_titles;
-- +goose StatementEnd