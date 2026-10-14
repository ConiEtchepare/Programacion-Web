package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	// Importar pgxpool para manejar la conexión nativa con pgx/v5
	"github.com/jackc/pgx/v5/pgxpool"
	
	// Paquete generado por sqlc
	"catalogo-LoDeKiki/db"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "index.html")
}

func main() {
	ctx := context.Background()

	// Obtener la URL de conexión desde el entorno o usar una por defecto
	dbSource := os.Getenv("DB_URL")
	if dbSource == "" {
		dbSource = "postgresql://root:secret@localhost:5432/lodekiki_db?sslmode=disable"
	}

	// Establecer la conexión usando pgxpool (nativo para sqlc con pgx/v5)
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatalf("Error al crear el pool de conexiones: %v", err)
	}
	defer pool.Close()

	// Verificar que responda al hacer ping
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("No se pudo conectar al servidor de base de datos: %v", err)
	}
	log.Println("¡Conexión exitosa a la base de datos!")

    // Crear la instancia del repositorio de sqlc (store) pasándole el pool
	store := db.New(pool)
	_ = store // Evita error si aún no se usa en los handlers

    // Crear el enrutador HTTP
	mux := http.NewServeMux()

	// Ruta principal (index.html)
	mux.HandleFunc("GET /", handleRoot)

	// Endpoints CRUD para productos
	mux.HandleFunc("POST /productos", handleCreateProducto(store))
	mux.HandleFunc("GET /productos", handleListProductos(store))
	mux.HandleFunc("GET /productos/{id}", handleGetProducto(store))
	mux.HandleFunc("PUT /productos/{id}", handleUpdateProducto(store))
	mux.HandleFunc("DELETE /productos/{id}", handleDeleteProducto(store))

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
