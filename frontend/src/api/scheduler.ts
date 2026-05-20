import type { Job, JobFormData, SchedulerRun, RunResponse } from "../types/scheduler.js";

function getBase(): string {
  const server = localStorage.getItem("agentsview-server-url");
  if (server) return `${server}/api/v1`;
  const baseEl = document.querySelector("base[href]");
  if (baseEl) {
    const base = new URL(document.baseURI).pathname.replace(/\/$/, "");
    return `${base}/api/v1`;
  }
  return "/api/v1";
}

function authHeaders(init?: RequestInit): RequestInit {
  const token = localStorage.getItem("agentsview-auth-token");
  if (!token) return init ?? {};
  const headers = new Headers(init?.headers);
  headers.set("Authorization", `Bearer ${token}`);
  return { ...init, headers };
}

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${getBase()}${path}`, authHeaders(init));
  if (!res.ok) {
    const body = await res.text().catch(() => "");
    throw new Error(body.trim() || `API ${res.status}`);
  }
  // Handle 204 No Content
  if (res.status === 204) return undefined as unknown as T;
  return res.json() as Promise<T>;
}

function buildQuery(params: Record<string, string | number | boolean | undefined | null>): string {
  const q = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "") {
      q.set(key, String(value));
    }
  }
  const qs = q.toString();
  return qs ? `?${qs}` : "";
}

export async function listJobs(): Promise<Job[]> {
  const res = await fetchJSON<{ jobs: Job[] }>("/scheduler/jobs");
  return res.jobs;
}

export async function getJob(id: string): Promise<Job> {
  return fetchJSON<Job>(`/scheduler/jobs/${id}`);
}

export async function createJob(data: JobFormData): Promise<Job> {
  return fetchJSON<Job>("/scheduler/jobs", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
}

export async function updateJob(id: string, data: Partial<JobFormData>): Promise<Job> {
  return fetchJSON<Job>(`/scheduler/jobs/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
}

export async function deleteJob(id: string): Promise<void> {
  return fetchJSON<void>(`/scheduler/jobs/${id}`, { method: "DELETE" });
}

export async function runJob(id: string): Promise<RunResponse> {
  return fetchJSON<RunResponse>(`/scheduler/jobs/${id}/run`, { method: "POST" });
}

export async function enableJob(id: string): Promise<Job> {
  return fetchJSON<Job>(`/scheduler/jobs/${id}/enable`, { method: "POST" });
}

export async function disableJob(id: string): Promise<Job> {
  return fetchJSON<Job>(`/scheduler/jobs/${id}/disable`, { method: "POST" });
}

export async function listRuns(jobId: string, limit?: number): Promise<SchedulerRun[]> {
  const res = await fetchJSON<{ runs: SchedulerRun[] }>(
    `/scheduler/runs${buildQuery({ job_id: jobId, limit })}`,
  );
  return res.runs;
}
