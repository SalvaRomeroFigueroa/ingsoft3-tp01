// Prioridad de una tarea según su estado y su fecha de vencimiento.
// (Tiene varios caminos adentro y, a propósito, ni un solo test.)
const MS_POR_DIA = 24 * 60 * 60 * 1000;

export function prioridadDe(tarea, ahora = new Date()) {
  if (!tarea) {
    return "ninguna";
  }
  if (tarea.status === "archivada" || tarea.status === "completada") {
    return "ninguna";
  }
  if (!tarea.due_date) {
    return "sin-fecha";
  }
  const vence = new Date(String(tarea.due_date).substring(0, 10));
  const hoy = new Date(ahora.toISOString().substring(0, 10));
  const dias = Math.round((vence - hoy) / MS_POR_DIA);
  if (dias < 0) {
    return "vencida";
  }
  if (dias === 0) {
    return "hoy";
  }
  if (dias <= 7) {
    return "esta-semana";
  }
  if (dias <= 30) {
    return "este-mes";
  }
  return "mas-adelante";
}

// Cuenta cuántas tareas hay en cada nivel de prioridad.
export function resumirPrioridades(tareas, ahora = new Date()) {
  const resumen = {
    vencida: 0,
    hoy: 0,
    "esta-semana": 0,
    "este-mes": 0,
    "mas-adelante": 0,
    "sin-fecha": 0,
    ninguna: 0,
  };
  for (const tarea of tareas) {
    resumen[prioridadDe(tarea, ahora)]++;
  }
  return resumen;
}