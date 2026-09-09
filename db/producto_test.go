package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateProducto(t *testing.T) {
	var precioNumeric pgtype.Numeric
	err := precioNumeric.Scan("2500.50")
	if err != nil {
		t.Fatalf("Error al parsear el precio: %v", err)
	}

	arg := CreateProductoParams{
		Nombre:    "Granola Artezanal con Almendras",
		Categoria: "Snacks Saludables",
		Precio:    precioNumeric,
	}

	producto, err := testQueries.CreateProducto(context.Background(), arg)
	if err != nil {
		t.Fatalf("Error al crear producto: %v", err)
	}

	if producto.Nombre != arg.Nombre {
		t.Errorf("Se esperaba nombre %s pero se obtuvo %s", arg.Nombre, producto.Nombre)
	}
}
