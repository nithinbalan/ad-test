-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_intiatives (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    created_at DATETIME NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_strategic_intiatives_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_intiatives_responsible_unit FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_intiatives_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_intiatives_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE strategic_intiatives;
-- +goose StatementEnd