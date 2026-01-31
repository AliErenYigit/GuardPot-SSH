import { Navigate } from "react-router-dom";
import { useContext } from "react";
import { AuthContext } from "./AuthContext";

export default function ProtectedRoute({ children }) {
  const { token, loading } = useContext(AuthContext);

  if (loading) return <div style={{ padding: 24 }}>Loading...</div>;
  if (!token) return <Navigate to="/login" replace />;
  return children;
}
