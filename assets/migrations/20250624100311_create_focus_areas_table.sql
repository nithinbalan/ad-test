-- +goose Up
-- +goose StatementBegin
CREATE TABLE focus_areas (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_focus_areas_status CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_focus_area_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_focus_area_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE focus_areas;
-- +goose StatementEnd
