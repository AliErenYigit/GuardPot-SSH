import { useEffect, useMemo, useRef, useState } from "react";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

function toWSBase(apiBase) {
  return apiBase.replace(/^http/, "ws");
}

function isMac() {
  return /Mac|iPhone|iPad|iPod/i.test(navigator.platform);
}

async function safeWriteClipboard(text) {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch (err) {
    console.log(err);
  }
  return false;
}

async function safeReadClipboard() {
  try {
    if (navigator.clipboard?.readText) {
      return await navigator.clipboard.readText();
    }
  } catch (err) {
    console.log(err);
  }
  return "";
}

export default function TerminalModal({ open, onClose, connId }) {
  const termElRef = useRef(null);
  const xtermRef = useRef(null);
  const fitRef = useRef(null);
  const wsRef = useRef(null);

  const [status, setStatus] = useState("idle"); // idle | connecting | connected | closed | error
  const [statusMsg, setStatusMsg] = useState("");

  const token = localStorage.getItem("accessToken") || "";
  const wsBase = useMemo(() => toWSBase(API_BASE), []);

  const connect = () => {
    const term = xtermRef.current;
    const fit = fitRef.current;
    if (!term || !fit) return;

    try {
      wsRef.current?.close();
    } catch (err) {
      console.log(err);
    }

    setStatus("connecting");
    setStatusMsg("");

    const wsUrl = `${wsBase}/api/v1/ssh/${connId}/ws?token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      setStatus("connected");
      term.writeln("\r\n[Connected]\r\n");

      // ✅ focus + fit
      fit.fit();
      term.focus();

      const dims = fit.proposeDimensions();
      if (dims?.cols && dims?.rows) {
        ws.send(JSON.stringify({ type: "resize", cols: dims.cols, rows: dims.rows }));
      }
    };

    ws.onmessage = (evt) => {
      term.write(evt.data);
    };

    ws.onerror = () => {
      setStatus("error");
      setStatusMsg("WebSocket error");
      term.writeln("\r\n[WS error]\r\n");
    };

    ws.onclose = () => {
      setStatus("closed");
      term.writeln("\r\n[Disconnected]\r\n");
    };
  };

  useEffect(() => {
    if (!open) return;

    const term = new Terminal({
      cursorBlink: true,
      scrollback: 5000,
      fontSize: 14,
      convertEol: true,
    });
    const fit = new FitAddon();
    term.loadAddon(fit);

    xtermRef.current = term;
    fitRef.current = fit;

    term.open(termElRef.current);
    fit.fit();

    // ✅ En kritik satır: terminal input için focus
    term.focus();

    setStatus("connecting");
    term.writeln("Connecting...\r\n");

    connect();

    const mac = isMac();

    const keyHandler = (e) => {
      const hasSelection = term.hasSelection?.() && term.getSelection?.();

      // COPY
      if (
        (!mac && e.ctrlKey && e.shiftKey && e.key.toLowerCase() === "c") ||
        (mac && e.metaKey && e.key.toLowerCase() === "c")
      ) {
        if (hasSelection) {
          const text = term.getSelection();
          safeWriteClipboard(text);
          e.preventDefault();
          return false;
        }
        return true; // selection yoksa Ctrl+C remote'a gitsin
      }

      // PASTE
      if (
        (!mac && e.ctrlKey && e.shiftKey && e.key.toLowerCase() === "v") ||
        (mac && e.metaKey && e.key.toLowerCase() === "v")
      ) {
        e.preventDefault();
        (async () => {
          const clip = await safeReadClipboard();
          if (clip && wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({ type: "data", data: clip }));
          }
          // paste sonrası tekrar focus
          term.focus();
        })();
        return false;
      }

      return true;
    };

    const container = termElRef.current;

    // ✅ Terminal alanı tıklanınca focus geri gelsin
    const onMouseDownFocus = () => {
      xtermRef.current?.focus();
    };
    container?.addEventListener("mousedown", onMouseDownFocus);

    const onKeyDown = (e) => keyHandler(e);
    container?.addEventListener("keydown", onKeyDown);

    const onContextMenu = (e) => {
      e.preventDefault();
      (async () => {
        const clip = await safeReadClipboard();
        if (clip && wsRef.current?.readyState === WebSocket.OPEN) {
          wsRef.current.send(JSON.stringify({ type: "data", data: clip }));
        }
        term.focus();
      })();
    };
    container?.addEventListener("contextmenu", onContextMenu);

    const disposable = term.onData((data) => {
      const ws = wsRef.current;
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: "data", data }));
      }
    });

    const onResize = () => {
      if (!fitRef.current || !wsRef.current) return;
      fitRef.current.fit();
      const dims = fitRef.current.proposeDimensions();
      if (dims?.cols && dims?.rows && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: "resize", cols: dims.cols, rows: dims.rows }));
      }
      // resize sonrası da focus kayabiliyor
      term.focus();
    };
    window.addEventListener("resize", onResize);

    return () => {
      window.removeEventListener("resize", onResize);

      try {
        disposable?.dispose();
      } catch (err) {
        console.log(err);
      }

      try {
        container?.removeEventListener("mousedown", onMouseDownFocus);
        container?.removeEventListener("keydown", onKeyDown);
        container?.removeEventListener("contextmenu", onContextMenu);
      } catch (err) {
        console.log(err);
      }

      try {
        wsRef.current?.close();
      } catch (err) {
        console.log(err);
      }

      try {
        term.dispose();
      } catch (err) {
        console.log(err);
      }

      xtermRef.current = null;
      fitRef.current = null;
      wsRef.current = null;

      setStatus("idle");
      setStatusMsg("");
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, connId]);

  if (!open) return null;

  const canReconnect = status === "closed" || status === "error";

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        background: "rgba(0,0,0,.55)",
        display: "grid",
        placeItems: "center",
        zIndex: 9999,
      }}
      onMouseDown={onClose}
    >
      <div
        style={{
          width: "min(1100px, 92vw)",
          height: "min(700px, 85vh)",
          background: "#111",
          borderRadius: 14,
          overflow: "hidden",
          border: "1px solid rgba(255,255,255,.15)",
          display: "flex",
          flexDirection: "column",
        }}
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div
          style={{
            padding: "10px 12px",
            color: "#fff",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            background: "#1a1a1a",
            borderBottom: "1px solid rgba(255,255,255,.12)",
          }}
        >
          <div style={{ display: "flex", gap: 12, alignItems: "center" }}>
            <div>Terminal — Connection #{connId}</div>
            <div style={{ fontSize: 12, opacity: 0.85 }}>
              Status: <b>{status}</b> {statusMsg ? `(${statusMsg})` : ""}
            </div>
          </div>

          <div style={{ display: "flex", gap: 8 }}>
            <button
              onClick={() => {
                try {
                  fitRef.current?.fit();
                  xtermRef.current?.focus();
                } catch (err) {
                  console.log(err);
                }
              }}
              style={{ padding: "6px 10px" }}
              title="Fit terminal to container"
            >
              Fit
            </button>

            <button
              onClick={() => {
                if (!canReconnect) return;
                xtermRef.current?.writeln("\r\n[Reconnecting...]\r\n");
                connect();
                xtermRef.current?.focus();
              }}
              disabled={!canReconnect}
              style={{ padding: "6px 10px" }}
              title="Reconnect"
            >
              Reconnect
            </button>

            <button onClick={onClose} style={{ padding: "6px 10px" }}>
              Close
            </button>
          </div>
        </div>

        <div
          style={{
            padding: "6px 12px",
            fontSize: 12,
            color: "rgba(255,255,255,.75)",
            background: "#141414",
            borderBottom: "1px solid rgba(255,255,255,.08)",
          }}
        >
          Copy: <b>Ctrl+Shift+C</b> · Paste: <b>Ctrl+Shift+V</b> · Right-click: paste
        </div>

        {/* ✅ Terminal area: tıklayınca focus */}
        <div style={{ flex: 1, padding: 8 }} onMouseDown={() => xtermRef.current?.focus()}>
          <div ref={termElRef} style={{ width: "100%", height: "100%" }} />
        </div>
      </div>
    </div>
  );
}
