import { apiFetch } from "./http";

export function register(email, password) {
  return apiFetch("/api/v1/auth/register", {
    method: "POST",
    body: { email, password },
  });
}

export function login(email, password) {
  return apiFetch("/api/v1/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export function me() {
  return apiFetch("/api/v1/me", { auth: true });
}

// ✅ Google login: idToken backend'e gider, backend accessToken döner


export function loginWithGoogle(idToken) {
  return apiFetch("/api/v1/auth/google", {
    method: "POST",
    body: { idToken },
  });
}


