import { useContext, useEffect, useRef, useState } from "react";
import { AuthContext } from "../auth/AuthContext";
import { Link, useNavigate } from "react-router-dom";
import "./login.css";

export default function Login() {
  const nav = useNavigate();
  const { login, loginWithGoogleToken } = useContext(AuthContext);

  const googleBtnRef = useRef(null);
  const [googleReady, setGoogleReady] = useState(false);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const [busy, setBusy] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();
    setErr("");
    setBusy(true);
    try {
      await login(email, password);
      nav("/");
    } catch (e) {
      setErr(e?.message || "Login failed");
    } finally {
      setBusy(false);
    }
  }

  // ✅ Google Button Init
  useEffect(() => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
    console.log(clientId);
    if (!clientId) return;

    // Google script yüklenmemişse buton gösterme
    const google = window.google?.accounts?.id;
    if (!google || !googleBtnRef.current) return;

    try {
      google.initialize({
        client_id: clientId,
        callback: async (resp) => {
          // resp.credential = id_token
          setErr("");
          setBusy(true);
          try {
            await loginWithGoogleToken(resp.credential);
            nav("/");
          } catch (e) {
            setErr(e?.message || "Google login failed");
          } finally {
            setBusy(false);
          }
        },
      });

      // Butonu render et
      google.renderButton(googleBtnRef.current, {
        theme: "outline",
        size: "large",
        shape: "pill",
        width: 380,
        text: "continue_with",
      });

      setGoogleReady(true);
    } catch (err) {
      // sessizce geç, normal login çalışsın
      console.log(err);
      setGoogleReady(false);
    }
  }, [loginWithGoogleToken, nav]);

  return (
    <div className="authPage">
      <div className="authBg" aria-hidden="true" />

      <div className="authShell">
        <div className="authCard">
          <div className="authHeader">
            <div className="authBadge">GuardPot</div>
            <h1 className="authTitle">Welcome back</h1>
            <p className="authSubtitle">
              Sign in to manage your SSH connections and live terminal.
            </p>
          </div>

          <form className="authForm" onSubmit={onSubmit}>
            <label className="authLabel">
              <span>Email</span>
              <input
                className="authInput"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
                inputMode="email"
                required
                disabled={busy}
              />
            </label>

            <label className="authLabel">
              <span>Password</span>
              <input
                className="authInput"
                placeholder="••••••••"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="current-password"
                required
                disabled={busy}
              />
            </label>

            {err && (
              <div className="authError" role="alert">
                <strong>Oops:</strong> {err}
              </div>
            )}

            <button className="authButton" disabled={busy} type="submit">
              {busy ? (
                <>
                  <span className="authSpinner" aria-hidden="true" />
                  Signing in...
                </>
              ) : (
                "Login"
              )}
            </button>

            {/* ✅ Divider */}
            <div className="authDivider" aria-hidden="true">
              <span>or</span>
            </div>

            {/* ✅ Google button placeholder */}
            <div
              className={`googleWrap ${busy ? "googleDisabled" : ""}`}
              aria-disabled={busy}
            >
              <div ref={googleBtnRef} />
              {!googleReady && (
                <div className="googleFallback">
                  Google sign-in is not available (missing script or Client ID).
                </div>
              )}
            </div>

            <div className="authFooter">
              <span>No account?</span>
              <Link className="authLink" to="/register">
                Register
              </Link>
            </div>
          </form>
        </div>

        <div className="authHint">
          <div className="authHintDot" />
          <span>Tip: Use a strong password and keep your token secure.</span>
        </div>
      </div>
    </div>
  );
}
