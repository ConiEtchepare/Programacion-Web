-include .env
export

.PHONY: test-api
test-api:
	./test-api.sh

env:
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Archivo .env creado automáticamente a partir de env.example"; \
	fi

# --- REGLA 1: Pruebas unitarias y de base de datos (Go tests) ---
test: env
	@echo "=== 1. Limpiando contenedores y volúmenes previos ==="
	docker compose down -v --remove-orphans || true

	@echo "=== 2. Generando código con sqlc ==="
	sqlc generate

	@echo "=== 3. Levantando contenedor de base de datos ==="
	docker compose up -d db

	@echo "=== 4. Esperando que la base de datos esté lista ==="
	@until docker exec postgres_tp2 pg_isready -U root -d lodekiki_db; do \
		sleep 1; \
	done

	@echo "=== 5. Ejecutando tests de Go ==="
	cp .env db/.env
	export $$(grep -v '^#' .env | xargs) && go test -v -count=1 ./...
	rm -f db/.env

	@echo "=== 6. Limpiando entorno de pruebas ==="
	docker compose down -v


# --- REGLA 2: Pruebas de la API REST (endpoints con request.bash) ---
test-api: env
	@echo "=== 1. Limpiando contenedores previos ==="
	docker compose down -v --remove-orphans || true

	@echo "=== 2. Generando código con sqlc ==="
	sqlc generate

	@echo "=== 3. Compilando la aplicación Go ==="
	go build -o lodekiki-app .

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

	@echo "=== 7. Ejecutando script de pruebas de endpoints (request.bash) ==="
	chmod +x request.bash
	./request.bash || (echo "Error en los endpoints"; kill $$APP_PID; docker compose down -v; exit 1)

	@echo "=== 8. Limpiando servidor y entorno ==="
	kill $$APP_PID || true
	docker compose down -v
	@echo "¡Pruebas de API finalizadas con éxito!"
