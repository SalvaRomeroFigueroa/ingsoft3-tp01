import { describe, it, expect } from "vitest";
import { validarTitulo, accionesPara, formatearVencimiento, escaparHtml, MAX_TITULO } from "../lib/tareas.js";

describe("validarTitulo", () => {
  // PARAMETRIZADO: una tabla de casos, incluidos los bordes del largo.
  it.each([
    { nombre: "título normal", titulo: "Comprar pan", ok: true },
    { nombre: "exactamente el máximo", titulo: "a".repeat(MAX_TITULO), ok: true },
    { nombre: "un caracter de más", titulo: "a".repeat(MAX_TITULO + 1), ok: false },
    { nombre: "solo espacios", titulo: "   ", ok: false },
    { nombre: "vacío", titulo: "", ok: false },
    { nombre: "sin valor (undefined)", titulo: undefined, ok: false },
  ])("$nombre → ok=$ok", ({ titulo, ok }) => {
    expect(validarTitulo(titulo).ok).toBe(ok);
  });

  // CASO DE ERROR: el mensaje dice por qué se rechazó.
  it("explica el motivo del rechazo", () => {
    expect(validarTitulo("   ").error).toMatch(/vacío/);
    expect(validarTitulo("a".repeat(MAX_TITULO + 1)).error).toMatch(/100/);
  });
});

describe("accionesPara", () => {
  it.each([
    ["pendiente", ["completar", "archivar", "eliminar"]],
    ["completada", ["reabrir", "archivar", "eliminar"]],
    ["archivada", ["eliminar"]],
    ["inventado", []],
  ])("estado %s ofrece %j", (estado, esperadas) => {
    expect(accionesPara(estado)).toEqual(esperadas);
  });

  it("una tarea archivada no puede volver a abrirse", () => {
    expect(accionesPara("archivada")).not.toContain("reabrir");
  });
});

describe("formatearVencimiento", () => {
  it("recorta la fecha ISO a YYYY-MM-DD", () => {
    expect(formatearVencimiento("2030-05-01T00:00:00Z")).toBe("2030-05-01");
  });
  it("devuelve cadena vacía si no hay fecha", () => {
    expect(formatearVencimiento(null)).toBe("");
  });
});

describe("escaparHtml", () => {
  it("neutraliza etiquetas y comillas", () => {
    expect(escaparHtml('<b>"x" & \'y\'</b>')).toBe("&lt;b&gt;&quot;x&quot; &amp; &#39;y&#39;&lt;/b&gt;");
  });
});
