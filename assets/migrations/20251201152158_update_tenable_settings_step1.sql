-- +goose Up
-- +goose NO TRANSACTION

ALTER TABLE tenable_settings
ADD status VARCHAR(10) NULL DEFAULT 'active';
