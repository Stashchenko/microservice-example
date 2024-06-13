package postgresutil

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CustomLogger implements the pgx.QueryTracer interface
type CustomLogger struct {
	logger *slog.Logger
}

func (c CustomLogger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	c.logger.Info("SQL Query", "query", data.SQL)
	return ctx
}

func (c CustomLogger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		c.logger.Info("Query error", data.Err)
	}
}

func Connect(ctx context.Context, host, user, password, dbname string, port int) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(buildDSN(host, user, password, dbname, port))
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	config.ConnConfig.Tracer = CustomLogger{logger: slog.New(slog.NewJSONHandler(os.Stdout, nil))}
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return db, nil
}

func buildDSN(host, user, password, dbname string, port int) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, password, dbname, port)
}
