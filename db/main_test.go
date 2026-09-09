package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connStr := "postgres://root:secret@localhost:5432/catalogo?sslmode=disable"
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos de pruebas: %v", err)
	}
	defer conn.Close(ctx)

	testQueries = New(conn)

	os.Exit(m.Run())
}
