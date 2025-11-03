package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pgvector/pgvector-go"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(" Nenhum .env encontrado.")
	}
}

func GetDB() *pgxpool.Pool {
	dbURI := os.Getenv("PSYCOPG_DB_URI")
	if dbURI == "" {
		log.Fatal(" PSYCOPG_DB_URI não definido no ambiente.")
	}

	config, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		log.Fatalf(" Erro ao parsear string de conexão: %v", err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterDefaultPgType(pgvector.Vector{}, "vector")
		return nil
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.HealthCheckPeriod = 2 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf(" Erro ao criar pool de conexões: %v", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf(" Banco inacessível: %v", err)
	}

	fmt.Println(" Conexão PostgreSQL (pgxpool) bem-sucedida")

	var dbName, user, version string
	dbpool.QueryRow(ctx, "SELECT current_database();").Scan(&dbName)
	dbpool.QueryRow(ctx, "SELECT current_user;").Scan(&user)
	dbpool.QueryRow(ctx, "SHOW server_version;").Scan(&version)
	log.Printf(" Banco conectado: %s | Usuário: %s | Versão: %s", dbName, user, version)

	rows, _ := dbpool.Query(ctx, "SHOW search_path;")
	defer rows.Close()
	for rows.Next() {
		var path string
		rows.Scan(&path)
		log.Printf(" search_path: %s", path)
	}

	return dbpool
}
