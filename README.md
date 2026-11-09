## Descripción del Proyecto

Este proyecto es la implementación de un sistema de gestión de catálogo para la dietética **"Lo De Kiki"**, desarrollado como parte del **Trabajo Práctico N°5** de la materia Programación Web. 

El objetivo principal de esta entrega fue migrar de una arquitectura tradicional de *Single Page Application (SPA)* comunicada mediante una API JSON, hacia una arquitectura de **Renderizado del Lado del Servidor (SSR)** utilizando Go y el motor de plantillas tipadas templ.

### Funcionalidades Implementadas
- **Visualización del Catálogo:** Renderizado dinámico de los productos almacenados en la base de datos directamente en el HTML.
- **Alta de Productos:** Formulario funcional con envío nativo (POST) para registrar nuevos ítems (Nombre, Categoría, Precio).
- **Diseño Personalizado:** Integración de la hoja de estilos mediante el despacho de archivos estáticos desde el servidor Go.

## Arquitectura y Decisiones de Diseño

### Renderizado del Lado del Servidor (SSR) con Templ
Se eligió templ porque permite escribir componentes similares a HTML o JSX, pero compilándolos a funciones seguras y fuertemente tipadas en Go. Esto elimina los errores de plantilla en tiempo de ejecución y mejora significativamente el SEO y el tiempo de carga inicial (TTV) en comparación con un cliente JavaScript.

### Patrón Post/Redirect/Get (PRG)
Para el manejo del formulario de creación de productos, se implementó el patrón **PRG**:
1. **POST:** El cliente envía los datos del formulario a /productos.
2. **Redirect:** El servidor procesa la inserción en la base de datos y responde con un código HTTP 303 See Other hacia /.
3. **GET:** El navegador realiza una nueva petición a / solicitando el catálogo actualizado.
*Beneficio:* Evita el reenvío accidental de datos (inserciones duplicadas) si el usuario recarga la página.

### Acceso a Datos
La capa de persistencia se maneja con **PostgreSQL** y **sqlc**. 
sqlc compila las consultas SQL nativas (schema.sql, queries.sql) y genera código Go optimizado y sin errores.

## Tecnologías Utilizadas

**Backend y Lógica:**
- Lenguaje: **Go**
- Enrutamiento: **net/http (ServeMux estándar)**
- Driver DB: **pgx/v5**

**Frontend y Vistas:**
- Motor SSR: **templ**
- Interfaz: **HTML5, CSS3 Custom (Plus Jakarta Sans)**

**Infraestructura y Herramientas:**
- Base de Datos: **PostgreSQL**
- Contenedores: **Docker & Docker Compose**
- Generador SQL: **sqlc**
- Automatización: **Make / Makefile**

## Ejecución Local

1. Ejecutar el proyecto mediante el comando: make run
Este comando compilará los templates, generará el código SQL, levantará PostgreSQL en Docker y arrancará el servidor.
2. Ingresar en el navegador a: http://localhost:8080
