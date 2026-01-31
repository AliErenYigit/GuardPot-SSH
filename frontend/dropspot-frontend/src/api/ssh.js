import { apiFetch } from "./http";

export function listSSH() {
  return apiFetch("/api/v1/ssh", { auth: true });
}

export function createSSH(payload) {
  return apiFetch("/api/v1/ssh", { method: "POST", body: payload, auth: true });
}

export function deleteSSH(id) {
  return apiFetch(`/api/v1/ssh/${id}`, { method: "DELETE", auth: true });
}

export function scanHostKey(id) {
  return apiFetch(`/api/v1/ssh/${id}/hostkey/scan`, { method: "POST", auth: true });
}

export function trustHostKey(id) {
  return apiFetch(`/api/v1/ssh/${id}/hostkey/trust`, { method: "POST", auth: true });
}

export function listKnownHosts() {
  return apiFetch(`/api/v1/ssh/known-hosts`, { auth: true });
}
