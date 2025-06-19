-- +goose Up
-- +goose StatementBegin
ALTER TABLE stakeholders
ADD type_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_process_dependency_types
FOREIGN KEY (type_id) REFERENCES process_dependency_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE stakeholders
ADD stakeholder_category_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_stakeholder_categories
FOREIGN KEY (stakeholder_category_id) REFERENCES stakeholder_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE stakeholders
ADD role NVARCHAR(100) NULL;

ALTER TABLE stakeholders
ADD influence_level_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_severity_levels1
FOREIGN KEY (influence_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE stakeholders
ADD impact_level_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_severity_levels2
FOREIGN KEY (impact_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE stakeholders
ADD contact_information NVARCHAR(MAX) NULL;

ALTER TABLE stakeholders
ADD organization_unit_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_organization_units
FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE stakeholders
ADD third_party_register_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_third_party_register
FOREIGN KEY (third_party_register_id) REFERENCES third_party_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE stakeholders
ADD client_id INT NULL;


ALTER TABLE stakeholders
ADD user_id INT NULL;

ALTER TABLE stakeholders
ADD CONSTRAINT FK_stakeholders_user3
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.stakeholder_influence_levels;
-- +goose StatementEnd

