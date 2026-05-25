<script lang="ts">
  import { onMount } from "svelte";
  import { router } from "../../stores/router.svelte.ts";

  let graphData = $state({ nodes: [] as any[], links: [] as any[], stats: {} as any });
  let loading = $state(true);
  let error = $state("");
  let selectedNodeId = $state<string | null>(null);
  let showLabels = $state(true);
  let showDetail = $state(false);
  let detailNode = $state<any>(null);
  let searchQuery = $state("");
  let d3Ready = $state(false);

  let filterTypes = $state({
    session: true,
    project: true,
    tool: true,
    domain: true,
    category: true,
  });

  let svgRef: SVGSVGElement;
  let wrapRef: HTMLDivElement;
  let width = $state(0);
  let height = $state(0);

  // D3 state (not reactive)
  let svg: any = null;
  let g: any = null;
  let zoom: any = null;
  let simulation: any = null;
  let linkSel: any = null;
  let nodeSel: any = null;
  let labelSel: any = null;

  const colorMap: Record<string, string> = {
    session:  "#8b5cf6",
    memory:   "#f59e0b",
    project:  "#3b82f6",
    domain:   "#a855f7",
    category: "#14b8a6",
    hub:      "#f97316",
  };

  // Load D3 from CDN dynamically once.
  async function loadD3() {
    if ((window as any).d3) return;
    return new Promise<void>((resolve, reject) => {
      const script = document.createElement("script");
      script.src = "/d3.v7.min.js";
      script.async = true;
      const timeout = setTimeout(() => reject(new Error("D3 load timeout")), 10000);
      script.onload = () => { clearTimeout(timeout); resolve(); };
      script.onerror = () => { clearTimeout(timeout); reject(new Error("Failed to load D3")); };
      document.head.appendChild(script);
    });
  }

  async function fetchGraph() {
    try {
      const types = Object.entries(filterTypes)
        .filter(([, v]) => v)
        .map(([k]) => k)
        .join(",");
      const res = await fetch(`/api/v1/memory/graph?types=${encodeURIComponent(types)}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      graphData = await res.json();
      loading = false;
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
      loading = false;
    }
  }

  onMount(() => {
    init();
    window.addEventListener("resize", resize);
    return () => {
      window.removeEventListener("resize", resize);
      if (simulation) simulation.stop();
    };
  });

  async function init() {
    try {
      await loadD3();
      d3Ready = true;
    } catch (e) {
      // D3 failed, but still try to show data in table view
      console.warn("D3 not available:", e);
    }
    await fetchGraph();
    resize();
    if (d3Ready) render();
  }

  function resize() {
    if (!wrapRef) return;
    const rect = wrapRef.getBoundingClientRect();
    width = rect.width;
    height = rect.height;
    if (svg) {
      svg.attr("width", width).attr("height", height);
      if (simulation) simulation.force("center", (window as any).d3.forceCenter(width / 2, height / 2)).alpha(0.3).restart();
    }
  }

  function render() {
    if (!d3Ready || !graphData.nodes.length || !svgRef) return;
    const d3 = (window as any).d3;
    if (!d3) return;

    if (simulation) simulation.stop();

    svg = d3.select(svgRef);
    svg.selectAll("*").remove();
    g = svg.append("g");

    zoom = d3.zoom().scaleExtent([0.1, 5]).on("zoom", (e: any) => g.attr("transform", e.transform));
    svg.call(zoom);

    simulation = d3.forceSimulation(graphData.nodes)
      .force("link", d3.forceLink(graphData.links).id((d: any) => d.id).distance((d: any) => {
        if (d.source.type === "hub" || d.target.type === "hub") return 220;
        if (d.source.type === "domain" || d.target.type === "domain") return 180;
        if (d.source.type === "category" || d.target.type === "category") return 150;
        if (d.source.type === "memory" || d.target.type === "memory") return 160;
        return 120 + 80 / (d.value || 1);
      }).strength(0.3))
      .force("charge", d3.forceManyBody().strength((d: any) => {
        if (d.type === "memory") return -500;
        if (d.type === "session") return -400;
        if (d.type === "project") return -600;
        if (d.type === "hub") return -800;
        if (d.type === "domain") return -500;
        if (d.type === "category") return -400;
        return -200;
      }).theta(0.9))
      .force("center", d3.forceCenter(width / 2, height / 2).strength(0.02))
      .force("collide", d3.forceCollide().radius((d: any) => d.r + 12).strength(1.0).iterations(3))
      .force("x", d3.forceX(width / 2).strength(0.01))
      .force("y", d3.forceY(height / 2).strength(0.01))
      .alphaDecay(0.02)
      .velocityDecay(0.4)
      .on("tick", ticked);

    linkSel = g.append("g").attr("class", "links")
      .selectAll("line").data(graphData.links).join("line")
      .attr("class", "Mem-link").attr("stroke-width", (d: any) => Math.max(1, Math.sqrt(d.value)));

    nodeSel = g.append("g").attr("class", "nodes")
      .selectAll("circle").data(graphData.nodes).join("circle")
      .attr("class", "Mem-node").attr("r", (d: any) => d.r).attr("fill", (d: any) => colorMap[d.type] || "#999")
      .call(d3.drag().on("start", dragstarted).on("drag", dragged).on("end", dragended));

    labelSel = g.append("g").attr("class", "labels")
      .selectAll("text").data(graphData.nodes).join("text")
      .attr("class", "Mem-label").text((d: any) => d.label)
      .attr("x", (d: any) => d.r + 3).attr("y", 3)
      .style("display", showLabels ? "block" : "none");

    nodeSel.on("mouseover", (event: any, d: any) => showTooltip(event, d))
      .on("mouseout", () => hideTooltip())
      .on("click", (event: any, d: any) => { event.stopPropagation(); selectNode(d); });
    svg.on("click", () => { selectedNodeId = null; showDetail = false; updateHighlight(); });

    function ticked() {
      linkSel.attr("x1", (d: any) => d.source.x).attr("y1", (d: any) => d.source.y)
             .attr("x2", (d: any) => d.target.x).attr("y2", (d: any) => d.target.y);
      nodeSel.attr("cx", (d: any) => d.x).attr("cy", (d: any) => d.y);
      labelSel.attr("x", (d: any) => d.x + d.r + 3).attr("y", (d: any) => d.y + 3);
    }
    function dragstarted(event: any, d: any) { if (!event.active) simulation.alphaTarget(0.3).restart(); d.fx = d.x; d.fy = d.y; }
    function dragged(event: any, d: any) { d.fx = event.x; d.fy = event.y; }
    function dragended(event: any, d: any) { if (!event.active) simulation.alphaTarget(0); d.fx = null; d.fy = null; }
  }

  let tooltip = $state({ x: 0, y: 0, text: "", show: false });
  function showTooltip(event: MouseEvent, d: any) {
    const rect = wrapRef.getBoundingClientRect();
    tooltip = {
      x: event.clientX - rect.left + 12,
      y: event.clientY - rect.top + 12,
      text: `${d.type}\n${d.label}${d.count ? "\ncount: " + d.count : ""}`,
      show: true,
    };
  }
  function hideTooltip() {
    tooltip = { ...tooltip, show: false };
  }

  function fitGraph() {
    if (!graphData.nodes.length || !svg) return;
    const xs = graphData.nodes.map((d: any) => d.x);
    const ys = graphData.nodes.map((d: any) => d.y);
    const pad = 40;
    const minX = Math.min(...xs) - pad, maxX = Math.max(...xs) + pad;
    const minY = Math.min(...ys) - pad, maxY = Math.max(...ys) + pad;
    const scale = Math.min(width / (maxX - minX), height / (maxY - minY), 2);
    const cx = (minX + maxX) / 2, cy = (minY + maxY) / 2;
    const d3 = (window as any).d3;
    svg.transition().duration(600).call(zoom.transform, d3.zoomIdentity.translate(width / 2, height / 2).scale(Math.max(scale, 0.2)).translate(-cx, -cy));
  }

  function resetGraph() {
    searchQuery = "";
    selectedNodeId = null;
    showDetail = false;
    if (svg) {
      const d3 = (window as any).d3;
      svg.transition().duration(500).call(zoom.transform, d3.zoomIdentity);
    }
  }

  function updateHighlight() {
    if (!nodeSel || !linkSel || !labelSel) return;
    if (!selectedNodeId) {
      nodeSel.classed("dim", false).classed("highlight", false);
      linkSel.classed("dim", false).classed("highlight", false);
      labelSel.classed("dim", false);
      return;
    }
    const con = new Set([selectedNodeId]);
    graphData.links.forEach((l: any) => {
      if (l.source.id === selectedNodeId || l.target.id === selectedNodeId) {
        con.add(l.source.id); con.add(l.target.id);
      }
    });
    nodeSel.classed("dim", (d: any) => !con.has(d.id)).classed("highlight", (d: any) => d.id === selectedNodeId);
    linkSel.classed("dim", (d: any) => {
      const sid = d.source.id || d.source;
      const tid = d.target.id || d.target;
      return sid !== selectedNodeId && tid !== selectedNodeId;
    }).classed("highlight", (d: any) => {
      const sid = d.source.id || d.source;
      const tid = d.target.id || d.target;
      return sid === selectedNodeId || tid === selectedNodeId;
    });
    labelSel.classed("dim", (d: any) => !con.has(d.id));
  }

  function selectNode(d: any) {
    selectedNodeId = d.id;
    showDetail = true;
    detailNode = d;
    updateHighlight();
  }

  function filteredNodes() {
    const q = searchQuery.toLowerCase().trim();
    if (!q) return;
    if (!nodeSel) return;
    nodeSel.style("opacity", (d: any) => (d.label.toLowerCase().includes(q) ? 1 : 0.15));
    linkSel?.style("opacity", 0.05);
    labelSel?.style("opacity", (d: any) => (d.label.toLowerCase().includes(q) ? 1 : 0.1));
  }

  function clearFilter() {
    searchQuery = "";
    if (!nodeSel) return;
    nodeSel.style("opacity", 1);
    linkSel?.style("opacity", null);
    labelSel?.style("opacity", 1);
  }

  function toggleFilter(key: string) {
    filterTypes = { ...filterTypes, [key]: !filterTypes[key] };
    fetchGraph().then(() => { resize(); if (d3Ready) render(); });
  }

  function navToSession(nodeId: string) {
    // Strip "s:" prefix from graph node ID. Pi DB stores session IDs
    // without the agent prefix, but agents-view expects "pi:<uuid>".
    let sessionId = nodeId.startsWith("s:") ? nodeId.slice(2) : nodeId;
    if (!sessionId.includes(":")) {
      sessionId = "pi:" + sessionId;
    }
    router.navigateToSession(sessionId);
  }

  $effect(() => {
    if (!searchQuery) clearFilter();
    else filteredNodes();
  });

  $effect(() => {
    if (showLabels !== undefined && labelSel) {
      labelSel.style("display", showLabels ? "block" : "none");
    }
  });
</script>

<div class="memory-page">
  <div class="memory-sidebar">
    <div class="memory-search">
      <input type="text" placeholder="Search nodes..." bind:value={searchQuery} />
    </div>
    <div class="memory-panel">
      <div class="memory-section-title">Node Types</div>
      <label class="memory-filter-item">
        <input type="checkbox" checked disabled />
        <span>Memory</span>
      </label>
      {#each Object.entries(filterTypes) as [key, val]}
        <label class="memory-filter-item">
          <input
            type="checkbox"
            checked={val}
            onchange={() => toggleFilter(key)}
          />
          <span>{key.charAt(0).toUpperCase() + key.slice(1)}</span>
        </label>
      {/each}
    </div>
    <div class="memory-panel">
      <div class="memory-section-title">Stats</div>
      {#if graphData.stats}
        <div class="memory-stat">{graphData.stats.nodes} nodes</div>
        <div class="memory-stat">{graphData.stats.links} links</div>
        <div class="memory-stat">{graphData.stats.messages} messages</div>
        <div class="memory-stat">{graphData.stats.memories} memories</div>
      {/if}
    </div>
    {#if showDetail && detailNode}
      <div class="memory-detail">
        <div class="memory-detail-header">
          <span class="memory-detail-type" style="background: {colorMap[detailNode.type]}22; color: {colorMap[detailNode.type]}">{detailNode.type}</span>
          <span class="memory-detail-title">{detailNode.label}</span>
          <button class="memory-detail-close" onclick={() => { showDetail = false; selectedNodeId = null; updateHighlight(); }}>&times;</button>
        </div>
        <div class="memory-detail-body">
          {#if detailNode.type === "session"}
            <div class="memory-field">
              <button class="memory-session-link" onclick={() => navToSession(detailNode.id)}>
                → View session
              </button>
            </div>
          {/if}
          <div class="memory-field">
            <div class="memory-field-label">ID</div>
            <div class="memory-field-value">{detailNode.id}</div>
          </div>
          {#if detailNode.count}
            <div class="memory-field">
              <div class="memory-field-label">Count</div>
              <div class="memory-field-value">{detailNode.count}</div>
            </div>
          {/if}
          {#if detailNode.time}
            <div class="memory-field">
              <div class="memory-field-label">Time</div>
              <div class="memory-field-value">{detailNode.time}</div>
            </div>
          {/if}
          {#if detailNode.db}
            <div class="memory-field">
              <div class="memory-field-label">Database</div>
              <div class="memory-field-value">{detailNode.db}</div>
            </div>
          {/if}
          {#if detailNode.raw && detailNode.raw.content}
            <div class="memory-field">
              <div class="memory-field-label">Content</div>
              <pre class="memory-field-pre">{detailNode.raw.content.slice(0, 2000)}</pre>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
  <div class="memory-graph-wrap" bind:this={wrapRef}>
    {#if loading}
      <div class="memory-loading">Loading memory graph...</div>
    {:else if error}
      <div class="memory-error">{error}</div>
    {:else if !graphData.nodes.length}
      <div class="memory-empty">No memory data found.</div>
    {:else}
      <svg class="memory-svg" bind:this={svgRef}></svg>
      <div class="memory-controls">
        <button onclick={fitGraph}>Fit</button>
        <button onclick={resetGraph}>Reset</button>
        <button onclick={() => showLabels = !showLabels}>{showLabels ? "Hide" : "Show"} Labels</button>
      </div>
    {/if}
    {#if tooltip.show}
      <div class="memory-tooltip" style="left: {tooltip.x}px; top: {tooltip.y}px">
        {tooltip.text}
      </div>
    {/if}
  </div>
</div>

<style>
  .memory-page {
    display: flex;
    height: 100%;
    overflow: hidden;
    background: var(--bg-canvas);
  }

  .memory-sidebar {
    width: 300px;
    min-width: 260px;
    max-width: 400px;
    background: var(--bg-surface);
    border-right: 1px solid var(--border-default);
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    flex-shrink: 0;
    z-index: 10;
  }

  .memory-search {
    padding: 10px 12px;
    border-bottom: 1px solid var(--border-default);
  }

  .memory-search input {
    width: 100%;
    padding: 6px 10px;
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-default);
    background: var(--bg-inset);
    color: var(--text-primary);
    font-size: 12px;
    outline: none;
  }

  .memory-search input:focus {
    border-color: var(--accent-blue);
  }

  .memory-filter-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 2px 0;
    font-size: 12px;
    color: var(--text-primary);
    cursor: pointer;
    user-select: none;
  }

  .memory-filter-item:hover {
    color: var(--text-primary);
  }

  .memory-filter-item input[type="checkbox"] {
    accent-color: var(--accent-blue);
    cursor: pointer;
    margin: 0;
  }

  .memory-panel {
    padding: 10px 12px;
    border-bottom: 1px solid var(--border-default);
  }

  .memory-section-title {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--text-muted);
    margin-bottom: 6px;
  }

  .memory-stat {
    font-size: 11px;
    color: var(--text-secondary);
    margin-bottom: 3px;
  }

  .memory-detail {
    padding: 12px;
    border-bottom: 1px solid var(--border-default);
    flex-shrink: 0;
  }

  .memory-detail-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 10px;
  }

  .memory-detail-type {
    text-transform: uppercase;
    font-size: 10px;
    letter-spacing: 0.04em;
    padding: 2px 8px;
    border-radius: 4px;
  }

  .memory-detail-title {
    font-weight: 600;
    font-size: 13px;
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .memory-detail-close {
    background: none;
    border: none;
    font-size: 18px;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0;
    line-height: 1;
  }

  .memory-detail-body {
    font-size: 12px;
    line-height: 1.5;
  }

  .memory-field {
    margin-bottom: 10px;
  }

  .memory-field-label {
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--text-muted);
    margin-bottom: 2px;
  }

  .memory-field-value {
    color: var(--text-primary);
    word-break: break-word;
  }

  .memory-field-pre {
    background: var(--bg-inset);
    border: 1px solid var(--border-default);
    border-radius: 6px;
    padding: 8px;
    font-size: 11px;
    color: var(--text-secondary);
    overflow-x: auto;
    white-space: pre-wrap;
    word-break: break-word;
    max-height: 300px;
    overflow-y: auto;
  }

  .memory-session-link {
    background: var(--accent-blue);
    color: white;
    border: none;
    padding: 5px 12px;
    border-radius: var(--radius-sm);
    font-size: 12px;
    cursor: pointer;
    transition: background 0.15s;
  }

  .memory-session-link:hover {
    background: var(--accent-blue-hover, #2563eb);
  }

  .memory-graph-wrap {
    flex: 1;
    position: relative;
    overflow: hidden;
  }

  .memory-svg {
    width: 100%;
    height: 100%;
    cursor: grab;
    display: block;
  }

  .memory-svg:active {
    cursor: grabbing;
  }

  .memory-loading,
  .memory-empty {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 13px;
    color: var(--text-muted);
  }

  .memory-error {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 13px;
    color: var(--accent-red);
    padding: 20px;
  }

  .memory-controls {
    position: absolute;
    bottom: 14px;
    left: 14px;
    display: flex;
    gap: 6px;
    z-index: 50;
  }

  .memory-controls button {
    background: var(--bg-surface);
    border: 1px solid var(--border-default);
    color: var(--text-secondary);
    padding: 5px 10px;
    border-radius: var(--radius-sm);
    font-size: 11px;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
  }

  .memory-controls button:hover {
    background: var(--bg-surface-hover);
    color: var(--text-primary);
  }

  .memory-tooltip {
    position: absolute;
    background: var(--bg-surface);
    border: 1px solid var(--border-default);
    padding: 5px 10px;
    border-radius: 6px;
    font-size: 11px;
    color: var(--text-primary);
    pointer-events: none;
    z-index: 100;
    max-width: 240px;
    box-shadow: var(--shadow-lg);
    white-space: pre-line;
  }

  :global(.Mem-node) {
    cursor: pointer;
    stroke: none;
  }

  :global(.Mem-node.dim) {
    opacity: 0.15;
  }

  :global(.Mem-node.highlight) {
    stroke: var(--text-primary);
    stroke-width: 2px;
    filter: brightness(1.2);
  }

  :global(.Mem-link) {
    stroke: var(--text-muted);
    stroke-opacity: 0.6;
  }

  :global(.Mem-link.dim) {
    opacity: 0.05;
  }

  :global(.Mem-link.highlight) {
    stroke: var(--accent-blue);
    stroke-opacity: 0.9;
  }

  :global(.Mem-label) {
    font-size: 11px;
    fill: var(--text-primary);
    pointer-events: none;
    text-shadow: 0 1px 2px rgba(255,255,255,0.9);
  }

  :global(.Mem-label.dim) {
    opacity: 0.1;
  }

  @media (max-width: 767px) {
    .memory-sidebar {
      display: none;
    }
  }
</style>
