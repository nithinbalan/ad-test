package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateIncidentStrategicObjectivesTableSql, downCreateIncidentStrategicObjectivesTableSql)
}

func upCreateIncidentStrategicObjectivesTableSql(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downCreateIncidentStrategicObjectivesTableSql(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
