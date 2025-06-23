-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_library_versions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_library_id INT NOT NULL,
    version NVARCHAR(50) NOT NULL,
    year INT NOT NULL,
    overview NVARCHAR(MAX) NULL,
    introduction NVARCHAR(MAX) NULL,
    domains NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_clv_compliance_libraries FOREIGN KEY (compliance_library_id) REFERENCES compliance_libraries(id) ON DELETE CASCADE ON UPDATE NO ACTION
);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_library_versions;
-- +goose StatementEnd