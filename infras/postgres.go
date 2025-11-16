package infras

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/fikriahmadf/outbox-examples/shared/failure"
)

// TransactionBlock contains a transaction block
type TransactionBlock func(db *sqlx.Tx, c chan error)

// PostgresConn wraps a pair of read/write PostgreSQL connections.
type PostgresConn struct {
	Read  *sqlx.DB
	Write *sqlx.DB
}

// ProvidePostgresConn is the provider for PostgresConn.
func ProvidePostgresConn(config *configs.Config) *PostgresConn {
	return &PostgresConn{
		Read:  CreatePostgresReadConn(*config),
		Write: CreatePostgresWriteConn(*config),
	}
}

// CreatePostgresReadConn creates a database connection for read access.
func CreatePostgresReadConn(config configs.Config) *sqlx.DB {
	return CreatePostgresDBConnection(
		"read",
		config.DB.Postgres.Read.User,
		config.DB.Postgres.Read.Password,
		config.DB.Postgres.Read.Host,
		config.DB.Postgres.Read.Port,
		config.DB.Postgres.Read.Name,
		config.DB.Postgres.Read.SSLMode,
		config.DB.Postgres.Read.MaxConnLifetime,
		config.DB.Postgres.Read.MaxIdleConn,
		config.DB.Postgres.Read.MaxOpenConn)

}

// CreatePostgresWriteConn creates a database connection for write access.
func CreatePostgresWriteConn(config configs.Config) *sqlx.DB {
	return CreatePostgresDBConnection(
		"write",
		config.DB.Postgres.Write.User,
		config.DB.Postgres.Write.Password,
		config.DB.Postgres.Write.Host,
		config.DB.Postgres.Write.Port,
		config.DB.Postgres.Write.Name,
		config.DB.Postgres.Write.SSLMode,
		config.DB.Postgres.Write.MaxConnLifetime,
		config.DB.Postgres.Write.MaxIdleConn,
		config.DB.Postgres.Write.MaxOpenConn)

}

// CreatePostgresDBConnection creates a database connection.
func CreatePostgresDBConnection(connType, username, password, host, port, dbName, sslmode string, maxConnLifetime time.Duration, maxIdleConn, maxOpenConn int) *sqlx.DB {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		username,
		password,
		dbName,
		sslmode)

	if password == "" {
		conn = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=%s",
			host,
			port,
			username,
			dbName,
			sslmode)
	}

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("type", connType).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Failed connecting to Postgres database")
	} else {
		log.
			Info().
			Str("type", connType).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Connected to Postgres database")
	}

	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)

	return db
}

// WithTransaction performs queries with transaction
func (m *PostgresConn) WithTransaction(block TransactionBlock) (err error) {
	e := make(chan error)
	tx, err := m.Write.Beginx()
	if err != nil {
		log.Err(err).Msg("error begin transaction")
		err = failure.InternalError(err)
		return
	}
	go block(tx, e)
	err = <-e
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			log.Err(errTx).Msg("error rollback transaction")
			err = failure.InternalError(err)
		}
		return
	}
	err = tx.Commit()
	return
}
