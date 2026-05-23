/** Matches the Go scheduler.Job struct in internal/scheduler/store.go */
export interface Job {
  id: string;
  name: string;
  cron: string;
  enabled: boolean;
  agent: string;
  prompt: string;
  model?: string;
  inherit_project_context: boolean;
  working_dir: string;
  spawn_mode: string;
  created_at: string;
  updated_at: string;
}

/** Input for creating or updating a scheduler job. */
export interface JobFormData {
  id?: string;
  name: string;
  cron: string;
  enabled?: boolean;
  agent: string;
  prompt: string;
  model?: string;
  inherit_project_context?: boolean;
  working_dir: string;
  spawn_mode?: string;
}

/** Matches the Go scheduler.SchedulerRun struct in scheduler_runs table. */
export interface SchedulerRun {
  id: string;
  job_id: string;
  session_id?: string;
  started_at: string;
  finished_at?: string;
  status: string; // running | completed | failed | killed
  exit_code?: number | null;
  error?: string;
}

export interface SchedulerSettings {
  timezone: string;
}

/** Response from POST /api/v1/scheduler/jobs/{id}/run */
export interface RunResponse {
  run_id: string;
  status: string;
  job_id: string;
}
