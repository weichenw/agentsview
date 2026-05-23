import { getBase, getAuthToken, ApiError } from "../lib/api/client.js";
import type { Job, JobFormData, SchedulerRun, RunResponse, SchedulerSettings } from "../types/scheduler.js";

function authHeaders(init?: RequestInit): RequestInit {
  const token = getAuthToken();
  if (!token) return init ?? {};
  const headers = new Headers(init?.headers);
  headers.set("Authorization", `Bearer ${token}`);
  return { ...init, headers };
}

async function responseErrorMessage(res: Response): Promise<string> {
  const body = await res.text().catch(() => "");
  const text = body.trim();
  if (!text) return `API ${res.status}`;
  try {
    const parsed = JSON.parse(text) as unknown;
    if (parsed !== null && typeof parsed === "object" && "error" in parsed && typeof parsed.error === "string" && parsed.error) {
      return parsed.error;
    }
  } catch { /* plain text */ }
  return text;
}

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${getBase()}${path}`, authHeaders(init));
  if (!res.ok) {
    throw new ApiError(res.status, await responseErrorMessage(res));
  }
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

export async function killRun(runId: string): Promise<void> {
  return fetchJSON<void>(`/scheduler/runs/${runId}/kill`, { method: "POST" });
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

export async function getSettings(): Promise<SchedulerSettings> {
  return fetchJSON<SchedulerSettings>("/scheduler/settings");
}

export async function updateSettings(data: SchedulerSettings): Promise<SchedulerSettings> {
  return fetchJSON<SchedulerSettings>("/scheduler/settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
}
