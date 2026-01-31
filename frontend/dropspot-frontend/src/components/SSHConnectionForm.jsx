import { useMemo, useState } from "react";
import { createSSH } from "../api/ssh";

export default function SSHConnectionForm({ onCreated }) {
  const [name, setName] = useState("");
  const [host, setHost] = useState("");
  const [port, setPort] = useState(22);
  const [username, setUsername] = useState("");
  const [authType, setAuthType] = useState("password");
  const [secret, setSecret] = useState("");

  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState("");

  const canSubmit = useMemo(() => {
    const p = Number(port);
    return (
      name.trim() &&
      host.trim() &&
      username.trim() &&
      secret.trim() &&
      Number.isFinite(p) &&
      p >= 1 &&
      p <= 65535
    );
  }, [name, host, port, username, secret]);

  async function submit(e) {
    e.preventDefault();
    setErr("");
    if (!canSubmit) {
      setErr("Please fill all fields correctly.");
      return;
    }
    setBusy(true);
    try {
      await createSSH({
        name: name.trim(),
        host: host.trim(),
        port: Number(port),
        username: username.trim(),
        authType,
        secret: secret.trim(),
      });

      setName("");
      setHost("");
      setPort(22);
      setUsername("");
      setAuthType("password");
      setSecret("");
      onCreated?.();
    } catch (e) {
      setErr(e.message || "Request failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="card">
      <h3 className="cardTitle">Add SSH Connection</h3>

      <form className="form" onSubmit={submit}>
        <input className="input" placeholder="Name (e.g. ubuntu-vm)" value={name} onChange={(e) => setName(e.target.value)} disabled={busy} />
        <input className="input" placeholder="Host (IP or domain)" value={host} onChange={(e) => setHost(e.target.value)} disabled={busy} />

        <div className="row2">
          <input className="input" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} disabled={busy} />
          <input className="input" placeholder="Port" type="number" min={1} max={65535} value={port} onChange={(e) => setPort(e.target.value)} disabled={busy} />
        </div>

        <select className="select" value={authType} onChange={(e) => { setAuthType(e.target.value); setSecret(""); }} disabled={busy}>
          <option value="password">Password</option>
          <option value="private_key">Private Key</option>
        </select>

        <textarea
          className="textarea"
          placeholder={authType === "password" ? "Password" : "Private key (PEM)"}
          value={secret}
          onChange={(e) => setSecret(e.target.value)}
          rows={authType === "password" ? 3 : 7}
          disabled={busy}
          style={{ fontFamily: authType === "private_key" ? "monospace" : "inherit" }}
        />

        {err && <div className="alertErr">{err}</div>}

        <button className="btn btnPrimary" disabled={busy || !canSubmit}>
          {busy ? "Saving..." : "Save connection"}
        </button>
      </form>

      <p className="helper">Secret backend’de AES-GCM ile şifrelenerek saklanır.</p>
    </div>
  );
}
