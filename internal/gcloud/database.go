package gcloud

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Vaansh/gore/internal/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

// Initialize db

func InitDatabase() (*sql.DB, error) {
	cfg := config.ReadDatabaseConfig()

	var (
		dbUser                 = cfg.Username
		dbPwd                  = cfg.Password
		dbName                 = cfg.Database
		instanceConnectionName = cfg.InstanceId
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	dbCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	var opts []cloudsqlconn.Option
	dialer, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	dbCfg.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return dialer.Dial(ctx, instanceConnectionName)
	}

	dbURI := stdlib.RegisterConnConfig(dbCfg)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	dbPool.SetMaxOpenConns(10)
	return dbPool, nil
}
