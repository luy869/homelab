import { EndpointStatus } from "../types";

export default function EndpointCard({ e }: { e: EndpointStatus }) {
  return (
    <div className="card">
      <div className="card-row">
        <span className={`dot ${e.ok ? "ok" : "error"}`} />
        <span className="card-name">{e.name}</span>
        {e.ok && <span className="latency">{e.latency_ms}ms</span>}
      </div>
    </div>
  );
}
