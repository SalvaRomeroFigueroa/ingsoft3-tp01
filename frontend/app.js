import { crearClienteApi } from "./lib/api.js";
import { validarTitulo, accionesPara, formatearVencimiento, escaparHtml } from "./lib/tareas.js";

// ---------------------------------------------------------------------------
// Frontend vanilla (sin frameworks) para la To-Do List.
// Toda la comunicación con el backend pasa por api(), que arma la URL con
// API_BASE (ver index.html) y agrega el header X-User-Id cuando hay un
// usuario logueado.
// ---------------------------------------------------------------------------

const API_BASE = window.API_BASE || "";

let currentUser = null; // { id, username }
let currentListId = null;
let lists = [];

// --- Helper de fetch ---------------------------------------------------

const api = crearClienteApi({
  fetchFn: (url, opciones) => fetch(url, opciones),
  baseUrl: API_BASE,
  getUserId: () => (currentUser ? currentUser.id : null),
});

// --- Elementos del DOM ---------------------------------------------------

const loginScreen = document.getElementById("login-screen");
const mainScreen = document.getElementById("main-screen");
const loginForm = document.getElementById("login-form");
const usernameInput = document.getElementById("username-input");
const loginError = document.getElementById("login-error");
const currentUsernameEl = document.getElementById("current-username");
const logoutBtn = document.getElementById("logout-btn");

const listsContainer = document.getElementById("lists-container");
const newListForm = document.getElementById("new-list-form");
const newListInput = document.getElementById("new-list-input");
const listError = document.getElementById("list-error");

const noListSelected = document.getElementById("no-list-selected");
const taskView = document.getElementById("task-view");
const currentListName = document.getElementById("current-list-name");
const deleteListBtn = document.getElementById("delete-list-btn");

const newTaskForm = document.getElementById("new-task-form");
const newTaskTitle = document.getElementById("new-task-title");
const newTaskDue = document.getElementById("new-task-due");
const taskError = document.getElementById("task-error");
const tasksContainer = document.getElementById("tasks-container");

// --- Utilidades de UI ---------------------------------------------------

function showError(el, message) {
  el.textContent = message;
  el.classList.remove("hidden");
}

function hideError(el) {
  el.classList.add("hidden");
}

function saveSession(user) {
  localStorage.setItem("todo_user", JSON.stringify(user));
}

function loadSession() {
  const raw = localStorage.getItem("todo_user");
  return raw ? JSON.parse(raw) : null;
}

function clearSession() {
  localStorage.removeItem("todo_user");
}

// --- Login ---------------------------------------------------------------

loginForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  hideError(loginError);
  const username = usernameInput.value.trim();
  if (!username) return;

  try {
    const user = await api("/api/users/login", {
      method: "POST",
      body: JSON.stringify({ username }),
    });
    currentUser = user;
    saveSession(user);
    enterApp();
  } catch (err) {
    showError(loginError, err.message);
  }
});

logoutBtn.addEventListener("click", () => {
  currentUser = null;
  currentListId = null;
  clearSession();
  mainScreen.classList.add("hidden");
  loginScreen.classList.remove("hidden");
  usernameInput.value = "";
});

function enterApp() {
  currentUsernameEl.textContent = currentUser.username;
  loginScreen.classList.add("hidden");
  mainScreen.classList.remove("hidden");
  refreshLists();
}

// --- Listas ---------------------------------------------------------------

async function refreshLists() {
  try {
    lists = await api("/api/lists");
    renderLists();
  } catch (err) {
    showError(listError, err.message);
  }
}

function renderLists() {
  listsContainer.innerHTML = "";
  lists.forEach((list) => {
    const li = document.createElement("li");
    li.className = list.id === currentListId ? "active" : "";
    li.innerHTML = `
      <span>${escaparHtml(list.name)}</span>
      ${list.pending_count > 0 ? `<span class="badge">${list.pending_count}</span>` : ""}
    `;
    li.addEventListener("click", () => selectList(list.id));
    listsContainer.appendChild(li);
  });
}

newListForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  hideError(listError);
  const name = newListInput.value.trim();
  if (!name) return;

  try {
    const list = await api("/api/lists", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
    newListInput.value = "";
    await refreshLists();
    selectList(list.id);
  } catch (err) {
    showError(listError, err.message);
  }
});

deleteListBtn.addEventListener("click", async () => {
  if (currentListId === null) return;
  if (!confirm("¿Eliminar esta lista?")) return;

  try {
    await api(`/api/lists/${currentListId}`, { method: "DELETE" });
    currentListId = null;
    noListSelected.classList.remove("hidden");
    taskView.classList.add("hidden");
    await refreshLists();
  } catch (err) {
    showError(taskError, err.message);
  }
});

async function selectList(listId) {
  currentListId = listId;
  renderLists();

  const list = lists.find((l) => l.id === listId);
  currentListName.textContent = list ? list.name : "";

  noListSelected.classList.add("hidden");
  taskView.classList.remove("hidden");

  await refreshTasks();
}

// --- Tareas ---------------------------------------------------------------

async function refreshTasks() {
  if (currentListId === null) return;
  try {
    const tasks = await api(`/api/lists/${currentListId}/tasks`);
    renderTasks(tasks);
  } catch (err) {
    showError(taskError, err.message);
  }
}

function renderTasks(tasks) {
  tasksContainer.innerHTML = "";
  tasks.forEach((task) => {
    const li = document.createElement("li");
    li.className = `task-item ${task.status}`;

    const venc = formatearVencimiento(task.due_date);
    const dueLabel = venc ? `<span class="task-due">Vence: ${venc}</span>` : "";

    li.innerHTML = `
      <span class="task-title">${escaparHtml(task.title)}</span>
      ${dueLabel}
      <span class="task-status ${task.status}">${task.status}</span>
      <span class="task-actions"></span>
    `;

    const actions = li.querySelector(".task-actions");
    addTaskActions(actions, task);

    tasksContainer.appendChild(li);
  });
}

function addTaskActions(container, task) {
  const botones = {
    completar: () => makeButton("Completar", "icon-btn", () => updateTaskStatus(task.id, "completada")),
    reabrir: () => makeButton("Reabrir", "icon-btn", () => updateTaskStatus(task.id, "pendiente")),
    archivar: () => makeButton("Archivar", "icon-btn", () => updateTaskStatus(task.id, "archivada")),
    eliminar: () => makeButton("Eliminar", "danger-btn", () => deleteTask(task.id)),
  };
  accionesPara(task.status).forEach((nombre) => container.appendChild(botones[nombre]()));
}

function makeButton(label, className, onClick) {
  const btn = document.createElement("button");
  btn.textContent = label;
  btn.className = className;
  btn.addEventListener("click", onClick);
  return btn;
}

newTaskForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  hideError(taskError);
  const title = newTaskTitle.value.trim();
  const chequeo = validarTitulo(title);
  if (!chequeo.ok) {
    showError(taskError, chequeo.error);
    return;
  }

  try {
    await api(`/api/lists/${currentListId}/tasks`, {
      method: "POST",
      body: JSON.stringify({
        title,
        due_date: newTaskDue.value || null,
      }),
    });
    newTaskTitle.value = "";
    newTaskDue.value = "";
    await refreshTasks();
    await refreshLists();
  } catch (err) {
    showError(taskError, err.message);
  }
});

async function updateTaskStatus(taskId, status) {
  hideError(taskError);
  try {
    await api(`/api/tasks/${taskId}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    });
    await refreshTasks();
    await refreshLists();
  } catch (err) {
    showError(taskError, err.message);
  }
}

async function deleteTask(taskId) {
  hideError(taskError);
  try {
    await api(`/api/tasks/${taskId}`, { method: "DELETE" });
    await refreshTasks();
    await refreshLists();
  } catch (err) {
    showError(taskError, err.message);
  }
}

// --- Utilidades ---------------------------------------------------------

// --- Arranque ---------------------------------------------------------------

(function init() {
  const savedUser = loadSession();
  if (savedUser) {
    currentUser = savedUser;
    enterApp();
  }
})();
