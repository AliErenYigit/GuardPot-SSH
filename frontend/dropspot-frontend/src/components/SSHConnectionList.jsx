import { useEffect, useState } from "react";
import { listSSH, deleteSSH } from "../api/ssh";
import TerminalModal from "./TerminalModal";

export default function SSHConnectionList({ refreshKey }) {
  const [items, setItems] = useState([]);
  const [busyId, setBusyId] = useState(null);
  const [err, setErr] = useState("");
  const [info, setInfo] = useState("");

  const [terminalConnId, setTerminalConnId] = useState(null);

  async function load() {
    setErr("");
    try {
      const res = await listSSH();
      setItems(res.data || []);
    } catch (e) {
      setErr(e.message || "Failed to load connections");
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [refreshKey]);

  // async function onScan(id) {
  //   setBusyId(id);
  //   setErr("");
  //   setInfo("");
  //   try {
  //     const res = await scanHostKey(id);
  //     setInfo(`Scan OK: ${res.data.fingerprint}`);
  //   } catch (e) {
  //     setErr(e.message || "Scan failed");
  //   } finally {
  //     setBusyId(null);
  //   }
  // }

  // async function onTrust(id) {
  //   setBusyId(id);
  //   setErr("");
  //   setInfo("");
  //   try {
  //     const res = await trustHostKey(id);
  //     setInfo(`Trusted: ${res.data.fingerprint}`);
  //     // İstersen trust sonrası listeyi yenile (güncel durum vs için)
  //     // await load();
  //   } catch (e) {
  //     setErr(e.message || "Trust failed");
  //   } finally {
  //     setBusyId(null);
  //   }
  // }

  async function onDelete(id) {
    setBusyId(id);
    setErr("");
    setInfo("");
    try {
      await deleteSSH(id);
      await load();
    } catch (e) {
      setErr(e.message || "Delete failed");
    } finally {
      setBusyId(null);
    }
  }

  return (
    <div className="card">
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12 }}>
        <h3 className="cardTitle" style={{ marginBottom: 0 }}>SSH Connections</h3>
        <button className="btn" onClick={load}>Refresh</button>
      </div>

      {err && <div className="alertErr">{err}</div>}
      {info && <div className="alertOk">{info}</div>}

      {items.length === 0 ? (
        <div className="sub" style={{ marginTop: 10 }}>No SSH connections</div>
      ) : (
        <div className="tableWrap" style={{ marginTop: 10 }}>
          <table className="table">
            <thead>
              <tr>
                <th style={{ width: 220 }}>Name</th>
                <th style={{ width: 260 }}>Host</th>
                <th style={{ width: 180 }}>User</th>
                <th>Actions</th>
              </tr>
            </thead>

            <tbody>
              {items.map((c) => (
                <tr key={c.id}>
                  <td>
                    <div style={{ fontWeight: 700 }}>{c.name}</div>
                    <div style={{ marginTop: 4 }}>
                      <span className="badge">{c.authType}</span>
                    </div>
                  </td>

                  <td>
                    <div style={{ fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace" }}>
                      {c.host}:{c.port}
                    </div>
                  </td>

                  <td>{c.username}</td>

                  <td>
                    <div className="actions">
                      {/* <button className="btn" onClick={() => onScan(c.id)} disabled={busyId === c.id}>
                        {busyId === c.id ? "..." : "Scan"}
                      </button> */}

                      {/* <button className="btn" onClick={() => onTrust(c.id)} disabled={busyId === c.id}>
                        {busyId === c.id ? "..." : "Trust"}
                      </button> */}

                      <button className="btn btnPrimary" onClick={() => setTerminalConnId(c.id)}>
                        Open Terminal
                      </button>

                      <button className="btn btnDanger" onClick={() => onDelete(c.id)} disabled={busyId === c.id}>
                        {busyId === c.id ? "..." : "Delete"}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <TerminalModal
        open={terminalConnId !== null}
        connId={terminalConnId}
        onClose={() => setTerminalConnId(null)}
      />
    </div>
  );
}
