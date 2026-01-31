const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

export function getToken() {
  return localStorage.getItem("accessToken");
}

export async function apiFetch(path, { method = "GET", body, auth = false } = {}) {
  const headers = { "Content-Type": "application/json" };
  if (auth) {
    const t = getToken();
    if (t) headers["Authorization"] = `Bearer ${t}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const data = await res.json().catch(() => null);

  if (!res.ok) {
    const msg = data?.message || data?.error?.message || "Request failed";
    const code = data?.error?.code || "error";
    const err = new Error(msg);
    err.status = res.status;
    err.code = code;
    err.payload = data;
    throw err;
  }

  return data;
}
