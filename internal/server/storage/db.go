package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	dbdir "github.com/tony-spark/metrico/db"
	"github.com/tony-spark/metrico/internal/server/models"
	"time"
)

type PgDatabaseManager struct {
	db  *sql.DB
	gdb *GaugeDB
	cdb *CounterDB
}

type GaugeDB struct {
	db *sql.DB
}

type CounterDB struct {
	db *sql.DB
}

func NewPgManager(dsn string) (*PgDatabaseManager, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := iofs.New(dbdir.EmbeddedDBFiles, "migrations")
	if err != nil {
		return nil, err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", driver, dsn)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &PgDatabaseManager{
		db:  db,
		gdb: &GaugeDB{db: db},
		cdb: &CounterDB{db: db},
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

func (pgm PgDatabaseManager) GaugeRepository() models.GaugeRepository {
	return pgm.gdb
}

func (pgm PgDatabaseManager) CounterRepository() models.CounterRepository {
	return pgm.cdb
}

func (pgm PgDatabaseManager) Close() error {
	return pgm.db.Close()
}

func (gdb GaugeDB) GetByName(ctx context.Context, name string) (*models.GaugeValue, error) {
	row := gdb.db.QueryRowContext(ctx, "SELECT name, value FROM gauges WHERE name = $1", name)
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

func (gdb GaugeDB) Save(ctx context.Context, name string, value float64) (*models.GaugeValue, error) {
	g := models.GaugeValue{
		Name:  name,
		Value: value,
	}

	result, err := gdb.db.ExecContext(ctx,
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

func (gdb GaugeDB) GetAll(ctx context.Context) ([]*models.GaugeValue, error) {
	gs := make([]*models.GaugeValue, 0)

	rows, err := gdb.db.QueryContext(ctx, `SELECT name, value FROM gauges`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.GaugeValue
		err := rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, err
		}

		gs = append(gs, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gs, nil
}

func (cdb CounterDB) GetByName(ctx context.Context, name string) (*models.CounterValue, error) {
	row := cdb.db.QueryRowContext(ctx, "SELECT name, value FROM counters WHERE name = $1", name)
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

func (cdb CounterDB) AddAndSave(ctx context.Context, name string, value int64) (*models.CounterValue, error) {
	result, err := cdb.db.ExecContext(ctx,
		`INSERT INTO counters(name, value) VALUES ($1, $2)
				ON CONFLICT (name) DO UPDATE 
				SET value = counters.value + excluded.value`,
		name, value)

	if err != nil {
		return nil, err
	}

	if err = checkOneAffected(result); err != nil {
		return nil, err
	}

	c, err := cdb.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (cdb CounterDB) Save(ctx context.Context, name string, value int64) (*models.CounterValue, error) {
	c := models.CounterValue{
		Name:  name,
		Value: value,
	}

	result, err := cdb.db.ExecContext(ctx,
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

func (cdb CounterDB) GetAll(ctx context.Context) ([]*models.CounterValue, error) {
	cs := make([]*models.CounterValue, 0)

	rows, err := cdb.db.QueryContext(ctx, `SELECT name, value FROM counters`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.CounterValue
		err := rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, err
		}

		cs = append(cs, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
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
