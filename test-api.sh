#!/bin/bash
set -e

echo "=== 1. Limpiando contenedores previos ==="
docker compose down -v --remove-orphans || true

echo "=== 2. Generando código con sqlc ==="
sqlc generate

echo "=== 3. Compilando la aplicación Go ==="
go build -o lodekiki-app .

echo "=== 4. Levantando contenedor de base de datos ==="
docker compose up -d db

echo "=== 5. Esperando que la base de datos esté lista ==="
until docker exec postgres_tp2 pg_isready -U root -d lodekiki_db; do
    sleep 1
done

# Pausa extra para asegurar que PostgreSQL acepte conexiones TCP
sleep 3

echo "=== 6. Iniciando la aplicación con las variables del .env ==="
# Cargamos las variables de entorno de forma segura para el proceso de Go
export $(grep -v '^#' .env | xargs)
./lodekiki-app &
APP_PID=$!

# Damos 3 segundos para que el servidor web levante en el puerto 8080
sleep 3

echo "=== 7. Ejecutando script de pruebas de endpoints (request.bash) ==="
chmod +x request.bash

if ./request.bash; then
    echo "¡Pruebas de endpoints exitosas!"
    kill -9 $APP_PID 2>/dev/null || true
    docker compose down -v
    echo "¡Pruebas de API finalizadas con éxito!"
    exit 0
else
    echo "¡ERROR: Fallaron las pruebas de los endpoints!"
    kill -9 $APP_PID 2>/dev/null || true
    docker compose down -v
    exit 1
fi
