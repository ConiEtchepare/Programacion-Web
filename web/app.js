const API_URL = "http://localhost:8080/productos";

// Elementos del DOM
const productoForm = document.getElementById("producto-form");
const productosLista = document.getElementById("productos-lista");
const modalEditar = document.getElementById("modal-editar");
const editForm = document.getElementById("edit-form");
const closeModalBtn = document.querySelector(".close-modal");
const btnCancelar = document.getElementById("btn-cancelar");

// 1. Cargar productos al iniciar la página (GET)
document.addEventListener("DOMContentLoaded", obtenerProductos);

async function obtenerProductos() {
    try {
        const respuesta = await fetch(API_URL);
        if (!respuesta.ok) throw new Error("Error al consultar la API");

        const productos = await respuesta.json();
        renderizarProductos(productos);
    } catch (error) {
        console.error("Error al obtener productos:", error);
        productosLista.innerHTML = `<li class="empty-msg">No se pudieron cargar los productos.</li>`;
    }
}

// Renderizar productos con botones de Editar y Eliminar
function renderizarProductos(productos) {
    productosLista.innerHTML = "";

    if (!productos || productos.length === 0) {
        productosLista.innerHTML = `<li class="empty-msg">No hay productos registrados en el catálogo.</li>`;
        return;
    }

    productos.forEach(prod => {
        const li = document.createElement("li");
        li.className = "producto-item";
        li.innerHTML = `
            <div class="producto-info">
                <strong>${prod.Nombre}</strong>
                <span class="categoria-badge">${prod.Categoria}</span>
                <span class="precio">$${prod.Precio}</span>
            </div>
            <div class="acciones-item">
                <button class="btn-edit" data-id="${prod.ID}" data-nombre="${prod.Nombre}" data-categoria="${prod.Categoria}" data-precio="${prod.Precio}">Editar</button>
                <button class="btn-delete" data-id="${prod.ID}">Eliminar</button>
            </div>
        `;
        productosLista.appendChild(li);
    });
}

// 2. Agregar un nuevo producto (POST)
if (productoForm) {
    productoForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const nombreVal = document.getElementById("nombre").value.trim();
        const categoriaVal = document.getElementById("categoria").value.trim();
        const precioVal = document.getElementById("precio").value.toString().trim().replace(',', '.');

        const nuevoProducto = {
            nombre: nombreVal,
            categoria: categoriaVal,
            precio: precioVal
        };

        try {
            const respuesta = await fetch(API_URL, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(nuevoProducto)
            });

            if (!respuesta.ok) throw new Error("Error al crear el producto");

            productoForm.reset();
            obtenerProductos();
        } catch (error) {
            console.error("Error al guardar:", error);
            alert("Hubo un error al guardar el producto.");
        }
    });
}

// 3. Abrir Modal de Edición o Eliminar (Delegación de Eventos)
if (productosLista) {
    productosLista.addEventListener("click", async (e) => {
        // Clic en Editar
        if (e.target.classList.contains("btn-edit")) {
            const id = e.target.getAttribute("data-id");
            const nombre = e.target.getAttribute("data-nombre");
            const categoria = e.target.getAttribute("data-categoria");
            const precio = e.target.getAttribute("data-precio");

            document.getElementById("edit-id").value = id;
            document.getElementById("edit-nombre").value = nombre;
            document.getElementById("edit-categoria").value = categoria;
            document.getElementById("edit-precio").value = precio;

            modalEditar.style.display = "flex";
        }

        // Clic en Eliminar (DELETE)
        if (e.target.classList.contains("btn-delete")) {
            const id = e.target.getAttribute("data-id");

            if (confirm(`¿Estás seguro de eliminar el producto #${id}?`)) {
                try {
                    const respuesta = await fetch(`${API_URL}/${id}`, { method: "DELETE" });
                    if (!respuesta.ok) throw new Error("Error al eliminar el producto");
                    obtenerProductos();
                } catch (error) {
                    console.error("Error al eliminar:", error);
                    alert("No se pudo eliminar el producto.");
                }
            }
        }
    });
}

// Cerrar Modal
const cerrarModal = () => { modalEditar.style.display = "none"; };
if (closeModalBtn) closeModalBtn.addEventListener("click", cerrarModal);
if (btnCancelar) btnCancelar.addEventListener("click", cerrarModal);
window.addEventListener("click", (e) => { if (e.target === modalEditar) cerrarModal(); });

// 4. Actualizar Producto (PUT)
if (editForm) {
    editForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const id = document.getElementById("edit-id").value;
        const nombreVal = document.getElementById("edit-nombre").value.trim();
        const categoriaVal = document.getElementById("edit-categoria").value.trim();
        const precioVal = document.getElementById("edit-precio").value.toString().trim().replace(',', '.');

        const productoActualizado = {
            nombre: nombreVal,
            categoria: categoriaVal,
            precio: precioVal
        };

        try {
            const respuesta = await fetch(`${API_URL}/${id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(productoActualizado)
            });

            if (!respuesta.ok) throw new Error("Error al actualizar el producto");

            cerrarModal();
            obtenerProductos();
        } catch (error) {
            console.error("Error al actualizar:", error);
            alert("Hubo un error al actualizar el producto.");
        }
    });
}
