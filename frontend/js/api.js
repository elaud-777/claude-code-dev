(function () {
  const API_BASE = window.API_BASE || "http://localhost:8000";

  function getToken() {
    return localStorage.getItem("token");
  }

  function setToken(token) {
    localStorage.setItem("token", token);
  }

  function clearToken() {
    localStorage.removeItem("token");
  }

  class ApiError extends Error {
    constructor(status, code, message, meta) {
      super(message);
      this.status = status;
      this.code = code;
      this.meta = meta;
    }
  }

  async function apiFetch(path, { method = "GET", body, auth = true } = {}) {
    const headers = { "Content-Type": "application/json" };
    if (auth) {
      const token = getToken();
      if (token) headers["Authorization"] = `Bearer ${token}`;
    }

    const res = await fetch(`${API_BASE}${path}`, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    let data = null;
    try {
      data = await res.json();
    } catch (_) {
      data = null;
    }

    if (!res.ok) {
      const err = data && data.error ? data.error : { code: "UNKNOWN", message: "알 수 없는 오류가 발생했습니다" };
      if (err.code === "TOKEN_EXPIRED" || res.status === 401) {
        clearToken();
        if (window.onAuthExpired) window.onAuthExpired();
      }
      throw new ApiError(res.status, err.code, err.message, err.meta);
    }

    return data;
  }

  window.api = { apiFetch, getToken, setToken, clearToken, ApiError };
})();
