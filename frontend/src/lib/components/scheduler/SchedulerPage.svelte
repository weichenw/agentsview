<script lang="ts">
  import { onMount } from "svelte";
  import * as api from "../../../api/scheduler.js";
  import type { Job, SchedulerRun } from "../../../types/scheduler.js";
  import { router } from "../../stores/router.svelte.js";

  let jobs: Job[] = $state([]);
  let runs: SchedulerRun[] = $state([]);
  let loading = $state(true);
  let selectedJob: Job | null = $state(null);
  let searchQuery = $state("");

  // Edit/create form state
  let editing = $state(false);
  let creating = $state(false);
  let showingRuns = $state(false);
  let formName = $state("");
  let formCron = $state("");
  let formAgent = $state("");
  let formPrompt = $state("");
  let formModel = $state("");
  let formWorkingDir = $state("");
  let formSpawnMode = $state("cmux");
  let formInheritContext = $state(false);
  let formEnabled = $state(true);

  const filteredJobs = $derived(
    searchQuery
      ? jobs.filter((j) =>
          j.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          j.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
          j.cron.includes(searchQuery))
      : jobs,
  );

  function cronHumanReadable(expr: string): string {
    const parts = expr.trim().split(/\s+/);
    if (parts.length === 6) parts.shift();
    if (parts.length !== 5) return expr;
    const [min, hour, dom, mon, dow] = parts;
    if (min === "0" && hour === "7" && dom === "*" && mon === "*" && dow === "*") return "Every day at 7:00 AM";
    if (min === "0" && hour === "0" && dom === "*" && mon === "*" && dow === "*") return "Every day at midnight";
    if (min === "0" && hour === "14" && dom === "*" && mon === "*" && dow === "1") return "Every Monday at 2:00 PM";
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

  onMount(() => loadJobs());

  async function loadJobs() {
    loading = true;
    try { jobs = await api.listJobs(); } catch { jobs = []; }
    finally { loading = false; }
  }

  async function handleRun(job: Job) {
    try {
      await api.runJob(job.id);
      if (showingRuns && selectedJob?.id === job.id) {
        runs = await api.listRuns(job.id, 20);
      }
    } catch { /* ignore */ }
  }

  async function handleToggle(job: Job) {
    try {
      if (job.enabled) { await api.disableJob(job.id); }
      else { await api.enableJob(job.id); }
      const updated = await api.listJobs();
      jobs = updated;
      if (selectedJob?.id === job.id) {
        selectedJob = updated.find((j) => j.id === job.id) ?? null;
      }
    } catch { /* ignore */ }
  }

  async function handleDelete(job: Job) {
    if (!confirm(`Delete job "${job.name}"?`)) return;
    try {
      await api.deleteJob(job.id);
      if (selectedJob?.id === job.id) { selectedJob = null; editing = false; showingRuns = false; }
      jobs = await api.listJobs();
    } catch { /* ignore */ }
  }

  async function handleCreate() {
    try {
      await api.createJob({
        name: formName, cron: formCron, agent: formAgent, prompt: formPrompt,
        model: formModel || undefined, working_dir: formWorkingDir,
        spawn_mode: formSpawnMode, inherit_project_context: formInheritContext, enabled: formEnabled,
      });
      resetForm(); creating = false;
      jobs = await api.listJobs();
    } catch { /* ignore */ }
  }

  async function handleUpdate() {
    if (!selectedJob) return;
    try {
      await api.updateJob(selectedJob.id, {
        name: formName, cron: formCron, agent: formAgent, prompt: formPrompt,
        model: formModel || undefined, working_dir: formWorkingDir,
        spawn_mode: formSpawnMode, inherit_project_context: formInheritContext, enabled: formEnabled,
      });
      editing = false;
      jobs = await api.listJobs();
      selectedJob = jobs.find((j) => j.id === selectedJob!.id) ?? null;
    } catch { /* ignore */ }
  }

  async function showRuns(job: Job) {
    selectedJob = job;
    showingRuns = true; editing = false; creating = false;
    try { runs = await api.listRuns(job.id, 20); }
    catch { runs = []; }
  }

  function editJob(job: Job) {
    selectedJob = job;
    fillForm(job);
    editing = true; creating = false; showingRuns = false;
  }

  function startCreate() {
    selectedJob = null;
    resetForm();
    creating = true; editing = false; showingRuns = false;
  }

  function navToSession(id: string) {
    router.navigateToSession(id);
  }
</script>

<div class="scheduler-page">
  <!-- ── Sidebar panel: job list ── -->
  <div class="sidebar-panel">
    <div class="controls">
      <input
        class="ctrl search-ctrl"
        type="text"
        placeholder="Search jobs..."
        bind:value={searchQuery}
      />
      <button class="create-btn" onclick={startCreate}>
        <svg width="11" height="11" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 1a.5.5 0 01.5.5V7h5.5a.5.5 0 010 1H8.5v5.5a.5.5 0 01-1 0V8H2a.5.5 0 010-1h5.5V1.5A.5.5 0 018 1z"/>
        </svg>
        New Job
      </button>
    </div>

    <div class="list-area">
      {#if loading}
        <div class="list-status">Loading...</div>
      {:else if filteredJobs.length === 0}
        <div class="empty-state">
          <div class="empty-glyph">
            <svg width="24" height="24" viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 0a8 8 0 110 16A8 8 0 018 0zm0 1.5a6.5 6.5 0 100 13 6.5 6.5 0 000-13zm-.5 2a.5.5 0 01.5.5V8a.5.5 0 00.5.5h2.5a.5.5 0 010 1H8a1 1 0 01-1-1V4a.5.5 0 01.5-.5z"/>
            </svg>
          </div>
          <span class="empty-text">
            {searchQuery ? "No matching jobs" : "No scheduled jobs"}
          </span>
        </div>
      {:else}
        {#each filteredJobs as job (job.id)}
          <button
            class="job-row"
            class:selected={selectedJob?.id === job.id}
            onclick={() => { selectedJob = job; editing = false; creating = false; showingRuns = false; }}
          >
            <span class="status-pip" class:pip-on={job.enabled} class:pip-off={!job.enabled}></span>
            <span class="row-body">
              <span class="row-title">{job.name}</span>
              <span class="row-meta">
                <span class="cron-expr">{job.cron}</span>
                <span class="cron-human">{cronHumanReadable(job.cron)}</span>
              </span>
            </span>
            <span class="row-agent">{job.agent}</span>
          </button>
        {/each}
      {/if}
    </div>
  </div>

  <!-- ── Content panel ── -->
  <main class="content-panel">
    {#if creating}
      <div class="reading-area">
        <header class="detail-header">
          <div class="header-top">
            <span class="header-badge badge-blue">New</span>
            <span class="header-title">Create Job</span>
          </div>
        </header>

        <div class="form-fields">
          <label>Name <input bind:value={formName} /></label>
          <label>
            Cron Expression
            <input bind:value={formCron} class="mono" placeholder="0 7 * * *" />
            {#if formCron}
              <span class="field-hint">{cronHumanReadable(formCron)}</span>
            {/if}
          </label>
          <label>Agent <input bind:value={formAgent} placeholder="claude" /></label>
          <label>Prompt <textarea bind:value={formPrompt} rows={4} class="mono" placeholder="Instructions for the agent..."></textarea></label>
          <label>Model <input bind:value={formModel} placeholder="anthropic/claude-sonnet-4 (optional)" /></label>
          <label>Working Directory <input bind:value={formWorkingDir} placeholder="/Users/me/project" /></label>
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
            <button class="generate-btn" onclick={handleCreate}>Create Job</button>
            <button class="cancel-btn" onclick={() => { creating = false; }}>Cancel</button>
          </div>
        </div>
      </div>

    {:else if editing && selectedJob}
      <div class="reading-area">
        <header class="detail-header">
          <div class="header-top">
            <span class="header-badge badge-blue">Edit</span>
            <span class="header-title">{selectedJob.name}</span>
          </div>
        </header>

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
          <label>Model <input bind:value={formModel} placeholder="optional" /></label>
          <label>Working Directory <input bind:value={formWorkingDir} /></label>
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
            <button class="generate-btn" onclick={handleUpdate}>Save Changes</button>
            <button class="cancel-btn" onclick={() => { editing = false; }}>Cancel</button>
            <button class="delete-btn" onclick={() => handleDelete(selectedJob!)}>Delete</button>
          </div>
        </div>
      </div>

    {:else if showingRuns && selectedJob}
      <div class="reading-area">
        <header class="detail-header">
          <div class="header-top">
            <span class="header-badge badge-purple">History</span>
            <span class="header-title">{selectedJob.name}</span>
          </div>
          <div class="header-actions">
            <button class="action-btn" onclick={() => handleRun(selectedJob!)} title="Run now">
              <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
                <path d="M4 2.5a.5.5 0 01.5.5v10a.5.5 0 01-.5.5H3a.5.5 0 01-.5-.5V3a.5.5 0 01.5-.5h1zm5.354.646l4 4a.5.5 0 010 .708l-4 4a.5.5 0 01-.708-.708L12.293 8 8.646 4.354a.5.5 0 01.708-.708z"/>
              </svg>
              Run now
            </button>
          </div>
        </header>

        {#if runs.length === 0}
          <div class="empty-prompt">
            <span>No run history yet.</span>
          </div>
        {:else}
          <div class="run-list">
            {#each runs as run (run.id)}
              <div class="run-row">
                <span class="run-pip" class:pip-completed={run.status === "completed"} class:pip-failed={run.status === "failed"} class:pip-killed={run.status === "killed"} class:pip-running={run.status === "running"}></span>
                <span class="run-cell run-time">{new Date(run.started_at).toLocaleString()}</span>
                <span class="run-cell run-status" class:status-completed={run.status === "completed"} class:status-failed={run.status === "failed"} class:status-killed={run.status === "killed"}>{run.status}</span>
                <span class="run-cell run-duration">{formatDuration(run.started_at, run.finished_at)}</span>
                <span class="run-cell run-relative">{formatRelativeTime(run.started_at)}</span>
                {#if run.session_id}
                  <button class="session-link" onclick={() => navToSession(run.session_id!)}>→ view</button>
                {/if}
                {#if run.status === "running"}
                  <button class="kill-btn" onclick={async () => { try { await api.killRun(run.id); runs = await api.listRuns(selectedJob!.id, 20); } catch { /* ignore */ }}} title="Kill this run">✕</button>
                {/if}
                {#if run.error}
                  <span class="run-error" title={run.error}>{run.error}</span>
                {/if}
              </div>
            {/each}
          </div>
        {/if}
      </div>

    {:else if selectedJob}
      <div class="reading-area">
        <header class="detail-header">
          <div class="header-top">
            <span class="header-badge" class:badge-blue={selectedJob.enabled} class:badge-muted={!selectedJob.enabled}>
              {selectedJob.enabled ? "Active" : "Disabled"}
            </span>
            <span class="header-title">{selectedJob.name}</span>
          </div>
          <div class="header-details">
            <span class="detail-chip">{selectedJob.agent}</span>
            <span class="detail-text">{selectedJob.model || "default model"}</span>
            <span class="detail-text">{selectedJob.spawn_mode}</span>
            <span class="detail-time">{formatRelativeTime(selectedJob.updated_at)}</span>
          </div>
        </header>

        <div class="job-detail-body">
          <div class="detail-field">
            <span class="detail-label">Cron</span>
            <span class="detail-value mono">{selectedJob.cron}</span>
            <span class="detail-hint">{cronHumanReadable(selectedJob.cron)}</span>
          </div>
          <div class="detail-field">
            <span class="detail-label">Prompt</span>
            <pre class="detail-prompt">{selectedJob.prompt}</pre>
          </div>
          {#if selectedJob.working_dir}
            <div class="detail-field">
              <span class="detail-label">Working Directory</span>
              <span class="detail-value mono">{selectedJob.working_dir}</span>
            </div>
          {/if}
          <div class="detail-field">
            <span class="detail-label">Inherit Context</span>
            <span class="detail-value">{selectedJob.inherit_project_context ? "Yes" : "No"}</span>
          </div>
        </div>

        <div class="detail-actions">
          <button class="action-btn" onclick={() => editJob(selectedJob!)}>Edit</button>
          <button class="action-btn" onclick={() => showRuns(selectedJob!)}>Run History</button>
          <button class="action-btn" onclick={() => handleRun(selectedJob!)}>Run Now</button>
          <button class="action-btn" onclick={() => handleToggle(selectedJob!)}>
            {selectedJob.enabled ? "Disable" : "Enable"}
          </button>
        </div>
      </div>

    {:else}
      <div class="content-empty">
        <div class="empty-prompt">
          <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 0a8 8 0 110 16A8 8 0 018 0zm0 1.5a6.5 6.5 0 100 13 6.5 6.5 0 000-13zm-.5 2a.5.5 0 01.5.5V8a.5.5 0 00.5.5h2.5a.5.5 0 010 1H8a1 1 0 01-1-1V4a.5.5 0 01.5-.5z"/>
          </svg>
          <span>Select a job or create a new one</span>
        </div>
      </div>
    {/if}
  </main>
</div>

<style>
  /* ── Layout ── (mirrors InsightsPage) */
  .scheduler-page {
    display: grid;
    grid-template-columns: 260px 1fr;
    height: calc(100vh - 40px - 24px);
    height: calc(100dvh - 40px - 24px);
    overflow: hidden;
  }

  /* ── Sidebar ── */
  .sidebar-panel {
    display: flex;
    flex-direction: column;
    border-right: 1px solid var(--border-default);
    background: var(--bg-surface);
    overflow: hidden;
  }

  /* ── Controls ── */
  .controls {
    padding: 10px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    border-bottom: 1px solid var(--border-default);
    flex-shrink: 0;
  }

  .ctrl {
    flex: 1;
    height: 28px;
    padding: 0 8px;
    background: var(--bg-inset);
    border: 1px solid var(--border-muted);
    border-radius: var(--radius-sm);
    font-size: 11px;
    color: var(--text-secondary);
    min-width: 0;
    transition: border-color 0.15s;
  }

  .ctrl:focus { outline: none; border-color: var(--accent-blue); }

  .create-btn {
    width: 100%;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 5px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: 600;
    background: var(--accent-blue);
    color: white;
    transition: opacity 0.12s;
  }

  .create-btn:hover { opacity: 0.92; }

  /* ── List Area ── */
  .list-area {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
  }

  .job-row {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    min-height: 44px;
    padding: 6px 12px;
    text-align: left;
    border-left: 2px solid transparent;
    transition: background 0.1s;
    cursor: pointer;
  }

  .job-row:hover { background: var(--bg-surface-hover); }
  .job-row.selected {
    background: var(--bg-surface-hover);
    border-left-color: var(--accent-blue);
  }

  .status-pip {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .pip-on { background: var(--accent-green, #22c55e); }
  .pip-off { background: var(--text-muted); }

  .row-body {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .row-title {
    font-size: 12px;
    font-weight: 500;
    color: var(--text-primary);
    line-height: 1.3;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .row-meta {
    display: flex;
    gap: 6px;
    font-size: 10px;
    color: var(--text-muted);
    line-height: 1.3;
  }

  .cron-expr { font-family: var(--font-mono); letter-spacing: -0.02em; }
  .cron-human { opacity: 0.7; }

  .row-agent {
    flex-shrink: 0;
    font-size: 10px;
    color: var(--text-muted);
    font-family: var(--font-mono);
    letter-spacing: -0.02em;
    white-space: nowrap;
    max-width: 55px;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .list-status {
    padding: 16px 12px;
    font-size: 11px;
    color: var(--text-muted);
    text-align: center;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 32px 16px;
    text-align: center;
  }

  .empty-glyph { color: var(--text-muted); opacity: 0.35; }
  .empty-text { font-size: 11px; color: var(--text-muted); line-height: 1.5; max-width: 160px; }

  /* ── Content Panel ── */
  .content-panel {
    overflow: hidden;
    display: flex;
    flex-direction: column;
    background: var(--bg-primary);
  }

  .reading-area {
    flex: 1;
    overflow-y: auto;
    padding: 24px 32px 48px;
  }

  .content-empty {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .empty-prompt {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    color: var(--text-muted);
    opacity: 0.5;
    font-size: 12px;
  }

  /* ── Detail Header ── */
  .detail-header {
    margin-bottom: 20px;
    padding-bottom: 14px;
    border-bottom: 1px solid var(--border-muted);
  }

  .header-top {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 8px;
    min-width: 0;
  }

  .header-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .header-badge {
    font-size: 9px;
    font-weight: 700;
    padding: 3px 8px;
    border-radius: 10px;
    color: white;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    flex-shrink: 0;
  }

  .badge-blue { background: var(--accent-blue); }
  .badge-purple { background: var(--accent-purple); }
  .badge-muted { background: var(--text-muted); }

  .header-details {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    color: var(--text-muted);
  }

  .detail-chip {
    padding: 1px 6px;
    border-radius: var(--radius-sm);
    background: var(--bg-inset);
    color: var(--text-secondary);
    font-size: 11px;
  }

  .detail-text { color: var(--text-muted); }
  .detail-time { margin-left: auto; font-variant-numeric: tabular-nums; }

  .header-actions { margin-top: 6px; display: flex; gap: 6px; }

  /* ── Detail Body ── */
  .job-detail-body {
    display: flex;
    flex-direction: column;
    gap: 14px;
    margin-bottom: 20px;
  }

  .detail-field {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .detail-label {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .detail-value {
    font-size: 13px;
    color: var(--text-primary);
  }

  .detail-value.mono { font-family: var(--font-mono); font-size: 12px; }

  .detail-hint {
    font-size: 11px;
    color: var(--text-muted);
  }

  .detail-prompt {
    font-size: 12px;
    color: var(--text-secondary);
    background: var(--bg-inset);
    padding: 8px 10px;
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-muted);
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.5;
    max-height: 200px;
    overflow-y: auto;
  }

  .detail-actions {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
  }

  .action-btn {
    height: 28px;
    padding: 0 12px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
    background: var(--bg-inset);
    border: 1px solid var(--border-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
    transition: background 0.1s, color 0.1s;
  }

  .action-btn:hover {
    background: var(--bg-surface-hover);
    color: var(--text-primary);
  }

  /* ── Form ── */
  .form-fields {
    display: flex;
    flex-direction: column;
    gap: 10px;
    max-width: 480px;
  }

  .form-fields label {
    font-size: 11px;
    font-weight: 500;
    color: var(--text-secondary);
    display: block;
  }

  .form-fields input, .form-fields textarea, .form-fields select {
    display: block;
    width: 100%;
    margin-top: 3px;
    padding: 6px 8px;
    font-size: 12px;
    border: 1px solid var(--border-muted);
    border-radius: var(--radius-sm);
    background: var(--bg-inset);
    color: var(--text-primary);
    box-sizing: border-box;
    transition: border-color 0.15s;
  }

  .form-fields input:focus, .form-fields textarea:focus, .form-fields select:focus {
    outline: none;
    border-color: var(--accent-blue);
  }

  .form-fields .mono { font-family: var(--font-mono); }

  .field-hint {
    font-size: 10px;
    color: var(--text-muted);
    margin-top: 2px;
    display: block;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
  }

  .checkbox-label input { width: auto; margin: 0; }

  .form-actions {
    display: flex;
    gap: 6px;
    margin-top: 6px;
    align-items: center;
  }

  .generate-btn {
    height: 28px;
    padding: 0 14px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: 600;
    background: var(--accent-blue);
    color: white;
    transition: opacity 0.12s;
  }

  .generate-btn:hover { opacity: 0.92; }

  .cancel-btn {
    height: 28px;
    padding: 0 14px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    color: var(--text-secondary);
    background: var(--bg-inset);
    border: 1px solid var(--border-muted);
    transition: background 0.1s;
  }

  .cancel-btn:hover { background: var(--bg-surface-hover); }

  .delete-btn {
    margin-left: auto;
    height: 28px;
    padding: 0 12px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    color: var(--accent-red, #ef4444);
    background: transparent;
    border: 1px solid transparent;
    transition: background 0.1s;
  }

  .delete-btn:hover {
    background: color-mix(in srgb, var(--accent-red) 10%, transparent);
  }

  /* ── Run List ── */
  .run-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .run-row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: var(--bg-surface);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-sm);
    font-size: 11px;
  }

  .run-pip {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .pip-completed { background: var(--accent-green, #22c55e); }
  .pip-failed { background: var(--accent-red, #ef4444); }
  .pip-killed { background: var(--accent-amber, #f59e0b); }
  .pip-running { background: var(--accent-blue, #3b82f6); }

  .run-cell {
    color: var(--text-secondary);
  }

  .run-time { white-space: nowrap; }
  .run-status { font-weight: 500; }
  .status-completed { color: var(--accent-green, #22c55e); }
  .status-failed { color: var(--accent-red, #ef4444); }
  .status-killed { color: var(--accent-amber, #f59e0b); }

  .run-duration, .run-relative { color: var(--text-muted); }
  .run-relative { flex: 1; min-width: 0; }

  .session-link {
    font-size: 11px;
    color: var(--accent-blue);
    text-decoration: none;
    white-space: nowrap;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
  }

  .session-link:hover { text-decoration: underline; }

  .kill-btn {
    background: none;
    border: 1px solid var(--accent-red, #ef4444);
    color: var(--accent-red, #ef4444);
    border-radius: 3px;
    font-size: 10px;
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    padding: 0;
    flex-shrink: 0;
    transition: background 0.12s;
  }

  .kill-btn:hover {
    background: color-mix(in srgb, var(--accent-red) 15%, transparent);
  }

  .run-error {
    color: var(--accent-red);
    font-size: 10px;
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
