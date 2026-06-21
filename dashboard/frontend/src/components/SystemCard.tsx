import { SystemMetrics } from "../types";

function Bar({
  label,
  value,
  max,
  unit,
}: {
  label: string;
  value: number;
  max: number;
  unit: string;
}) {
  const pct = max > 0 ? (value / max) * 100 : 0;
  const cls = pct > 90 ? "crit" : pct > 70 ? "warn" : "";
  return (
    <div className="bar-wrap">
      <div className="bar-label">
        <span>{label}</span>
        <span>
          {value.toFixed(1)} / {max.toFixed(1)} {unit}
        </span>
      </div>
      <div className="bar-track">
        <div className={`bar-fill ${cls}`} style={{ width: `${pct}%` }} />
      </div>
    </div>
  );
}

export default function SystemCard({ s }: { s: SystemMetrics }) {
  const cpuCls =
    s.cpu_percent > 90 ? "crit" : s.cpu_percent > 70 ? "warn" : "";
  return (
    <div className="card">
      <div className="bar-wrap">
        <div className="bar-label">
          <span>CPU</span>
          <span>{s.cpu_percent.toFixed(1)}%</span>
        </div>
        <div className="bar-track">
          <div
            className={`bar-fill ${cpuCls}`}
            style={{ width: `${s.cpu_percent}%` }}
          />
        </div>
      </div>
      <Bar label="RAM" value={s.ram_used_gb} max={s.ram_total_gb} unit="GB" />
    </div>
  );
}
