package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type PgPool struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Schema   string
	DbPool   *pgxpool.Pool
}

var pgPool *PgPool

func dbMigrate() {
	workingDir, _ := os.Getwd()
	log.Info().Msg(workingDir)
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=disable",
			pgPool.Username, pgPool.Password, pgPool.Host, pgPool.Port, pgPool.Database, pgPool.Schema))
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}

func LotChanged(lot *Lot) bool {
	if pgPool == nil {
		initDb()
	}

	query := `select id, status, name, description, create_date, bid_end_time,
       image, is_stopped, has_appeals, is_annulled from lot where id = $1`

	row := pgPool.DbPool.QueryRow(context.TODO(), query, lot.Id)
	storedLot := Lot{}

	if err := row.Scan(&storedLot.Id, &storedLot.Status, &storedLot.Name, &storedLot.Description, &storedLot.CreateDate, &storedLot.BidEndTime,
		&storedLot.Image, &storedLot.IsStopped, &storedLot.HasAppeals, &storedLot.IsAnnulled); err != nil {
		if err == pgx.ErrNoRows {
			return true
		}
	}

	return lot.Status != storedLot.Status ||
		lot.Name != storedLot.Name ||
		lot.Description != storedLot.Description ||
		!lot.CreateDate.Equal(storedLot.CreateDate) ||
		!lot.BidEndTime.Equal(storedLot.BidEndTime) ||
		!bytes.Equal(lot.Image, storedLot.Image) ||
		lot.IsStopped != storedLot.IsStopped ||
		lot.HasAppeals != storedLot.HasAppeals ||
		lot.IsAnnulled != storedLot.IsAnnulled
}

func StoreLot(lot *Lot) error {
	if pgPool == nil {
		initDb()
	}

	tx, err := pgPool.DbPool.Begin(context.TODO())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	statement := `insert into lot
    		(id, status, name, description, create_date, bid_end_time, image, is_stopped, has_appeals, is_annulled)
			values
    		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    		on conflict (id)
			do update set status = $2, name = $3, description = $4, create_date = $5, bid_end_time = $6, image = $7, is_stopped = $8, has_appeals = $9, is_annulled = $10`
	if _, err = tx.Exec(context.TODO(), statement, lot.Id, lot.Status, lot.Name, lot.Description, lot.CreateDate,
		lot.BidEndTime, lot.Image, lot.IsStopped, lot.HasAppeals, lot.IsAnnulled); err != nil {
		return err
	}
	return nil

}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func initDb() {
	dbPort, _ := strconv.Atoi(getenv("DB_PORT", "5432"))
	pgPool = &PgPool{
		Host:     getenv("DB_HOST", "localhost"),
		Port:     dbPort,
		Username: getenv("DB_USERNAME", "postgres"),
		Password: getenv("DB_PASSWORD", "password"),
		Database: getenv("DB_DATABASE", "postgres"),
		Schema:   getenv("DB_SCHEMA", "public"),
	}
	connstring := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s search_path=%s target_session_attrs=read-write",
		pgPool.Host, pgPool.Port, pgPool.Database, pgPool.Username, pgPool.Password, pgPool.Schema)

	connConfig, err := pgxpool.ParseConfig(connstring)
	if err != nil {
		log.Panic().Err(err).Msg("Unable to parse config")
		os.Exit(1)
	}

	if pgPool.DbPool, err = pgxpool.NewWithConfig(context.TODO(), connConfig); err != nil {
		log.Panic().Err(err).Msg("Unable to create connection pool")
		os.Exit(1)
	}

	dbMigrate()
}
