// Package db, PostgreSQL bağlantı havuzunu ve şema bazlı goose
// migration'larını yönetir.
package db

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"

	"github.com/Libre-University/platform-api/core"
)

// DB bağlantı havuzunu sarar.
type DB struct{ pool *pgxpool.Pool }

// Open havuzu açar ve bağlantıyı doğrular.
func Open(ctx context.Context, url string) (*DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &DB{pool: pool}, nil
}

// Ping bağlantıyı denetler.
func (d *DB) Ping(ctx context.Context) error { return d.pool.Ping(ctx) }

// Close havuzu kapatır.
func (d *DB) Close() { d.pool.Close() }

// Migrate verilen şemayı oluşturur ve fsys içindeki goose migration'larını
// o şemaya ait sürüm tablosuyla (<schema>.goose_db_version) uygular. Her
// modül böylece kendi migration geçmişine sahip olur (ADR-0012).
func (d *DB) Migrate(ctx context.Context, schema string, fsys fs.FS) error {
	if !core.ValidName(schema) {
		return fmt.Errorf("db: geçersiz şema adı %q", schema)
	}
	if _, err := d.pool.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS "+schema); err != nil {
		return fmt.Errorf("db: şema %s oluşturulamadı: %w", schema, err)
	}
	sqldb := stdlib.OpenDBFromPool(d.pool)
	defer func() { _ = sqldb.Close() }()

	p, err := goose.NewProvider(database.DialectPostgres, sqldb, fsys,
		goose.WithTableName(schema+".goose_db_version"),
	)
	if err != nil {
		return fmt.Errorf("db: goose: %w", err)
	}
	if _, err := p.Up(ctx); err != nil {
		return fmt.Errorf("db: %s migration: %w", schema, err)
	}
	return nil
}
