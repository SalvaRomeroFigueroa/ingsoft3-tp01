import { describe, it, expect, vi } from "vitest";
import { crearClienteApi } from "../lib/api.js";

const respuesta = (status, cuerpo) => ({
  ok: status >= 200 && status < 300,
  status,
  json: async () => cuerpo,
});

describe("crearClienteApi", () => {
  // CON MOCK: se reemplaza fetch por un doble y se verifica la INTERACCIÓN
  // (a quién llamó, cuántas veces y con qué header), no solo el resultado.
  it("llama a fetch una vez, con la URL completa y el header del usuario", async () => {
    const fetchFn = vi.fn().mockResolvedValue(respuesta(200, [{ id: 1 }]));
    const api = crearClienteApi({ fetchFn, baseUrl: "http://x", getUserId: () => 7 });

    const datos = await api("/api/lists");

    expect(datos).toEqual([{ id: 1 }]);
    expect(fetchFn).toHaveBeenCalledTimes(1);
    const [url, opciones] = fetchFn.mock.calls[0];
    expect(url).toBe("http://x/api/lists");
    expect(opciones.headers["X-User-Id"]).toBe("7");
  });

  it("no manda X-User-Id si no hay usuario logueado", async () => {
    const fetchFn = vi.fn().mockResolvedValue(respuesta(200, []));
    await crearClienteApi({ fetchFn })("/api/users/login");
    expect(fetchFn.mock.calls[0][1].headers["X-User-Id"]).toBeUndefined();
  });

  // CASOS DE ERROR
  it("rechaza con el mensaje del servidor cuando responde un error", async () => {
    const fetchFn = vi.fn().mockResolvedValue(
      respuesta(409, { error: "no se puede eliminar una lista con tareas pendientes" })
    );
    await expect(crearClienteApi({ fetchFn })("/api/lists/1", { method: "DELETE" })).rejects.toThrow(
      "no se puede eliminar una lista con tareas pendientes"
    );
  });

  it("usa un mensaje genérico si el error no trae cuerpo", async () => {
    const fetchFn = vi.fn().mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error("sin cuerpo");
      },
    });
    await expect(crearClienteApi({ fetchFn })("/x")).rejects.toThrow("Error 500");
  });

  it("devuelve null en una respuesta 204 sin cuerpo", async () => {
    const fetchFn = vi.fn().mockResolvedValue({
      ok: true,
      status: 204,
      json: async () => {
        throw new Error("sin cuerpo");
      },
    });
    expect(await crearClienteApi({ fetchFn })("/x", { method: "DELETE" })).toBeNull();
  });
});