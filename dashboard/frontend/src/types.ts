export interface ContainerStatus {
  name: string;
  state: string;
  status: string;
}

export interface EndpointStatus {
  name: string;
  url: string;
  ok: boolean;
  latency_ms: number;
}

export interface SystemMetrics {
  cpu_percent: number;
  ram_used_gb: number;
  ram_total_gb: number;
}

export interface DashboardStatus {
  containers: ContainerStatus[];
  endpoints: EndpointStatus[];
  system: SystemMetrics;
  updated_at: string;
}
