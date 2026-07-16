const { apiFetch, getToken, setToken, clearToken } = window.api;

const state = {
  user: null,
  team: null,
  members: [],
  tasks: [],
  filter: "all",
  lastMessageAt: null,
  chatPollTimer: null,
  isSignup: false,
  currentTaskId: null,
  pendingConfirm: null,
};

// ---------- View switching ----------
function showView(id) {
  ["view-auth", "view-team-select", "view-app"].forEach((v) => {
    document.getElementById(v).classList.toggle("hidden", v !== id);
  });
}

function showToast(message) {
  const toast = document.getElementById("toast");
  toast.textContent = message;
  toast.classList.remove("hidden");
  setTimeout(() => toast.classList.add("hidden"), 2500);
}

window.onAuthExpired = () => {
  stopChatPolling();
  showToast("인증이 만료되었습니다. 다시 로그인해주세요.");
  showView("view-auth");
};

// ---------- Auth ----------
const authForm = document.getElementById("auth-form");
const authTitle = document.getElementById("auth-title");
const authToggleLink = document.getElementById("auth-toggle-link");
const authToggleText = document.getElementById("auth-toggle-text");
const authSubmit = document.getElementById("auth-submit");

function resetAuthErrors() {
  ["auth-email-error", "auth-password-error", "auth-general-error"].forEach((id) => {
    const el = document.getElementById(id);
    el.classList.add("hidden");
    el.textContent = "";
  });
}

authToggleLink.addEventListener("click", (e) => {
  e.preventDefault();
  state.isSignup = !state.isSignup;
  authTitle.textContent = state.isSignup ? "회원가입" : "로그인";
  authSubmit.textContent = state.isSignup ? "가입하기" : "로그인";
  authToggleText.textContent = state.isSignup ? "이미 계정이 있으신가요?" : "계정이 없으신가요?";
  authToggleLink.textContent = state.isSignup ? "로그인" : "회원가입";
  resetAuthErrors();
});

authForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  resetAuthErrors();
  const email = document.getElementById("auth-email").value.trim();
  const password = document.getElementById("auth-password").value;

  authSubmit.disabled = true;
  authSubmit.textContent = "처리 중…";

  try {
    const path = state.isSignup ? "/auth/signup" : "/auth/login";
    const data = await apiFetch(path, { method: "POST", body: { email, password }, auth: false });
    setToken(data.token);
    state.user = data.user;
    await afterLogin();
  } catch (err) {
    const el = document.getElementById("auth-general-error");
    el.textContent = err.message || "오류가 발생했습니다";
    el.classList.remove("hidden");
  } finally {
    authSubmit.disabled = false;
    authSubmit.textContent = state.isSignup ? "가입하기" : "로그인";
  }
});

async function afterLogin() {
  if (state.user.team_id) {
    await enterTeam(state.user.team_id);
  } else {
    document.getElementById("team-select-email").textContent = state.user.email;
    showView("view-team-select");
  }
}

async function tryRestoreSession() {
  const token = getToken();
  if (!token) {
    showView("view-auth");
    return;
  }
  try {
    state.user = await apiFetch("/auth/me");
    await afterLogin();
  } catch (_) {
    clearToken();
    showView("view-auth");
  }
}

function logout() {
  apiFetch("/auth/logout", { method: "POST" }).catch(() => {});
  clearToken();
  stopChatPolling();
  state.user = null;
  state.team = null;
  showView("view-auth");
}

document.getElementById("team-select-logout").addEventListener("click", (e) => { e.preventDefault(); logout(); });
document.getElementById("app-logout").addEventListener("click", (e) => { e.preventDefault(); logout(); });
document.getElementById("mobile-logout").addEventListener("click", (e) => { e.preventDefault(); logout(); });

// ---------- Team select ----------
document.getElementById("team-create-btn").addEventListener("click", async () => {
  const name = document.getElementById("team-name-input").value.trim();
  const errEl = document.getElementById("team-create-error");
  errEl.classList.add("hidden");
  try {
    const team = await apiFetch("/teams", { method: "POST", body: { name } });
    state.user.team_id = team.id;
    await enterTeam(team.id);
  } catch (err) {
    errEl.textContent = err.message;
    errEl.classList.remove("hidden");
  }
});

document.getElementById("team-join-btn").addEventListener("click", async () => {
  const code = document.getElementById("invite-code-input").value.trim().toUpperCase();
  const errEl = document.getElementById("team-join-error");
  errEl.classList.add("hidden");
  try {
    const info = await apiFetch("/teams/join", { method: "POST", body: { invite_code: code } });
    state.user.team_id = info.id;
    await enterTeam(info.id);
  } catch (err) {
    errEl.textContent = err.message;
    errEl.classList.remove("hidden");
  }
});

async function enterTeam(teamId) {
  state.team = await apiFetch(`/teams/${teamId}`);
  document.getElementById("app-team-name").textContent = state.team.name;
  document.getElementById("app-email").textContent = state.user.email;
  showView("view-app");
  switchTab("kanban");
  await Promise.all([loadTasks(), loadMembers()]);
  startChatPolling();
}

// ---------- Tabs / mobile menu ----------
function switchTab(tab) {
  document.querySelectorAll(".tab-panel").forEach((p) => p.classList.add("hidden"));
  document.getElementById(`tab-${tab}`).classList.remove("hidden");
  document.querySelectorAll(".tab-btn").forEach((b) => {
    b.classList.toggle("bg-teal-700", b.dataset.tab === tab);
    b.classList.toggle("text-white", b.dataset.tab === tab);
    b.classList.toggle("bg-gray-100", b.dataset.tab !== tab);
  });
  document.getElementById("mobile-menu").classList.add("hidden");
  if (tab === "chat") scrollChatToBottom();
}

document.querySelectorAll(".tab-btn, .tab-btn-mobile").forEach((btn) => {
  btn.addEventListener("click", () => switchTab(btn.dataset.tab));
});

document.getElementById("hamburger-btn").addEventListener("click", () => {
  document.getElementById("mobile-menu").classList.toggle("hidden");
});

// ---------- Kanban ----------
document.querySelectorAll(".filter-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    state.filter = btn.dataset.filter;
    document.querySelectorAll(".filter-btn").forEach((b) => {
      b.classList.toggle("bg-gray-800", b === btn);
      b.classList.toggle("text-white", b === btn);
      b.classList.toggle("bg-gray-100", b !== btn);
    });
    loadTasks();
  });
});

async function loadTasks() {
  const query = state.filter !== "all" ? `?filter=${state.filter}` : "";
  state.tasks = await apiFetch(`/teams/${state.team.id}/tasks${query}`);
  renderTasks();
}

function renderTasks() {
  const columns = { TODO: [], DOING: [], DONE: [] };
  state.tasks.forEach((t) => columns[t.status].push(t));

  document.querySelectorAll(".kanban-column").forEach((col) => {
    const status = col.dataset.status;
    col.querySelector(".count").textContent = columns[status].length;
    const list = col.querySelector(".task-list");
    list.innerHTML = "";
    columns[status].forEach((task) => {
      const card = document.createElement("div");
      card.className = "task-card bg-white border rounded p-2 cursor-pointer shadow-sm";
      card.draggable = true;
      card.dataset.taskId = task.id;
      const assignee = memberEmail(task.assignee_id);
      card.innerHTML = `
        <div class="font-medium text-sm">${escapeHtml(task.title)}</div>
        <div class="text-xs text-gray-400">#${task.id} · ${assignee ? "@" + assignee : "⚠미할당"}</div>
      `;
      card.addEventListener("click", () => openTaskModal(task.id));
      card.addEventListener("dragstart", (e) => {
        e.dataTransfer.setData("text/plain", task.id);
      });
      list.appendChild(card);
    });
  });
}

function memberEmail(userId) {
  if (!userId) return null;
  const m = state.members.find((m) => m.id === userId);
  return m ? m.email.split("@")[0] : null;
}

function escapeHtml(s) {
  const div = document.createElement("div");
  div.textContent = s;
  return div.innerHTML;
}

document.querySelectorAll(".kanban-column").forEach((col) => {
  col.addEventListener("dragover", (e) => e.preventDefault());
  col.addEventListener("drop", async (e) => {
    e.preventDefault();
    const taskId = e.dataTransfer.getData("text/plain");
    const newStatus = col.dataset.status;
    try {
      await apiFetch(`/tasks/${taskId}/status`, { method: "PATCH", body: { status: newStatus } });
      await loadTasks();
    } catch (err) {
      showToast(err.message);
    }
  });
});

document.querySelectorAll(".add-task-btn").forEach((btn) => {
  btn.addEventListener("click", async () => {
    const title = prompt("태스크 제목을 입력하세요 (1-100자)");
    if (!title) return;
    try {
      await apiFetch(`/teams/${state.team.id}/tasks`, { method: "POST", body: { title } });
      await loadTasks();
    } catch (err) {
      showToast(err.message);
    }
  });
});

// ---------- Task modal ----------
const taskModal = document.getElementById("task-modal");

function populateAssigneeOptions(selectedId) {
  const select = document.getElementById("modal-assignee-select");
  select.innerHTML = '<option value="">미할당</option>';
  state.members.forEach((m) => {
    const opt = document.createElement("option");
    opt.value = m.id;
    opt.textContent = m.email;
    if (m.id === selectedId) opt.selected = true;
    select.appendChild(opt);
  });
}

function openTaskModal(taskId) {
  const task = state.tasks.find((t) => t.id === Number(taskId));
  if (!task) return;
  state.currentTaskId = task.id;
  document.getElementById("modal-task-id").textContent = `#${task.id}`;
  document.getElementById("modal-task-title-display").textContent = task.title;
  document.getElementById("modal-title-input").value = task.title;
  populateAssigneeOptions(task.assignee_id);
  document.querySelectorAll(".modal-status-btn").forEach((b) => {
    b.classList.toggle("bg-teal-700", b.dataset.status === task.status);
    b.classList.toggle("text-white", b.dataset.status === task.status);
  });
  document.getElementById("modal-error").classList.add("hidden");
  taskModal.classList.remove("hidden");
  taskModal.dataset.selectedStatus = task.status;
}

document.getElementById("modal-close-btn").addEventListener("click", () => taskModal.classList.add("hidden"));

document.querySelectorAll(".modal-status-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    taskModal.dataset.selectedStatus = btn.dataset.status;
    document.querySelectorAll(".modal-status-btn").forEach((b) => {
      b.classList.toggle("bg-teal-700", b === btn);
      b.classList.toggle("text-white", b === btn);
    });
  });
});

document.getElementById("modal-save-btn").addEventListener("click", async () => {
  const errEl = document.getElementById("modal-error");
  errEl.classList.add("hidden");
  const title = document.getElementById("modal-title-input").value.trim();
  const assigneeVal = document.getElementById("modal-assignee-select").value;
  const newStatus = taskModal.dataset.selectedStatus;
  const task = state.tasks.find((t) => t.id === state.currentTaskId);

  try {
    await apiFetch(`/tasks/${state.currentTaskId}`, {
      method: "PUT",
      body: { title, assignee_id: assigneeVal ? Number(assigneeVal) : null },
    });
    if (newStatus !== task.status) {
      await apiFetch(`/tasks/${state.currentTaskId}/status`, { method: "PATCH", body: { status: newStatus } });
    }
    taskModal.classList.add("hidden");
    await loadTasks();
  } catch (err) {
    errEl.textContent = err.message;
    errEl.classList.remove("hidden");
  }
});

document.getElementById("modal-delete-btn").addEventListener("click", () => {
  const task = state.tasks.find((t) => t.id === state.currentTaskId);
  askConfirm(`'#${task.id} ${task.title}' — 되돌릴 수 없습니다`, async () => {
    try {
      await apiFetch(`/tasks/${state.currentTaskId}`, { method: "DELETE" });
      taskModal.classList.add("hidden");
      await loadTasks();
    } catch (err) {
      showToast(err.message);
    }
  });
});

// ---------- Confirm dialog ----------
const confirmDialog = document.getElementById("confirm-dialog");

function askConfirm(message, onConfirm) {
  document.getElementById("confirm-message").textContent = message;
  state.pendingConfirm = onConfirm;
  confirmDialog.classList.remove("hidden");
}

document.getElementById("confirm-cancel-btn").addEventListener("click", () => {
  confirmDialog.classList.add("hidden");
  state.pendingConfirm = null;
});

document.getElementById("confirm-ok-btn").addEventListener("click", async () => {
  confirmDialog.classList.add("hidden");
  if (state.pendingConfirm) await state.pendingConfirm();
  state.pendingConfirm = null;
});

// ---------- Members ----------
async function loadMembers() {
  state.members = await apiFetch(`/teams/${state.team.id}/members`);
  document.getElementById("member-count").textContent = state.members.length;
  const list = document.getElementById("member-list");
  list.innerHTML = "";
  state.members.forEach((m) => {
    const row = document.createElement("div");
    row.className = "flex justify-between items-center border rounded px-3 py-2";
    row.innerHTML = `
      <div>
        <div class="text-sm font-medium">${escapeHtml(m.email)}${m.is_owner ? ' <span class="text-teal-700">★ owner</span>' : ""}</div>
        <div class="text-xs text-gray-400">${m.created_at.slice(0, 10)}</div>
      </div>
    `;
    list.appendChild(row);
  });
}

// ---------- Chat ----------
const chatMessagesEl = document.getElementById("chat-messages");
const chatInput = document.getElementById("chat-input");
const chatCounter = document.getElementById("chat-counter");

chatInput.addEventListener("input", () => {
  const len = chatInput.value.length;
  chatCounter.textContent = `${len} / 1000`;
  chatCounter.classList.toggle("text-red-600", len > 1000);
  document.getElementById("chat-send-btn").disabled = len > 1000 || len === 0;
});

document.getElementById("chat-send-btn").addEventListener("click", sendMessage);
chatInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter") sendMessage();
});

async function sendMessage() {
  const content = chatInput.value;
  if (!content || content.length > 1000) return;
  try {
    await apiFetch(`/teams/${state.team.id}/messages`, { method: "POST", body: { content } });
    chatInput.value = "";
    chatCounter.textContent = "0 / 1000";
    await pollMessages();
  } catch (err) {
    showToast(err.message);
  }
}

function renderMessage(m) {
  const isMe = m.user_id === state.user.id;
  const wrapper = document.createElement("div");
  wrapper.className = `flex ${isMe ? "justify-end" : "justify-start"}`;
  const time = new Date(m.created_at).toLocaleTimeString("ko-KR", { hour: "2-digit", minute: "2-digit" });
  const deleteBtn = isMe
    ? `<button class="delete-msg-btn text-xs text-red-500 ml-2" data-msg-id="${m.id}">🗑</button>`
    : "";
  wrapper.innerHTML = `
    <div class="max-w-xs">
      <div class="text-xs text-gray-400 mb-1">${isMe ? "" : escapeHtml(m.user_email.split("@")[0]) + " · "}${time}${deleteBtn}</div>
      <div class="${isMe ? "bg-teal-700 text-white" : "bg-gray-100"} rounded px-3 py-2 text-sm break-words">${escapeHtml(m.content)}</div>
    </div>
  `;
  chatMessagesEl.appendChild(wrapper);
  wrapper.querySelectorAll(".delete-msg-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      try {
        await apiFetch(`/messages/${btn.dataset.msgId}`, { method: "DELETE" });
        wrapper.remove();
      } catch (err) {
        showToast(err.message);
      }
    });
  });
}

function scrollChatToBottom() {
  chatMessagesEl.scrollTop = chatMessagesEl.scrollHeight;
}

async function pollMessages() {
  try {
    const since = state.lastMessageAt ? `?since=${encodeURIComponent(state.lastMessageAt)}` : "";
    const messages = await apiFetch(`/teams/${state.team.id}/messages${since}`);
    if (messages.length === 0 && state.lastMessageAt) return;
    if (!state.lastMessageAt) chatMessagesEl.innerHTML = "";
    messages.forEach((m) => {
      renderMessage(m);
      state.lastMessageAt = m.created_at;
    });
    if (messages.length > 0) scrollChatToBottom();
    document.getElementById("chat-status").textContent = "● 5초마다 새로고침";
  } catch (err) {
    document.getElementById("chat-status").textContent = "⚠ 연결 끊김 · 재시도 중";
  }
}

function startChatPolling() {
  state.lastMessageAt = null;
  pollMessages();
  stopChatPolling();
  state.chatPollTimer = setInterval(pollMessages, 5000);
}

function stopChatPolling() {
  if (state.chatPollTimer) clearInterval(state.chatPollTimer);
  state.chatPollTimer = null;
}

// ---------- Init ----------
tryRestoreSession();
