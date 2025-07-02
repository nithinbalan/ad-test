-- +goose Up
-- +goose StatementBegin
CREATE TABLE competencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    competency_group_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_competencies_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_competencies_group FOREIGN KEY (competency_group_id) REFERENCES competency_groups(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS competencies;
-- +goose StatementEnd
