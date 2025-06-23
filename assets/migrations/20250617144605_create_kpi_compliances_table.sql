-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_compliances (
    id INT PRIMARY KEY IDENTITY(1, 1),
    kpi_register_id INT NOT NULL,
    compliance_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_kpi_compliances_kpi_register FOREIGN KEY (kpi_register_id) REFERENCES kpi_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_compliances_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_compliances_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE kpi_compliances;

-- +goose StatementEnd
