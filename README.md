Programación Web — Catálogo Lo De Kiki

Para esta segunda entrega continuamos con el desarrollo de la aplicación web y el servicio de persistencia desarrollado en Go. 
La entrega cuenta con los siguientes archivos:

Programacion-Web/
  db/                     - Código generado por sqlc y tests de integración
      db.go               - Estructura e interfaz de base de datos
      models.go           - Modelos/Structs en Go (Producto)
      queries.sql.go      - Métodos CRUD autogenerados
      main_test.go        - Configuración principal del entorno de pruebas
      producto_test.go    - Tests de integración del catálogo (testint)
  docker-compose.yml      - Configuración del contenedor de PostgreSQL
  Makefile                - Automatización del flujo de generación y ejecución de tests
  queries.sql             - Consultas SQL para sqlc (CRUD)
  schema.sql              - Definición de la tabla 'producto'
  sqlc.yaml               - Configuración de sqlc
  go.mod                  - Módulo de Go y dependencias
  README.md               - Documentación del proyecto

Para llevar a cabo dicha entrega, tuvimos que contar con las siguientes herramientes:
- Go (ya instalado para la primer entrega)
- Docker y Docker Compose
- sqlc
- Make

Testing: 
  Para la ejecución de pruebas de integración se levanta un contenedor PostgreSQL en segundo plano, genera el código con sqlc, ejecuta las pruebas unitarias e integrales sobre la base de datos y finalmente limpia el entorno de contenedores.
  
  Para poder realizar el testeo abrimos la terminal desde la raiz del proyecto y utilizamos el comando: sudo make test
  
  Para verificar la correcta ejecución el sistema realiza los siguientes pasos mostrandolos por la terminal:
    1. Limpieza previa:apaga y elimina contenedores/volúmenes antiguos (docker compose down -v).
    2. Generación de código:mapea schema.sql y queries.sql a código nativo en Go usando sqlc generate.
    3. Levantado de la base de datos: inicia el contenedor PostgreSQL definido en docker-compose.yml.
    4. Verificación de disponibilidad:espera a que PostgreSQL esté aceptando conexiones (pg_isready).
    5. Ejecución de tests (testing):corre la suite de pruebas ignorando la caché:
       bash
       go test -v -count=1 ./db/...
    6. Limpieza del entorno:al finalizar, apaga el contenedor de PostgreSQL para liberar recursos (docker compose down -v).

Pruebas que realizamos:
En `db/producto_test.go` se evalúan las siguientes operaciones CRUD sobre la entidad `Producto`:

- Creación (CreateProducto) inserción de productos de dietética (ej. Almendra non pareil) utilizando tipos precisos como `pgtype.Numeric` para el precio.
- Lectura (GetProducto / ListProductos):consulta individual por ID y listado de catálogo.
- Actualización (UpdateProducto): modificación de precios y categorías.
- Baja (DeleteProducto): eliminación de registros de la base de datos.
