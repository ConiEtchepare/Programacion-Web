-- name: CreateProducto :one
INSERT INTO producto (
    nombre, categoria, precio
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetProducto :one
SELECT * FROM producto
WHERE id = $1 LIMIT 1;

-- name: ListProductos :many
SELECT * FROM producto
ORDER BY nombre;

-- name: UpdateProducto :one
UPDATE producto
SET 
    nombre = $1,
    categoria = $2,
    precio = $3
WHERE id = $4
RETURNING *;

-- name: DeleteProducto :exec
DELETE FROM producto
WHERE id = $1;
