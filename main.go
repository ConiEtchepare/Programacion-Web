package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Previene que se sirva index.html para rutas no encontradas (ej. /favicon.ico)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Establece el Content-Type explícitamente
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Sirve el archivo index.html ubicado en la raíz del proyecto
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", handleRoot)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
