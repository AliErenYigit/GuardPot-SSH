import { useState } from "react";
import { register } from "../api/auth";
import { Link, useNavigate } from "react-router-dom";
import "./login.css"; // Login ile aynı CSS

export default function Register() {
  const nav = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const [busy, setBusy] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();
    setErr("");
    setBusy(true);
    try {
      await register(email, password);
      nav("/login");
    } catch (e) {
      setErr(e?.message || "Registration failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="authPage">
      <div className="authBg" aria-hidden="true" />

      <div className="authShell">
        <div className="authCard">
          <div className="authHeader">
            <div className="authBadge">GuardPot</div>
            <h1 className="authTitle">Create your account</h1>
            <p className="authSubtitle">
              Get started with secure SSH connections and live terminal access.
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
              />
            </label>

            <label className="authLabel">
              <span>Password</span>
              <input
                className="authInput"
                placeholder="Minimum 8 characters"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="new-password"
                minLength={8}
                required
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
                  <span className="authSpinner" />
                  Creating account...
                </>
              ) : (
                "Create account"
              )}
            </button>

            <div className="authFooter">
              <span>Already have an account?</span>
              <Link className="authLink" to="/login">
                Login
              </Link>
            </div>
          </form>
        </div>

        <div className="authHint">
          <div className="authHintDot" />
          <span>Password must be at least 8 characters long.</span>
        </div>
      </div>
    </div>
  );
}
