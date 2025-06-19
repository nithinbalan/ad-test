-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_activities (
    id INT PRIMARY KEY IDENTITY(1,1),
    process_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    sequence INT NOT NULL,
    estimated_duration NVARCHAR(100),
    expected_output NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
 
    CONSTRAINT FK_process_activities_process FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_process_activities_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_activities_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS process_activities;
-- +goose StatementEnd
