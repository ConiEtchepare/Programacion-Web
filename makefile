.PHONY: test

test:
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

	@echo "=== 5. Ejecutando tests ==="
	go test -v -count=1 ./db/...

	@echo "=== 6. Limpiando entorno de pruebas ==="
	docker compose down -v
