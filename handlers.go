package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"catalogo-LoDeKiki/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// 1. POST /productos - Crear un registro (201 Created)
func handleCreateProducto(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Nombre    string `json:"nombre"`
			Categoria string `json:"categoria"`
			Precio    string `json:"precio"` // Recibimos el precio como string para convertirlo de forma segura a Numeric
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.Nombre == "" {
			http.Error(w, "El campo 'nombre' es obligatorio", http.StatusBadRequest)
			return
		}

		// Convertir el string de precio a pgtype.Numeric
		var precioNumeric pgtype.Numeric
		if err := precioNumeric.Scan(req.Precio); err != nil {
			http.Error(w, "Formato de precio inválido", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		producto, err := store.CreateProducto(ctx, db.CreateProductoParams{
			Nombre:    req.Nombre,
			Categoria: req.Categoria,
			Precio:    precioNumeric,
		})
		if err != nil {
			http.Error(w, "Error al crear el producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(producto)
	}
}

// 2. GET /productos - Listar todos los registros
func handleListProductos(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		productos, err := store.ListProductos(ctx)
		if err != nil {
			http.Error(w, "Error al obtener los productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(productos)
	}
}

// 3. GET /productos/{id} - Obtener un registro por ID
func handleGetProducto(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		producto, err := store.GetProducto(ctx, id)
		if err != nil {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(producto)
	}
}

// 4. PUT /productos/{id} - Actualizar un registro
func handleUpdateProducto(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req struct {
			Nombre    string `json:"nombre"`
			Categoria string `json:"categoria"`
			Precio    string `json:"precio"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if req.Nombre == "" {
			http.Error(w, "El campo 'nombre' es obligatorio", http.StatusBadRequest)
			return
		}

		var precioNumeric pgtype.Numeric
		if err := precioNumeric.Scan(req.Precio); err != nil {
			http.Error(w, "Formato de precio inválido", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		producto, err := store.UpdateProducto(ctx, db.UpdateProductoParams{
			ID:        id,
			Nombre:    req.Nombre,
			Categoria: req.Categoria,
			Precio:    precioNumeric,
		})
		if err != nil {
			http.Error(w, "Error al actualizar el producto", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(producto)
	}
}

// 5. DELETE /productos/{id} - Eliminar un registro (204 No Content)
func handleDeleteProducto(store *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		err = store.DeleteProducto(ctx, id)
		if err != nil {
			http.Error(w, "Error al eliminar el producto", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
