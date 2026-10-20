// Cliente de la API. Recibe fetch, la URL base y una función que dice
// quién es el usuario, en vez de usar los globales del navegador: así en
// los tests se le pasa un doble en lugar de salir a la red.
export function crearClienteApi({ fetchFn, baseUrl = "", getUserId = () => null }) {
  return async function request(path, options = {}) {
    const headers = { "Content-Type": "application/json", ...(options.headers || {}) };
    const userId = getUserId();
    if (userId !== null && userId !== undefined) {
      headers["X-User-Id"] = String(userId);
    }

    const res = await fetchFn(`${baseUrl}${path}`, { ...options, headers });

    let data = null;
    try {
      data = await res.json();
    } catch (_) {
      // respuestas sin cuerpo (204) no traen JSON
    }

    if (!res.ok) {
      throw new Error((data && data.error) || `Error ${res.status}`);
    }
    return data;
  };
}