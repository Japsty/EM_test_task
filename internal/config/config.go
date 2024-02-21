package config

import (
	"database/sql"
)

func CheckMigrations(db *sql.DB) {
	//if os.Getenv("MIGRATE_DB") == "true" {
	//	provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)
	//	if err != nil {
	//		log.Fatal("Main failed to create NewProvider for migration")
	//	}
	//	_, err = provider.Up(context.Background())
	//	if err != nil {
	//		slog.Info("Failed to up migration")
	//		return
	//}
}
