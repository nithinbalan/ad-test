-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS system_register_job_titles;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS system_register_operational_locations;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS system_register_recovery_needed_tools;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS system_register_recovery_locations;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS system_register_slas;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS system_register_processes;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS system_registers;
-- +goose StatementEnd

