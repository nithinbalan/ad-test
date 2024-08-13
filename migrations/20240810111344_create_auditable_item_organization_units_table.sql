-- +goose Up
-- +goose StatementBegin
CREATE TABLE auditable_item_organization_units
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    auditable_item_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_auditable_item_organization_units_auditable_items FOREIGN KEY (auditable_item_id) REFERENCES auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_item_organization_units_organization_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_item_organization_units_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE auditable_item_organization_units;

-- +goose StatementEnd