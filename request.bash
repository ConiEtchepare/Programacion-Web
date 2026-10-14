#!/bin/bash

# URL base de nuestra API
BASE_URL="http://localhost:8080"

echo "=== 1. CREAR UN PRODUCTO (POST /productos) ==="
curl -X POST "$BASE_URL/productos" \
  -H "Content-Type: application/json" \
  -d '{"nombre": "Castaña de caju", "categoria": "Frutos secos", "precio": "12500.50"}'
echo -e "\n"

echo "=== 2. LISTAR TODOS LOS PRODUCTOS (GET /productos) ==="
curl -X GET "$BASE_URL/productos"
echo -e "\n"

echo "=== 3. OBTENER UN PRODUCTO POR ID (GET /productos/1) ==="
curl -X GET "$BASE_URL/productos/1"
echo -e "\n"

echo "=== 4. ACTUALIZAR UN PRODUCTO (PUT /productos/1) ==="
curl -X PUT "$BASE_URL/productos/1" \
  -H "Content-Type: application/json" \
  -d '{"nombre": "Castaña de caju w4 Brasil", "categoria": "Frutos secos", "precio": "15000.00"}'
echo -e "\n"

echo "=== 5. ELIMINAR UN PRODUCTO (DELETE /productos/1) ==="
curl -v -X DELETE "$BASE_URL/productos/1"
echo -e "\n"
