// Lógica pura del frontend: no toca el DOM ni la red, por eso se testea
// con valores de entrada y salida, sin dobles.

export const MAX_TITULO = 100;

// Espejo de ValidateTaskTitle del backend. Sirve para avisar antes de
// mandar el pedido; quien decide de verdad sigue siendo el backend.
export function validarTitulo(titulo) {
  const limpio = (titulo ?? "").trim();
  if (limpio === "") {
    return { ok: false, error: "el título no puede estar vacío" };
  }
  if (limpio.length > MAX_TITULO) {
    return { ok: false, error: `el título no puede superar los ${MAX_TITULO} caracteres` };
  }
  return { ok: true, error: null };
}

// Qué botones ofrece la UI según el estado de la tarea. Refleja las
// transiciones permitidas del backend: "archivada" es un estado terminal.
export function accionesPara(estado) {
  switch (estado) {
    case "pendiente":
      return ["completar", "archivar", "eliminar"];
    case "completada":
      return ["reabrir", "archivar", "eliminar"];
    case "archivada":
      return ["eliminar"];
    default:
      return [];
  }
}

// La API devuelve la fecha como ISO ("2030-05-01T00:00:00Z"); la UI
// muestra solo YYYY-MM-DD.
export function formatearVencimiento(fecha) {
  if (!fecha) return "";
  return String(fecha).substring(0, 10);
}

// Escapa texto antes de meterlo en innerHTML (evita inyección de HTML).
export function escaparHtml(texto) {
  return String(texto)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}