-include .env
export

.PHONY: env generate build run test test-api clean

env:
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Archivo .env creado automáticamente a partir de env.example"; \
	fi

# --- REGLA DE GENERACIÓN DE CÓDIGO (SQLC + TEMPL) ---
generate:
	sqlc generate
	templ generate

# --- REGLA DE COMPILACIÓN ---
build: generate
	go build -o lodekiki-app ./cmd/api

# --- REGLA DE EJECUCIÓN RÁPIDA (SSR) ---
run: env generate
	@echo "=== Levantando Base de Datos ==="
	docker compose up -d db
	@echo "=== Esperando conexión ==="
	@until docker exec postgres_tp2 pg_isready -U root -d lodekiki_db; do \
		sleep 1; \
	done
	@echo "=== Iniciando servidor SSR en http://localhost:8080 ==="
	export $$(grep -v '^#' .env | xargs) && go run ./cmd/api

# --- REGLA 1: Pruebas unitarias y de base de datos (Go tests) ---
test: env
	@echo "=== 1. Limpiando contenedores y volúmenes previos ==="
	docker compose down -v --remove-orphans || true

	@echo "=== 2. Generando código con sqlc y templ ==="
	$(MAKE) generate

	@echo "=== 3. Levantando contenedor de base de datos ==="
	docker compose up -d db

	@echo "=== 4. Esperando que la base de datos esté lista ==="
	@until docker exec postgres_tp2 pg_isready -U root -d lodekiki_db; do \
		sleep 1; \
	done

	@echo "=== 5. Ejecutando tests de Go ==="
	cp .env db/.env 2>/dev/null || true
	export $$(grep -v '^#' .env | xargs) && go test -v -count=1 ./...
	rm -f db/.env

	@echo "=== 6. Limpiando entorno de pruebas ==="
	docker compose down -v

# --- REGLA 2: Pruebas de la API REST (script de endpoints) ---
test-api: env
	@echo "=== 1. Limpiando contenedores previos ==="
	docker compose down -v --remove-orphans || true

	@echo "=== 2. Generando código con sqlc y templ ==="
	$(MAKE) generate

	@echo "=== 3. Compilando la aplicación Go ==="
	go build -o lodekiki-app ./cmd/api

	@echo "=== 4. Levantando contenedor de base de datos ==="
	docker compose up -d db

	@echo "=== 5. Esperando que la base de datos esté lista ==="
	@until docker exec postgres_tp2 pg_isready -U root -d lodekiki_db; do \
		sleep 1; \
	done

	@echo "=== 6. Iniciando la aplicación en segundo plano ==="
	./lodekiki-app &
	APP_PID=$$!
	sleep 2

	@echo "=== 7. Ejecutando script de pruebas de endpoints ==="
	chmod +x scripts/request.bash 2>/dev/null || chmod +x request.bash
	./scripts/request.bash 2>/dev/null || ./request.bash || (echo "Error en los endpoints"; kill $$APP_PID; docker compose down -v; exit 1)

	@echo "=== 8. Limpiando servidor y entorno ==="
	kill $$APP_PID || true
	docker compose down -v
	@echo "¡Pruebas de API finalizadas con éxito!"

clean:
	rm -f lodekiki-app
	docker compose down -v
