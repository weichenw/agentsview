import React, { useState, useEffect, useCallback } from "react";
import * as api from "../../api/scheduler.js";
import type { Job, SchedulerRun } from "../../types/scheduler.js";

/* ───── helper: human-readable cron ───── */
function cronHumanReadable(expr: string): string {
  // Basic 5/6-field cron -> English approximation
  const parts = expr.trim().split(/\s+/);
  // If 6 fields (with seconds), shift off the seconds field
  if (parts.length === 6) parts.shift();
  if (parts.length !== 5) return expr;

  const [min, hour, dom, mon, dow] = parts;
  if (min === "0" && hour === "7" && dom === "*" && mon === "*" && dow === "*")
    return "Every day at 7:00 AM";
  if (min === "0" && hour === "0" && dom === "*" && mon === "*" && dow === "*")
    return "Every day at midnight";
  if (min === "0" && hour === "14" && dom === "*" && mon === "*" && dow === "1")
    return "Every Monday at 2:00 PM";
  if (min === "*/5" && hour === "*" && dom === "*" && mon === "*" && dow === "*")
    return "Every 5 minutes";
  return expr;
}

function formatRelativeTime(ts: string): string {
  if (!ts) return "";
  const diff = Date.now() - new Date(ts).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m ago`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  return new Date(ts).toLocaleDateString();
}

function formatDuration(start: string, finish?: string): string {
  if (!start || !finish) return "";
  const ms = new Date(finish).getTime() - new Date(start).getTime();
  const secs = Math.floor(ms / 1000);
  if (secs < 60) return `${secs}s`;
  return `${Math.floor(secs / 60)}m ${secs % 60}s`;
}

/* ───── views ───── */

type View = "list" | "detail" | "create" | "runs";

export default function SchedulerPage() {
  const [view, setView] = useState<View>("list");
  const [jobs, setJobs] = useState<Job[]>([]);
  const [selectedJob, setSelectedJob] = useState<Job | null>(null);
  const [runs, setRuns] = useState<SchedulerRun[]>([]);
  const [loading, setLoading] = useState(true);

  // Job form state
  const [form, setForm] = useState({
    name: "",
    cron: "",
    agent: "",
    prompt: "",
    model: "",
    working_dir: "",
    spawn_mode: "cmux",
    inherit_project_context: false,
    enabled: true,
  });

  const loadJobs = useCallback(async () => {
    try {
      const j = await api.listJobs();
      setJobs(j);
    } catch (e) {
      console.error("Failed to load jobs", e);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadJobs();
  }, [loadJobs]);

  async function handleRun(job: Job) {
    try {
      await api.runJob(job.id);
      // Reload runs if viewing runs for this job
      if (view === "runs" && selectedJob?.id === job.id) {
        const r = await api.listRuns(job.id, 10);
        setRuns(r);
      }
    } catch (e) {
      console.error("Failed to run job", e);
    }
  }

  async function handleToggleEnabled(job: Job) {
    try {
      if (job.enabled) {
        await api.disableJob(job.id);
      } else {
        await api.enableJob(job.id);
      }
      await loadJobs();
    } catch (e) {
      console.error("Failed to toggle job", e);
    }
  }

  async function handleDelete(job: Job) {
    if (!confirm(`Delete job "${job.name}"?`)) return;
    try {
      await api.deleteJob(job.id);
      await loadJobs();
      if (selectedJob?.id === job.id) {
        setSelectedJob(null);
        setView("list");
      }
    } catch (e) {
      console.error("Failed to delete job", e);
    }
  }

  async function handleCreate() {
    try {
      await api.createJob({
        name: form.name,
        cron: form.cron,
        agent: form.agent,
        prompt: form.prompt,
        model: form.model || undefined,
        working_dir: form.working_dir,
        spawn_mode: form.spawn_mode,
        inherit_project_context: form.inherit_project_context,
        enabled: form.enabled,
      });
      setForm({ name: "", cron: "", agent: "", prompt: "", model: "", working_dir: "", spawn_mode: "cmux", inherit_project_context: false, enabled: true });
      await loadJobs();
      setView("list");
    } catch (e) {
      console.error("Failed to create job", e);
    }
  }

  async function handleUpdate() {
    if (!selectedJob) return;
    try {
      await api.updateJob(selectedJob.id, {
        name: form.name,
        cron: form.cron,
        agent: form.agent,
        prompt: form.prompt,
        model: form.model || undefined,
        working_dir: form.working_dir,
        spawn_mode: form.spawn_mode,
        inherit_project_context: form.inherit_project_context,
        enabled: form.enabled,
      });
      await loadJobs();
      setView("list");
    } catch (e) {
      console.error("Failed to update job", e);
    }
  }

  async function viewRuns(job: Job) {
    setSelectedJob(job);
    try {
      const r = await api.listRuns(job.id, 20);
      setRuns(r);
    } catch {
      setRuns([]);
    }
    setView("runs");
  }

  function editJob(job: Job) {
    setSelectedJob(job);
    setForm({
      name: job.name,
      cron: job.cron,
      agent: job.agent,
      prompt: job.prompt,
      model: job.model || "",
      working_dir: job.working_dir,
      spawn_mode: job.spawn_mode,
      inherit_project_context: job.inherit_project_context,
      enabled: job.enabled,
    });
    setView("detail");
  }

  /* ───── render: job list ───── */
  if (view === "list") {
    return (
      <div style={{ padding: "20px", maxWidth: "900px", margin: "0 auto" }}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "16px" }}>
          <h2 style={{ margin: 0, fontSize: "18px", fontWeight: 600 }}>Scheduler</h2>
          <button onClick={() => { setForm({ name: "", cron: "", agent: "", prompt: "", model: "", working_dir: "", spawn_mode: "cmux", inherit_project_context: false, enabled: true }); setView("create"); }}
            style={{ padding: "6px 14px", fontSize: "13px", fontWeight: 500, background: "var(--accent-blue, #3b82f6)", color: "#fff", border: "none", borderRadius: "6px", cursor: "pointer" }}>
            + New Job
          </button>
        </div>

        {loading ? (
          <p style={{ color: "var(--text-muted, #888)", fontSize: "13px" }}>Loading...</p>
        ) : jobs.length === 0 ? (
          <p style={{ color: "var(--text-muted, #888)", fontSize: "13px" }}>No scheduled jobs. Create one to get started.</p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
            {jobs.map((job) => (
              <div key={job.id}
                style={{ display: "flex", alignItems: "center", gap: "12px", padding: "12px 14px", background: "var(--bg-surface, #fff)", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "8px" }}>
                <div style={{ flex: 1, minWidth: 0, cursor: "pointer" }} onClick={() => editJob(job)}>
                  <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                    <span style={{ width: "8px", height: "8px", borderRadius: "50%", background: job.enabled ? "var(--accent-green, #22c55e)" : "var(--text-muted, #888)", flexShrink: 0 }} />
                    <span style={{ fontSize: "14px", fontWeight: 600, color: "var(--text-primary, #111)" }}>{job.name}</span>
                    <span style={{ fontSize: "11px", color: "var(--text-muted, #888)", fontFamily: "monospace" }}>{job.cron}</span>
                  </div>
                  <div style={{ fontSize: "11px", color: "var(--text-muted, #888)", marginTop: "2px" }}>
                    {cronHumanReadable(job.cron)}
                  </div>
                </div>

                <button onClick={() => handleRun(job)}
                  style={{ padding: "4px 10px", fontSize: "11px", background: "var(--bg-surface-hover, #f3f4f6)", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "4px", cursor: "pointer", color: "var(--text-secondary, #555)" }}
                  title="Run now">
                  ▶ Run
                </button>

                <button onClick={() => handleToggleEnabled(job)}
                  style={{ padding: "4px 10px", fontSize: "11px", background: "transparent", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "4px", cursor: "pointer", color: "var(--text-muted, #888)" }}>
                  {job.enabled ? "Disable" : "Enable"}
                </button>

                <button onClick={() => viewRuns(job)}
                  style={{ padding: "4px 10px", fontSize: "11px", background: "transparent", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "4px", cursor: "pointer", color: "var(--text-muted, #888)" }}>
                  History
                </button>

                <button onClick={() => handleDelete(job)}
                  style={{ padding: "4px 10px", fontSize: "11px", background: "transparent", border: "1px solid transparent", borderRadius: "4px", cursor: "pointer", color: "var(--accent-red, #ef4444)" }}>
                  ✕
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }

  /* ───── render: create/edit form ───── */
  if (view === "create" || view === "detail") {
    const isEdit = view === "detail";
    return (
      <div style={{ padding: "20px", maxWidth: "600px", margin: "0 auto" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "12px", marginBottom: "20px" }}>
          <button onClick={() => setView("list")}
            style={{ background: "none", border: "none", cursor: "pointer", fontSize: "16px", color: "var(--text-muted, #888)", padding: "0" }}>
            ← Back
          </button>
          <h2 style={{ margin: 0, fontSize: "18px", fontWeight: 600 }}>{isEdit ? "Edit Job" : "New Job"}</h2>
        </div>

        <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Name
            <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", boxSizing: "border-box" }} />
          </label>

          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Cron Expression
            <input value={form.cron} onChange={(e) => setForm({ ...form, cron: e.target.value })}
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", fontFamily: "monospace", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", boxSizing: "border-box" }} />
            {form.cron && (
              <span style={{ fontSize: "11px", color: "var(--text-muted, #888)", marginTop: "4px", display: "block" }}>
                {cronHumanReadable(form.cron)}
              </span>
            )}
          </label>

          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Agent
            <input value={form.agent} onChange={(e) => setForm({ ...form, agent: e.target.value })}
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", boxSizing: "border-box" }} />
          </label>

          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Prompt
            <textarea value={form.prompt} onChange={(e) => setForm({ ...form, prompt: e.target.value })} rows={4}
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", fontFamily: "monospace", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", resize: "vertical", boxSizing: "border-box" }} />
          </label>

          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Model (optional)
            <input value={form.model} onChange={(e) => setForm({ ...form, model: e.target.value })} placeholder="e.g. anthropic/claude-sonnet-4"
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", boxSizing: "border-box" }} />
          </label>

          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Working Directory
            <input value={form.working_dir} onChange={(e) => setForm({ ...form, working_dir: e.target.value })} placeholder="e.g. /Users/me/projects/myapp"
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", boxSizing: "border-box" }} />
          </label>

          <label style={{ fontSize: "12px", fontWeight: 500, color: "var(--text-secondary, #555)" }}>
            Spawn Mode
            <select value={form.spawn_mode} onChange={(e) => setForm({ ...form, spawn_mode: e.target.value })}
              style={{ display: "block", width: "100%", marginTop: "4px", padding: "8px 10px", fontSize: "13px", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", background: "var(--bg-inset, #f9fafb)", color: "var(--text-primary, #111)", boxSizing: "border-box" }}>
              <option value="cmux">cmux (visible terminal)</option>
              <option value="subprocess">subprocess (headless)</option>
            </select>
          </label>

          <label style={{ display: "flex", alignItems: "center", gap: "8px", fontSize: "13px", color: "var(--text-secondary, #555)" }}>
            <input type="checkbox" checked={form.inherit_project_context} onChange={(e) => setForm({ ...form, inherit_project_context: e.target.checked })} />
            Inherit project context
          </label>

          <label style={{ display: "flex", alignItems: "center", gap: "8px", fontSize: "13px", color: "var(--text-secondary, #555)" }}>
            <input type="checkbox" checked={form.enabled} onChange={(e) => setForm({ ...form, enabled: e.target.checked })} />
            Enabled
          </label>

          <div style={{ display: "flex", gap: "8px", marginTop: "8px" }}>
            <button onClick={isEdit ? handleUpdate : handleCreate}
              style={{ padding: "8px 20px", fontSize: "13px", fontWeight: 500, background: "var(--accent-blue, #3b82f6)", color: "#fff", border: "none", borderRadius: "6px", cursor: "pointer" }}>
              {isEdit ? "Save Changes" : "Create Job"}
            </button>
            <button onClick={() => setView("list")}
              style={{ padding: "8px 20px", fontSize: "13px", background: "transparent", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", cursor: "pointer", color: "var(--text-secondary, #555)" }}>
              Cancel
            </button>
          </div>
        </div>
      </div>
    );
  }

  /* ───── render: run history ───── */
  if (view === "runs" && selectedJob) {
    return (
      <div style={{ padding: "20px", maxWidth: "600px", margin: "0 auto" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "12px", marginBottom: "20px" }}>
          <button onClick={() => setView("list")}
            style={{ background: "none", border: "none", cursor: "pointer", fontSize: "16px", color: "var(--text-muted, #888)", padding: "0" }}>
            ← Back
          </button>
          <h2 style={{ margin: 0, fontSize: "18px", fontWeight: 600 }}>Run History — {selectedJob.name}</h2>
        </div>

        {runs.length === 0 ? (
          <p style={{ color: "var(--text-muted, #888)", fontSize: "13px" }}>No run history yet.</p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
            {runs.map((run) => (
              <div key={run.id}
                style={{ display: "flex", alignItems: "center", gap: "12px", padding: "10px 14px", background: "var(--bg-surface, #fff)", border: "1px solid var(--border-default, #e5e7eb)", borderRadius: "6px", fontSize: "12px" }}>
                <span style={{
                  width: "8px", height: "8px", borderRadius: "50%", flexShrink: 0,
                  background: run.status === "completed" ? "var(--accent-green, #22c55e)"
                    : run.status === "failed" ? "var(--accent-red, #ef4444)"
                    : run.status === "killed" ? "var(--accent-amber, #f59e0b)"
                    : "var(--accent-blue, #3b82f6)"
                }} />
                <span style={{ color: "var(--text-secondary, #555)", whiteSpace: "nowrap" }}>
                  {new Date(run.started_at).toLocaleString()}
                </span>
                <span style={{
                  fontWeight: 500,
                  color: run.status === "completed" ? "var(--accent-green, #22c55e)"
                    : run.status === "failed" ? "var(--accent-red, #ef4444)"
                    : run.status === "killed" ? "var(--accent-amber, #f59e0b)"
                    : "var(--accent-blue, #3b82f6)"
                }}>
                  {run.status}
                </span>
                <span style={{ color: "var(--text-muted, #888)" }}>
                  {formatDuration(run.started_at, run.finished_at)}
                </span>
                <span style={{ color: "var(--text-muted, #888)", flex: 1, minWidth: 0 }}>
                  {formatRelativeTime(run.started_at)}
                </span>
                {run.session_id && (
                  <a href={`/sessions/${encodeURIComponent(run.session_id)}`}
                    style={{ fontSize: "11px", color: "var(--accent-blue, #3b82f6)", textDecoration: "none", whiteSpace: "nowrap" }}>
                    → view session
                  </a>
                )}
                {run.error && (
                  <span style={{ color: "var(--accent-red, #ef4444)", fontSize: "11px", maxWidth: "200px", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }} title={run.error}>
                    {run.error}
                  </span>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }

  return null;
}
