import { useContext, useState } from "react";
import { AuthContext } from "../auth/AuthContext";
import SSHConnectionForm from "../components/SSHConnectionForm";
import SSHConnectionList from "../components/SSHConnectionList";

export default function Dashboard() {
  const { user, logout } = useContext(AuthContext);
  const [refreshKey, setRefreshKey] = useState(0);

  return (
    <div className="container">
      <div className="topbar">
        <div>
          <h1 className="h1">Dashboard</h1>
          <div className="sub">
            Logged in as: <b>{user?.email}</b>
          </div>
        </div>

        <button className="btn btnDanger" onClick={logout}>
          Logout
        </button>
      </div>

      <div className="grid">
        <SSHConnectionForm onCreated={() => setRefreshKey((k) => k + 1)} />
        <SSHConnectionList refreshKey={refreshKey} />
      </div>
    </div>
  );
}
