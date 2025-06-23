-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_objective_kpis (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategic_objective_id INT NOT NULL,
    kpi_register_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    CONSTRAINT FK_so_kpi_strategic_objective FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_so_kpi_kpi_register FOREIGN KEY (kpi_register_id) REFERENCES kpi_register(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_so_kpi_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE strategic_objective_kpis;
-- +goose StatementEnd
