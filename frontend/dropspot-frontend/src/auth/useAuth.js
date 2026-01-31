import { useEffect, useMemo, useState } from "react";
import * as authApi from "../api/auth";

export function useAuthProvider() {
  const [token, setToken] = useState(localStorage.getItem("accessToken") || "");
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  function persistToken(t) {
    if (t) localStorage.setItem("accessToken", t);
    else localStorage.removeItem("accessToken");
    setToken(t || "");
  }

  async function refreshMe(t = token) {
    if (!t) {
      setUser(null);
      return;
    }

    try {
      // auth:true => apiFetch Authorization header eklemeli
      const res = await authApi.me();
      setUser(res.data);
    } catch {
      // token invalid -> temizle
      persistToken("");
      setUser(null);
    }
  }

  // ✅ Email/Password Login
  async function login(email, password) {
    const res = await authApi.login(email, password);
    const accessToken = res?.data?.accessToken;
    if (!accessToken) throw new Error("accessToken not returned from login");

    persistToken(accessToken);
    await refreshMe(accessToken);
  }

  // ✅ Google Login (idToken -> backend -> accessToken)
  async function loginWithGoogleToken(idToken) {
    const res = await authApi.loginWithGoogle(idToken);
    const accessToken = res?.data?.accessToken;
    if (!accessToken) throw new Error("accessToken not returned from google login");

    persistToken(accessToken);
    await refreshMe(accessToken);
  }

  function logout() {
    persistToken("");
    setUser(null);
  }

  useEffect(() => {
    (async () => {
      setLoading(true);
      await refreshMe();
      setLoading(false);
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return useMemo(
    () => ({
      token,
      user,
      loading,
      isAuthenticated: !!token,
      login,
      loginWithGoogleToken,
      logout,
      refreshMe,
    }),
    [token, user, loading]
  );
}
