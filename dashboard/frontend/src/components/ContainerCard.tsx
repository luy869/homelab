import { ContainerStatus } from "../types";

export default function ContainerCard({ c }: { c: ContainerStatus }) {
  return (
    <div className="card">
      <div className="card-row">
        <span className={`dot ${c.state === "running" ? "ok" : "error"}`} />
        <span className="card-name">{c.name}</span>
      </div>
      <div className="card-sub">{c.status}</div>
    </div>
  );
}
