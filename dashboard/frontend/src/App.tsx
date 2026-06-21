import { useEffect, useState } from "react";
import { DashboardStatus } from "./types";
import ContainerCard from "./components/ContainerCard";
import EndpointCard from "./components/EndpointCard";
import SystemCard from "./components/SystemCard";

export default function App() {
  const [status, setStatus] = useState<DashboardStatus | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    const load = () =>
      fetch("/api/status")
        .then((r) => r.json())
        .then((d: DashboardStatus) => {
          setStatus(d);
          setError(false);
        })
        .catch(() => setError(true));

    load();
    const id = setInterval(load, 10_000);
    return () => clearInterval(id);
  }, []);

  if (error) return <div className="center">API 接続エラー</div>;
  if (!status) return <div className="center">読み込み中…</div>;

  const updatedAt = new Date(status.updated_at).toLocaleTimeString("ja-JP");

  return (
    <div className="dashboard">
      <header>
        <h1>luy-XA7C-R38</h1>
        <span className="updated">updated {updatedAt}</span>
      </header>

      {status.containers.length > 0 && (
        <section>
          <h2>Containers</h2>
          <div className="grid">
            {status.containers.map((c) => (
              <ContainerCard key={c.name} c={c} />
            ))}
          </div>
        </section>
      )}

      {status.endpoints && status.endpoints.length > 0 && (
        <section>
          <h2>Services</h2>
          <div className="grid">
            {status.endpoints.map((e) => (
              <EndpointCard key={e.name} e={e} />
            ))}
          </div>
        </section>
      )}

      <section>
        <h2>System</h2>
        <div className="grid" style={{ gridTemplateColumns: "minmax(280px, 480px)" }}>
          <SystemCard s={status.system} />
        </div>
      </section>
    </div>
  );
}
