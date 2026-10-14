# Programacion-Web

Para esta tercer entrega continuamos con el desarrollo de la aplicación web exponiendo la lógica como API para la gestión del catálogo de productos de Lo De Kiki.

### Requisitos previos
Para ejecutar el proyecto y las pruebas debemos tener instalado:
- Go
- Docker y Docker Compose
- sqlc
- Make

### Pruebas de endpoints
Las pruebas se ejecutan de forma automatizada mediante un script que levanta el entorno, compila la aplicación, ejecuta todas las peticiones CRUD a través del archivo 'request.bash' y realiza la limpieza de contenedores al finalizar. 
Comando para ejecutar las pruebas: ./test-api.sh

