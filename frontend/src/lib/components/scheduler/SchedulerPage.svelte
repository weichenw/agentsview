<script lang="ts">
  import { onMount } from "svelte";
  import * as api from "../../api/scheduler.js";
  import type { Job, SchedulerRun } from "../../types/scheduler.js";
  import { router } from "../../stores/router.svelte.js";

  type View = "list" | "detail" | "create" | "runs";

  let view: View = $state("list");
  let jobs: Job[] = $state([]);
  let selectedJob: Job | null = $state(null);
  let runs: SchedulerRun[] = $state([]);
  let loading = $state(true);

  // Form state
  let formName = $state("");
  let formCron = $state("");
  let formAgent = $state("");
  let formPrompt = $state("");
  let formModel = $state("");
  let formWorkingDir = $state("");
  let formSpawnMode = $state("cmux");
  let formInheritContext = $state(false);
  let formEnabled = $state(true);

  function cronHumanReadable(expr: string): string {
    const parts = expr.trim().split(/\s+/);
    if (parts.length === 6) parts.shift();
    if (parts.length !== 5) return expr;
    const [min, hour, dom, mon, dow] = parts;
    if (min === "0" && hour === "7" && dom === "*" && mon === "*" && dow === "*") return "Every day at 7:00 AM";
    if (min === "0" && hour === "0" && dom === "*" && mon === "*" && dow === "*") return "Every day at midnight";
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

  function resetForm() {
    formName = ""; formCron = ""; formAgent = ""; formPrompt = "";
    formModel = ""; formWorkingDir = ""; formSpawnMode = "cmux";
    formInheritContext = false; formEnabled = true;
  }

  function fillForm(job: Job) {
    formName = job.name; formCron = job.cron; formAgent = job.agent;
    formPrompt = job.prompt; formModel = job.model ?? "";
    formWorkingDir = job.working_dir; formSpawnMode = job.spawn_mode;
    formInheritContext = job.inherit_project_context; formEnabled = job.enabled;
  }

  onMount(() => { loadJobs(); });

  async function loadJobs() {
    loading = true;
    try { jobs = await api.listJobs(); } catch { jobs = []; }
    finally { loading = false; }
  }

  async function handleRun(job: Job) {
    try {
      await api.runJob(job.id);
      if (view === "runs" && selectedJob?.id === job.id) {
        runs = await api.listRuns(job.id, 10);
      }
    } catch { /* ignore */ }
  }

  async function handleToggle(job: Job) {
    try {
      if (job.enabled) { await api.disableJob(job.id); }
      else { await api.enableJob(job.id); }
      await loadJobs();
    } catch { /* ignore */ }
  }

  async function handleDelete(job: Job) {
    if (!confirm(`Delete job "${job.name}"?`)) return;
    try {
      await api.deleteJob(job.id);
      await loadJobs();
      if (selectedJob?.id === job.id) { selectedJob = null; view = "list"; }
    } catch { /* ignore */ }
  }

  async function handleCreate() {
    try {
      await api.createJob({ name: formName, cron: formCron, agent: formAgent, prompt: formPrompt, model: formModel || undefined, working_dir: formWorkingDir, spawn_mode: formSpawnMode, inherit_project_context: formInheritContext, enabled: formEnabled });
      resetForm(); await loadJobs(); view = "list";
    } catch { /* ignore */ }
  }

  async function handleUpdate() {
    if (!selectedJob) return;
    try {
      await api.updateJob(selectedJob.id, { name: formName, cron: formCron, agent: formAgent, prompt: formPrompt, model: formModel || undefined, working_dir: formWorkingDir, spawn_mode: formSpawnMode, inherit_project_context: formInheritContext, enabled: formEnabled });
      await loadJobs(); view = "list";
    } catch { /* ignore */ }
  }

  async function viewRuns(job: Job) {
    selectedJob = job;
    try { runs = await api.listRuns(job.id, 20); }
    catch { runs = []; }
    view = "runs";
  }

  function editJob(job: Job) {
    selectedJob = job;
    fillForm(job);
    view = "detail";
  }

  function navToSession(id: string) {
    router.navigateToSession(id);
  }
</script>

<div class="scheduler-page">
  {#if view === "list"}
    <div class="scheduler-header">
      <h2>Scheduler</h2>
      <button class="btn-primary" onclick={() => { resetForm(); view = "create"; }}>
        + New Job
      </button>
    </div>

    {#if loading}
      <p class="empty-state">Loading...</p>
    {:else if jobs.length === 0}
      <p class="empty-state">No scheduled jobs. Create one to get started.</p>
    {:else}
      <div class="job-list">
        {#each jobs as job (job.id)}
          <div class="job-card">
            <div class="job-main" onclick={() => editJob(job)} role="button" tabindex="0" onkeydown={(e) => { if (e.key === 'Enter') editJob(job); }}>
              <span class="status-dot" class:enabled={job.enabled} class:disabled={!job.enabled}></span>
              <span class="job-name">{job.name}</span>
              <span class="job-cron">{job.cron}</span>
              <span class="job-cron-preview">{cronHumanReadable(job.cron)}</span>
            </div>
            <div class="job-actions">
              <button class="btn-sm" onclick={() => handleRun(job)} title="Run now">▶ Run</button>
              <button class="btn-sm" onclick={() => handleToggle(job)}>{job.enabled ? "Disable" : "Enable"}</button>
              <button class="btn-sm" onclick={() => viewRuns(job)}>History</button>
              <button class="btn-sm btn-danger" onclick={() => handleDelete(job)}>✕</button>
            </div>
          </div>
        {/each}
      </div>
    {/if}

  {:else if view === "create" || view === "detail"}
    {@const isEdit = view === "detail"}
    <div class="form-page">
      <div class="form-header">
        <button class="btn-back" onclick={() => view = "list"}>← Back</button>
        <h2>{isEdit ? "Edit Job" : "New Job"}</h2>
      </div>

      <div class="form-fields">
        <label>Name <input bind:value={formName} /></label>
        <label>
          Cron Expression
          <input bind:value={formCron} class="mono" />
          {#if formCron}
            <span class="field-hint">{cronHumanReadable(formCron)}</span>
          {/if}
        </label>
        <label>Agent <input bind:value={formAgent} /></label>
        <label>Prompt <textarea bind:value={formPrompt} rows={4} class="mono"></textarea></label>
        <label>Model (optional) <input bind:value={formModel} placeholder="e.g. anthropic/claude-sonnet-4" /></label>
        <label>Working Directory <input bind:value={formWorkingDir} placeholder="/Users/me/projects/myapp" /></label>
        <label>
          Spawn Mode
          <select bind:value={formSpawnMode}>
            <option value="cmux">cmux (visible terminal)</option>
            <option value="subprocess">subprocess (headless)</option>
          </select>
        </label>
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={formInheritContext} />
          Inherit project context
        </label>
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={formEnabled} />
          Enabled
        </label>

        <div class="form-actions">
          <button class="btn-primary" onclick={isEdit ? handleUpdate : handleCreate}>
            {isEdit ? "Save Changes" : "Create Job"}
          </button>
          <button class="btn-secondary" onclick={() => view = "list"}>Cancel</button>
        </div>
      </div>
    </div>

  {:else if view === "runs" && selectedJob}
    <div class="form-page">
      <div class="form-header">
        <button class="btn-back" onclick={() => view = "list"}>← Back</button>
        <h2>Run History — {selectedJob.name}</h2>
      </div>

      {#if runs.length === 0}
        <p class="empty-state">No run history yet.</p>
      {:else}
        <div class="runs-list">
          {#each runs as run (run.id)}
            <div class="run-row">
              <span class="status-dot" class:completed={run.status === "completed"} class:failed={run.status === "failed"} class:killed={run.status === "killed"} class:running={run.status === "running"}></span>
              <span class="run-time">{new Date(run.started_at).toLocaleString()}</span>
              <span class="run-status" class:completed={run.status === "completed"} class:failed={run.status === "failed"} class:killed={run.status === "killed"}>{run.status}</span>
              <span class="run-duration">{formatDuration(run.started_at, run.finished_at)}</span>
              <span class="run-relative">{formatRelativeTime(run.started_at)}</span>
              {#if run.session_id}
                <a href="/sessions/{encodeURIComponent(run.session_id)}" class="session-link">→ view session</a>
              {/if}
              {#if run.error}
                <span class="run-error" title={run.error}>{run.error}</span>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .scheduler-page {
    padding: 20px;
    max-width: 900px;
    margin: 0 auto;
  }

  .scheduler-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .scheduler-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary, #111);
  }

  .btn-primary {
    padding: 6px 14px;
    font-size: 13px;
    font-weight: 500;
    background: var(--accent-blue, #3b82f6);
    color: #fff;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }

  .btn-secondary {
    padding: 8px 20px;
    font-size: 13px;
    background: transparent;
    border: 1px solid var(--border-default, #e5e7eb);
    border-radius: 6px;
    cursor: pointer;
    color: var(--text-secondary, #555);
  }

  .btn-sm {
    padding: 4px 10px;
    font-size: 11px;
    background: var(--bg-surface-hover, #f3f4f6);
    border: 1px solid var(--border-default, #e5e7eb);
    border-radius: 4px;
    cursor: pointer;
    color: var(--text-secondary, #555);
  }

  .btn-danger {
    background: transparent;
    border-color: transparent;
    color: var(--accent-red, #ef4444);
  }

  .btn-back {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 16px;
    color: var(--text-muted, #888);
    padding: 0;
  }

  .empty-state {
    color: var(--text-muted, #888);
    font-size: 13px;
  }

  .job-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .job-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 14px;
    background: var(--bg-surface, #fff);
    border: 1px solid var(--border-default, #e5e7eb);
    border-radius: 8px;
  }

  .job-main {
    flex: 1;
    min-width: 0;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .job-actions {
    display: flex;
    gap: 6px;
    flex-shrink: 0;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .status-dot.enabled { background: var(--accent-green, #22c55e); }
  .status-dot.disabled { background: var(--text-muted, #888); }

  .job-name {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary, #111);
  }

  .job-cron {
    font-size: 11px;
    color: var(--text-muted, #888);
    font-family: monospace;
  }

  .job-cron-preview {
    font-size: 11px;
    color: var(--text-muted, #888);
  }

  .form-page {
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
  }

  .form-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;
  }

  .form-header h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }

  .form-fields {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .form-fields label {
    font-size: 12px;
    font-weight: 500;
    color: var(--text-secondary, #555);
    display: block;
  }

  .form-fields input, .form-fields textarea, .form-fields select {
    display: block;
    width: 100%;
    margin-top: 4px;
    padding: 8px 10px;
    font-size: 13px;
    border: 1px solid var(--border-default, #e5e7eb);
    border-radius: 6px;
    background: var(--bg-inset, #f9fafb);
    color: var(--text-primary, #111);
    box-sizing: border-box;
  }

  .form-fields .mono {
    font-family: monospace;
  }

  .field-hint {
    font-size: 11px;
    color: var(--text-muted, #888);
    margin-top: 4px;
    display: block;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: var(--text-secondary, #555);
  }

  .checkbox-label input {
    width: auto;
    margin: 0;
  }

  .form-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }

  .runs-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .run-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 14px;
    background: var(--bg-surface, #fff);
    border: 1px solid var(--border-default, #e5e7eb);
    border-radius: 6px;
    font-size: 12px;
  }

  .run-time {
    color: var(--text-secondary, #555);
    white-space: nowrap;
  }

  .run-status {
    font-weight: 500;
  }

  .run-status.completed { color: var(--accent-green, #22c55e); }
  .run-status.failed { color: var(--accent-red, #ef4444); }
  .run-status.killed { color: var(--accent-amber, #f59e0b); }

  .run-duration, .run-relative {
    color: var(--text-muted, #888);
  }

  .run-relative {
    flex: 1;
    min-width: 0;
  }

  .session-link {
    font-size: 11px;
    color: var(--accent-blue, #3b82f6);
    text-decoration: none;
    white-space: nowrap;
  }

  .run-error {
    color: var(--accent-red, #ef4444);
    font-size: 11px;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .status-dot.completed { background: var(--accent-green, #22c55e); }
  .status-dot.failed { background: var(--accent-red, #ef4444); }
  .status-dot.killed { background: var(--accent-amber, #f59e0b); }
  .status-dot.running { background: var(--accent-blue, #3b82f6); }
</style>
