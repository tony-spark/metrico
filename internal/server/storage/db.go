package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"

	dbdir "github.com/tony-spark/metrico/db"
	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
)

type PgDatabaseManager struct {
	db  *sql.DB
	mdb MetricDВ
}

type MetricDВ struct {
	db *sql.DB
}

func NewPgManager(dsn string) (*PgDatabaseManager, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	driver, err := iofs.New(dbdir.EmbeddedDBFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("could not find db migrations: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("could not initialize db migrations: %w", err)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("could not execute db migrations: %w", err)
	}

	return &PgDatabaseManager{
		db:  db,
		mdb: MetricDВ{db: db},
	}, nil
}

func (pgm PgDatabaseManager) Check(ctx context.Context) (bool, error) {
	ct, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := pgm.db.PingContext(ct); err != nil {
		return false, err
	}
	return true, nil
}

func (pgm PgDatabaseManager) MetricRepository() models.MetricRepository {
	return pgm.mdb
}

func (pgm PgDatabaseManager) Close() error {
	return pgm.db.Close()
}

func (db MetricDВ) GetGaugeByName(ctx context.Context, name string) (*models.GaugeValue, error) {
	row := db.db.QueryRowContext(ctx, "SELECT name, value FROM gauges WHERE name = $1", name)
	var g models.GaugeValue

	err := row.Scan(&g.Name, &g.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (db MetricDВ) SaveGauge(ctx context.Context, name string, value float64) (*models.GaugeValue, error) {
	g := models.GaugeValue{
		Name:  name,
		Value: value,
	}

	result, err := db.db.ExecContext(ctx,
		`INSERT INTO gauges(name, value) VALUES ($1, $2)
				ON CONFLICT (name) DO UPDATE 
				SET value = excluded.value`,
		name, value)

	if err != nil {
		return nil, err
	}

	if err = checkOneAffected(result); err != nil {
		return nil, err
	}

	return &g, nil
}

func (db MetricDВ) SaveAllGauges(ctx context.Context, gs []models.GaugeValue) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO gauges(name, value) VALUES ($1, $2)
				ON CONFLICT (name) DO UPDATE 
				SET value = excluded.value`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, g := range gs {
		if _, err := stmt.ExecContext(ctx, g.Name, g.Value); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db MetricDВ) getAllGauges(ctx context.Context) ([]models.GaugeValue, error) {
	gs := make([]models.GaugeValue, 0)

	rows, err := db.db.QueryContext(ctx, `SELECT name, value FROM gauges`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.GaugeValue
		err = rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, err
		}

		gs = append(gs, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gs, nil
}

func (db MetricDВ) GetCounterByName(ctx context.Context, name string) (*models.CounterValue, error) {
	row := db.db.QueryRowContext(ctx, "SELECT name, value FROM counters WHERE name = $1", name)
	var g models.CounterValue

	err := row.Scan(&g.Name, &g.Value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (db MetricDВ) AddAndSaveCounter(ctx context.Context, name string, value int64) (*models.CounterValue, error) {
	row := db.db.QueryRowContext(ctx,
		`INSERT INTO counters(name, value) VALUES ($1, $2)
				ON CONFLICT (name) DO UPDATE 
				SET value = counters.value + excluded.value
				RETURNING counters.name, counters.value`,
		name, value)

	var c models.CounterValue

	err := row.Scan(&c.Name, &c.Value)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db MetricDВ) AddAndSaveAllCounters(ctx context.Context, cs []models.CounterValue) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO counters(name, value) VALUES ($1, $2)
				ON CONFLICT (name) DO UPDATE 
				SET value = counters.value + excluded.value`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range cs {
		if _, err := stmt.ExecContext(ctx, c.Name, c.Value); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db MetricDВ) SaveCounter(ctx context.Context, name string, value int64) (*models.CounterValue, error) {
	c := models.CounterValue{
		Name:  name,
		Value: value,
	}

	result, err := db.db.ExecContext(ctx,
		`INSERT INTO counters(name, value) VALUES ($1, $2)
				ON CONFLICT (name) DO UPDATE 
				SET value = excluded.value`,
		name, value)

	if err != nil {
		return nil, err
	}

	if err = checkOneAffected(result); err != nil {
		return nil, err
	}

	return &c, nil
}

func (db MetricDВ) getAllCounters(ctx context.Context) ([]models.CounterValue, error) {
	cs := make([]models.CounterValue, 0)

	rows, err := db.db.QueryContext(ctx, `SELECT name, value FROM counters`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.CounterValue
		err = rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, err
		}

		cs = append(cs, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}

func (db MetricDВ) GetAll(ctx context.Context) ([]model.Metric, error) {
	ms := make([]model.Metric, 0)

	gs, err := db.getAllGauges(ctx)
	if err != nil {
		return nil, err
	}

	cs, err := db.getAllCounters(ctx)
	if err != nil {
		return nil, err
	}

	for _, g := range gs {
		ms = append(ms, g)
	}
	for _, c := range cs {
		ms = append(ms, c)
	}

	return ms, nil
}

func checkOneAffected(r sql.Result) error {
	rows, err := r.RowsAffected()

	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}

	return nil
}
