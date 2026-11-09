package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"catalogo-LoDeKiki/internal/db"
	"catalogo-LoDeKiki/views"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	dbUrl := os.Getenv("DB_URL")
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer pool.Close()

	store := db.New(pool)
	mux := http.NewServeMux()

    // Servir archivos estáticos (CSS, imágenes)
    mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		productos, err := store.ListProductos(r.Context())
		if err != nil {
			http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		views.IndexPage(productos).Render(r.Context(), w)
	})

	mux.HandleFunc("POST /productos", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
			return
		}

		nombre := r.FormValue("nombre")
		categoria := r.FormValue("categoria")
		precioStr := r.FormValue("precio")

		// Convertir el string de precio a pgtype.Numeric
		var precioNum pgtype.Numeric
		if err := precioNum.Scan(precioStr); err != nil {
			http.Error(w, "Formato de precio inválido", http.StatusBadRequest)
			return
		}

		_, err := store.CreateProducto(r.Context(), db.CreateProductoParams{
			Nombre:    nombre,
			Categoria: categoria,
			Precio:    precioNum,
		})
		if err != nil {
			http.Error(w, "Error al guardar el producto", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("Servidor SSR escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
