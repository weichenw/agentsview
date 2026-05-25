<script lang="ts">
  import { onMount } from "svelte";

  interface LogFile {
    name: string;
    path: string;
    lines: string[];
    total_lines: number;
  }

  let files: LogFile[] = $state([]);
  let activeFile = $state(0);
  let loading = $state(true);
  let error = $state("");
  let autoScroll = $state(true);
  let logRef: HTMLDivElement | undefined = $state(undefined);

  const POLL_MS = 5000;
  const LINES = 1000;

  async function fetchLogs() {
    try {
      const res = await fetch(`/api/v1/logs?lines=${LINES}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      files = data.files ?? [];
      // Keep active index valid if files changed.
      if (activeFile >= files.length) {
        activeFile = 0;
      }
      error = "";
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    fetchLogs();
    const interval = window.setInterval(fetchLogs, POLL_MS);
    return () => window.clearInterval(interval);
  });

  $effect(() => {
    if (autoScroll && logRef && files[activeFile]?.lines.length > 0) {
      logRef.scrollTop = logRef.scrollHeight;
    }
  });

  function copyAll() {
    const text = files[activeFile]?.lines.join("\n") ?? "";
    navigator.clipboard.writeText(text).catch(() => {});
  }

  const active = $derived(files[activeFile]);
</script>

<div class="logs-page">
  <header class="logs-header">
    <div class="logs-meta">
      <span class="logs-title">Logs</span>
      {#if active}
        <span class="logs-subtitle" title={active.path}>{active.name}</span>
        <span class="logs-count">{active.total_lines.toLocaleString()} lines</span>
      {/if}
    </div>
    <div class="logs-actions">
      <label class="logs-toggle">
        <input type="checkbox" bind:checked={autoScroll} />
        Auto-scroll
      </label>
      <button class="logs-btn" onclick={copyAll}>Copy</button>
      <button class="logs-btn" onclick={fetchLogs}>Refresh</button>
    </div>
  </header>

  {#if files.length > 1}
    <div class="logs-tabs">
      {#each files as f, i}
        <button
          class="logs-tab"
          class:active={i === activeFile}
          onclick={() => activeFile = i}
        >
          {f.name}
          <span class="logs-tab-count">{f.total_lines}</span>
        </button>
      {/each}
    </div>
  {/if}

  {#if error}
    <div class="logs-error">{error}</div>
  {/if}

  <div class="logs-body" bind:this={logRef}>
    {#if loading && (!active || active.lines.length === 0)}
      <div class="logs-empty">Loading logs...</div>
    {:else if !active || active.lines.length === 0}
      <div class="logs-empty">No log lines yet.</div>
    {:else}
      {#each active.lines as line, i}
        <div class="logs-line">
          <span class="logs-line-num">{active.total_lines - active.lines.length + i + 1}</span>
          <span class="logs-line-text">{line}</span>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .logs-page {
    display: flex;
    flex-direction: column;
    height: 100%;
    overflow: hidden;
    background: var(--bg-canvas);
  }

  .logs-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    border-bottom: 1px solid var(--border-default);
    background: var(--bg-surface);
    flex-shrink: 0;
    gap: 10px;
  }

  .logs-meta {
    display: flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
  }

  .logs-title {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .logs-subtitle {
    font-size: 10px;
    color: var(--text-muted);
    font-family: var(--font-mono);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 300px;
  }

  .logs-count {
    font-size: 10px;
    color: var(--text-muted);
  }

  .logs-actions {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-shrink: 0;
  }

  .logs-toggle {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 10px;
    color: var(--text-muted);
    cursor: pointer;
    user-select: none;
  }

  .logs-toggle input {
    width: 12px;
    height: 12px;
  }

  .logs-btn {
    height: 24px;
    padding: 0 10px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
    background: var(--bg-inset);
    border: 1px solid var(--border-default);
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
  }

  .logs-btn:hover {
    background: var(--bg-surface-hover);
    color: var(--text-primary);
  }

  .logs-tabs {
    display: flex;
    gap: 2px;
    padding: 6px 14px 0;
    border-bottom: 1px solid var(--border-default);
    background: var(--bg-surface);
    flex-shrink: 0;
  }

  .logs-tab {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 4px 10px;
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    border-bottom: 2px solid transparent;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
  }

  .logs-tab:hover {
    background: var(--bg-surface-hover);
    color: var(--text-primary);
  }

  .logs-tab.active {
    color: var(--accent-blue);
    border-bottom-color: var(--accent-blue);
    background: color-mix(in srgb, var(--accent-blue) 6%, transparent);
  }

  .logs-tab-count {
    font-size: 9px;
    padding: 1px 4px;
    border-radius: var(--radius-sm);
    background: var(--bg-inset);
    color: var(--text-muted);
  }

  .logs-error {
    padding: 8px 14px;
    font-size: 11px;
    color: var(--accent-red);
    background: color-mix(in srgb, var(--accent-red) 6%, transparent);
    border-bottom: 1px solid var(--border-default);
    flex-shrink: 0;
  }

  .logs-body {
    flex: 1;
    overflow-y: auto;
    padding: 8px 0;
    font-family: var(--font-mono);
    font-size: 11px;
    line-height: 1.5;
    color: var(--text-secondary);
  }

  .logs-empty {
    padding: 40px 14px;
    text-align: center;
    color: var(--text-muted);
    font-size: 12px;
  }

  .logs-line {
    display: flex;
    gap: 8px;
    padding: 1px 14px;
    white-space: pre;
  }

  .logs-line:hover {
    background: var(--bg-surface-hover);
  }

  .logs-line-num {
    color: var(--text-muted);
    text-align: right;
    min-width: 44px;
    flex-shrink: 0;
    user-select: none;
  }

  .logs-line-text {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
