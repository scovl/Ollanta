"use strict";(()=>{function n(e){return e.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;")}function V(e){return[{key:"details",label:"Details"},{key:"rule",label:"Rule"},{key:"ai-fix",label:"Fix with AI"}].map(t=>`<button class="detail-tab${e===t.key?" active":""}" data-detail-tab="${t.key}">${t.label}</button>`).join("")}var L,c=[],g=[],f=[],p=null,b=-1,U="overview",$="details",I="",k=!1,v=null,s=G(),h="all",C="all",P="all",j="",z={blocker:0,critical:1,major:2,minor:3,info:4},S={blocker:"#ef4444",critical:"#f97316",major:"#eab308",minor:"#22c55e",info:"#64748b"},q={bug:"Bug",code_smell:"Code Smell",vulnerability:"Vulnerability",security_hotspot:"Hotspot"};async function ee(){try{let e=await fetch("/report.json");if(!e.ok)throw new Error(`HTTP ${e.status}`);L=await e.json(),c=L.issues??[],te(),ae(),se(),ne(),le(),oe(),de(),w(),pe(),X(),re(),xe(),l("tab-issue-count").textContent=String(c.length),l("tab-file-count").textContent=String(new Set(c.map(i=>i.component_path)).size)}catch(e){l("app").innerHTML=`<div class="error">Failed to load report: ${String(e)}</div>`}}document.addEventListener("DOMContentLoaded",ee);function te(){let e=L.metadata,i=new Date(e.analysis_date).toLocaleString();l("project-key").textContent=e.project_key,l("scan-date").textContent=i,l("scan-version").textContent=`v${e.version}`,l("elapsed").textContent=`${e.elapsed_ms}ms`}function ie(){let e=L.measures,i=[{metric:"Bugs",operator:"=",threshold:0,value:e.bugs,passed:e.bugs===0},{metric:"Vulnerabilities",operator:"=",threshold:0,value:e.vulnerabilities,passed:e.vulnerabilities===0}];return{status:i.every(a=>a.passed)?"passed":"failed",conditions:i}}function ae(){let e=ie(),i=l("gate-hero");i.classList.remove("gate-loading"),i.classList.add(e.status==="passed"?"gate-passed":"gate-failed"),l("gate-icon").textContent=e.status==="passed"?"\u2713":"\u2717",l("gate-status").textContent=e.status==="passed"?"Passed":"Failed";let t=e.conditions.map(a=>{let o=a.passed?"cond-pass":"cond-fail",r=a.passed?"\u2713":"\u2717";return`<div class="gate-cond ${o}">
      <span class="gate-cond-icon">${r}</span>
      <span class="gate-cond-metric">${n(a.metric)}</span>
      <span class="gate-cond-value">${a.value}</span>
    </div>`}).join("");l("gate-conditions").innerHTML=t}function se(){let e=L.measures;x("m-bugs",e.bugs),x("m-vulns",e.vulnerabilities),x("m-smells",e.code_smells),x("m-ncloc",e.ncloc),x("m-files",e.files),x("m-comments",e.comments),F("card-bugs",e.bugs,[0,1,5]),F("card-vulns",e.vulnerabilities,[0,1,3]),F("card-smells",e.code_smells,[0,10,50]),A("card-ncloc","card-neutral"),A("card-files","card-neutral"),A("card-comments","card-neutral")}function x(e,i){l(e).textContent=i.toLocaleString()}function F(e,i,t){i<=t[0]?A(e,"card-green"):i<=t[1]?A(e,"card-yellow"):A(e,"card-red")}function ne(){let e=M(c,d=>d.severity),i=Math.max(1,...Object.values(e)),t=["blocker","critical","major","minor","info"],a="",o="",r=c.length||1;for(let d of t){let y=e[d]??0,E=y/i*100,T=S[d]??"#64748b";a+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${d}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${E}%;background:${T}"></div></div>
      <span class="sev-bar-count">${y}</span>
    </div>`,y>0&&(o+=`<div class="sev-segment" style="width:${y/r*100}%;background:${T}" title="${d}: ${y}"></div>`)}l("sev-bars").innerHTML=a,l("sev-proportional").innerHTML=o;let u=M(c,d=>d.type),H=Math.max(1,...Object.values(u)),Q={bug:"#ef4444",vulnerability:"#f97316",code_smell:"#22c55e",security_hotspot:"#eab308"},B="";for(let[d,y]of Object.entries(q)){let E=u[d]??0,T=E/H*100,Z=Q[d]??"#64748b";B+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${y}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${T}%;background:${Z}"></div></div>
      <span class="sev-bar-count">${E}</span>
    </div>`}l("type-bars").innerHTML=B}function le(){let e=M(c,t=>t.component_path),i=Object.entries(e).sort((t,a)=>a[1]-t[1]).slice(0,10);if(!i.length){l("hotspot-files").innerHTML='<div class="empty-state">No issues found</div>';return}l("hotspot-files").innerHTML=i.map(([t,a])=>{let o=N(t);return`<div class="hotspot-row" data-path="${n(t)}">
      <span class="hotspot-file" title="${n(t)}">${n(o)}</span>
      <span class="hotspot-count">${a}</span>
    </div>`}).join(""),l("hotspot-files").querySelectorAll(".hotspot-row").forEach(t=>{t.addEventListener("click",()=>{let a=t.dataset.path;J("files"),ve(a)})})}function oe(){let e=Object.entries(L.measures.by_language).sort((t,a)=>a[1]-t[1]),i=Math.max(1,e[0]?.[1]??1);if(!e.length){l("by-lang").innerHTML='<span class="empty-state">No language data</span>';return}l("by-lang").innerHTML=e.map(([t,a])=>`<div class="lang-row">
      <span class="lang-name">${n(t)}</span>
      <div class="lang-bar-track"><div class="lang-bar-fill" style="width:${a/i*100}%"></div></div>
      <span class="lang-count">${a} files</span>
    </div>`).join("")}function re(){document.querySelectorAll(".tab").forEach(e=>{e.addEventListener("click",()=>{let i=e.dataset.tab;J(i)})})}function J(e){U=e,document.querySelectorAll(".tab").forEach(i=>i.classList.remove("active")),document.querySelector(`.tab[data-tab="${e}"]`)?.classList.add("active"),document.querySelectorAll(".panel").forEach(i=>i.classList.add("hidden")),l(`panel-${e}`).classList.remove("hidden")}function de(){let e=[...new Set(c.map(t=>t.rule_key))].sort((t,a)=>t.localeCompare(a)),i=l("filter-rule");e.forEach(t=>{let a=document.createElement("option");a.value=t,a.textContent=t,i.appendChild(a)}),l("filter-severity").addEventListener("change",t=>{h=t.target.value,w()}),l("filter-type").addEventListener("change",t=>{C=t.target.value,w()}),i.addEventListener("change",t=>{P=t.target.value,w()}),l("search").addEventListener("input",t=>{j=t.target.value.toLowerCase(),w()}),W()}function W(){let e=M(c,t=>t.severity),i=["blocker","critical","major","minor","info"];l("sev-chips").innerHTML=i.map(t=>{let a=e[t]??0,o=S[t]??"#64748b";return`<div class="sev-chip${h===t?" active":""}" data-sev="${t}"
      style="--chip-color:${o};--chip-bg:${o}15">
      <span class="chip-dot" style="background:${o}"></span>
      ${t}
      <span class="chip-count">${a}</span>
    </div>`}).join(""),l("sev-chips").querySelectorAll(".sev-chip").forEach(t=>{t.addEventListener("click",()=>{let a=t.dataset.sev;h=h===a?"all":a,l("filter-severity").value=h,w(),W()})})}function w(){g=c.filter(e=>!(h!=="all"&&e.severity!==h||C!=="all"&&e.type!==C||P!=="all"&&e.rule_key!==P||j&&!`${e.component_path} ${e.message} ${e.rule_key}`.toLowerCase().includes(j))),g.sort((e,i)=>{let t=z[e.severity]??99,a=z[i.severity]??99;return t-a}),b=-1,ce()}function ce(){let e=l("issue-list"),i=g.length===1?"issue":"issues";if(l("issue-count").textContent=`${g.length} ${i}`,!g.length){e.innerHTML='<div class="empty-state">No issues match the current filters.</div>';return}e.innerHTML=g.map((t,a)=>{let o=S[t.severity]??"#64748b",r=N(t.component_path),u=t.end_line&&t.end_line!==t.line?`L${t.line}\u2013${t.end_line}`:`L${t.line}`,H=q[t.type]??t.type;return`<div class="issue-row" data-idx="${a}">
      <span class="issue-sev">
        <span class="issue-sev-dot" style="background:${o}"></span>
        ${n(t.severity)}
      </span>
      <span class="issue-type">${n(H)}</span>
      <div class="issue-main">
        <span class="issue-msg">${n(t.message)}</span>
        <span class="issue-file" title="${n(t.component_path)}">${n(r)}:${u}</span>
      </div>
      <span class="issue-rule">${n(t.rule_key)}</span>
    </div>`}).join(""),e.querySelectorAll(".issue-row").forEach(t=>{t.addEventListener("click",()=>{let a=Number.parseInt(t.dataset.idx,10);R(a)})})}function pe(){let e=new Map;for(let i of c){let t=i.component_path;e.has(t)||e.set(t,[]),e.get(t).push(i)}f=[...e.entries()].sort((i,t)=>t[1].length-i[1].length).map(([i,t])=>({path:i,shortPath:N(i),issues:[...t].sort((a,o)=>a.line-o.line),expanded:!1}))}function X(){let e=l("file-tree");if(!f.length){e.innerHTML='<div class="empty-state">No issues found</div>';return}e.innerHTML=f.map((i,t)=>`<div class="file-group${i.expanded?" expanded":""}" data-gi="${t}">
      <div class="file-group-header">
        <span class="file-group-chevron">\u25B6</span>
        <span class="file-group-name" title="${n(i.path)}">${n(i.shortPath)}</span>
        <span class="file-group-count">${i.issues.length}</span>
      </div>
      <div class="file-group-issues" style="${i.expanded?"":"display:none"}">
        ${i.issues.map((a,o)=>{let r=S[a.severity]??"#64748b";return`<div class="file-issue" data-gi="${t}" data-ii="${o}">
            <span class="issue-sev">
              <span class="issue-sev-dot" style="background:${r}"></span>
              ${n(a.severity)}
            </span>
            <span class="issue-msg">${n(a.message)}</span>
            <span class="file-issue-line">L${a.line}</span>
          </div>`}).join("")}
      </div>
    </div>`).join(""),e.querySelectorAll(".file-group-header").forEach(i=>{i.addEventListener("click",()=>{let t=i.closest(".file-group"),a=Number.parseInt(t.dataset.gi,10);f[a].expanded=!f[a].expanded,t.classList.toggle("expanded");let o=t.querySelector(".file-group-issues");o.style.display=f[a].expanded?"":"none"})}),e.querySelectorAll(".file-issue").forEach(i=>{i.addEventListener("click",t=>{t.stopPropagation();let a=Number.parseInt(i.dataset.gi,10),o=Number.parseInt(i.dataset.ii,10),r=f[a].issues[o];O(r)})})}function ve(e){let i=f.findIndex(a=>a.path===e);if(i<0)return;f[i].expanded=!0,X(),document.querySelector(`.file-group[data-gi="${i}"]`)?.scrollIntoView({behavior:"smooth",block:"start"})}function R(e){b=e,p=g[e]??null,document.querySelectorAll(".issue-row").forEach(i=>i.classList.remove("selected")),document.querySelector(`.issue-row[data-idx="${e}"]`)?.classList.add("selected"),p&&O(p)}function O(e){p=e,$="details",I="",k=!0,s=G(),l("detail-title").textContent=e.rule_key,_(e),l("detail-panel").classList.add("open"),l("detail-overlay").classList.add("open"),ue(e.rule_key)}async function ue(e){try{let i=await fetch(`/rules/${encodeURIComponent(e)}`);if(!i.ok)throw new Error("not found");let t=await i.json(),a="";t.rationale&&(a+=`<div class="detail-section">
        <div class="detail-section-title">Why is this a problem?</div>
        <div class="rule-rationale">${n(t.rationale)}</div>
      </div>`),t.description&&t.description!==t.rationale&&(a+=`<div class="detail-section">
        <div class="detail-section-title">Description</div>
        <div class="rule-rationale">${n(t.description)}</div>
      </div>`),t.noncompliant_code&&(a+=`<div class="detail-section">
        <div class="detail-section-title">\u2718 Noncompliant Code</div>
        <pre class="rule-code noncompliant"><code>${n(t.noncompliant_code)}</code></pre>
      </div>`),t.compliant_code&&(a+=`<div class="detail-section">
        <div class="detail-section-title">\u2714 Compliant Code</div>
        <pre class="rule-code compliant"><code>${n(t.compliant_code)}</code></pre>
      </div>`),I=a||'<div class="detail-empty">No additional rule details available.</div>'}catch{I='<div class="detail-empty">Rule details are not available for this issue.</div>'}finally{k=!1,p?.rule_key===e&&_(p)}}function D(){l("detail-panel").classList.remove("open"),l("detail-overlay").classList.remove("open"),p=null,I="",k=!1,s=G(),document.querySelectorAll(".issue-row").forEach(e=>e.classList.remove("selected"))}function _(e){let i=`
    <div class="detail-tabs">
      ${V($)}
    </div>
    <div class="detail-tab-panel${$==="details"?"":" hidden"}" data-detail-panel="details">
      ${fe(e)}
    </div>
    <div class="detail-tab-panel${$==="rule"?"":" hidden"}" data-detail-panel="rule">
      ${k?'<div class="detail-loading">Loading rule details\u2026</div>':I}
    </div>
    <div class="detail-tab-panel${$==="ai-fix"?"":" hidden"}" data-detail-panel="ai-fix">
      ${Y(e,s,v??[])}
    </div>
  `;l("detail-body").innerHTML=i,ye(e)}function fe(e){let i=S[e.severity]??"#64748b",t=q[e.type]??e.type,a=e.end_line&&e.end_line!==e.line?`${e.line}:${e.column} \u2013 ${e.end_line}:${e.end_column}`:`${e.line}:${e.column}`,o=`
    <div class="detail-section">
      <div class="detail-msg">${n(e.message)}</div>
    </div>
    <div class="detail-section">
      <div class="detail-section-title">Properties</div>
      <div class="detail-field">
        <span class="detail-field-label">Severity</span>
        <span class="detail-field-value"><span class="issue-sev-dot" style="background:${i};display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:6px"></span>${n(e.severity)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Type</span>
        <span class="detail-field-value">${n(t)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Rule</span>
        <span class="detail-field-value" style="font-family:var(--font-mono);color:var(--accent)">${n(e.rule_key)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Status</span>
        <span class="detail-field-value">${n(e.status)}</span>
      </div>
      ${e.engine_id?`<div class="detail-field">
        <span class="detail-field-label">Engine</span>
        <span class="detail-field-value">${n(e.engine_id)}</span>
      </div>`:""}
      ${e.tags?.length?`<div class="detail-field">
        <span class="detail-field-label">Tags</span>
        <span class="detail-field-value">${e.tags.map(r=>n(r)).join(", ")}</span>
      </div>`:""}
    </div>
    <div class="detail-section">
      <div class="detail-section-title">Location</div>
      <div class="detail-field">
        <span class="detail-field-label">File</span>
        <span class="detail-field-value" style="font-family:var(--font-mono);font-size:12px;word-break:break-all">${n(e.component_path)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Lines</span>
        <span class="detail-field-value" style="font-family:var(--font-mono)">${a}</span>
      </div>
    </div>`;return e.secondary_locations?.length&&(o+=`<div class="detail-section">
      <div class="detail-section-title">Related Locations (${e.secondary_locations.length})</div>
      <div class="detail-loc-list">
        ${e.secondary_locations.map(r=>`
          <div class="detail-loc-item">
            <div class="detail-loc-file">${n(r.file_path||e.component_path)}:${r.start_line}</div>
            ${r.message?`<div class="detail-loc-msg">${n(r.message)}</div>`:""}
          </div>
        `).join("")}
      </div>
    </div>`),o}function Y(e){let i=e.end_line&&e.end_line!==e.line?`-${e.end_line}`:"",t=(v?.length??0)>0,a=(v??[]).map(u=>`<option value="${n(u.id)}"${s.selectedAgentId===u.id?" selected":""}>${n(u.label)} \xB7 ${n(u.model)}</option>`).join(""),o=ge(t,a),r=me();return`
    <div class="detail-section">
      <div class="detail-section-title">Fix with AI</div>
      <div class="detail-msg ai-fix-callout">Ollanta prepara o contexto da issue, envia apenas o trecho relevante para o agente escolhido e mostra um preview antes de qualquer escrita no seu c\xF3digo.</div>
    </div>

    <div class="detail-section">
      <div class="detail-field detail-field-stack">
        <span class="detail-field-label">Target</span>
        <span class="detail-field-value detail-mono-block">${n(e.component_path)}:${e.line}${i}</span>
      </div>
      <div class="detail-field detail-field-stack">
        <span class="detail-field-label">Issue</span>
        <span class="detail-field-value">${n(e.message)}</span>
      </div>
    </div>

    <div class="detail-section">
      <div class="detail-section-title">Agent</div>
      ${o}
      ${s.statusMessage?`<div class="ai-fix-status ai-fix-status-ok">${n(s.statusMessage)}</div>`:""}
      ${s.errorMessage?`<div class="ai-fix-status ai-fix-status-error">${n(s.errorMessage)}</div>`:""}
    </div>

    <div class="detail-section">
      <div class="detail-section-title">Preview</div>
      ${r}
    </div>
  `}function ge(e,i){if(s.loadingAgents)return'<div class="detail-loading">Loading AI agents\u2026</div>';if(!e)return'<div class="detail-empty">No AI agent is configured for the local scanner.</div>';let t=s.loadingPreview?"Generating\u2026":"Generate fix",a=s.loadingPreview?" disabled":"";return`<div class="ai-fix-controls">
      <select id="ai-agent-select" class="ai-fix-select">${i}</select>
      <button id="ai-generate-fix" class="ai-fix-button"${a}>${t}</button>
    </div>`}function me(){if(!s.preview)return'<div class="detail-empty">Generate a fix preview to inspect the patch before Ollanta edits your local file.</div>';let e=s.preview.summary||"Generated fix preview",i=s.preview.explanation?`<div class="rule-rationale">${n(s.preview.explanation)}</div>`:"",t=s.applying?"Applying\u2026":"Apply to file",a=s.applying?" disabled":"";return`
    <div class="ai-fix-preview-meta">
      <div><strong>Agent:</strong> ${n(s.preview.agent.label)}</div>
      <div><strong>Summary:</strong> ${n(e)}</div>
    </div>
    ${i}
    <pre class="rule-code ai-fix-diff"><code>${n(s.preview.diff)}</code></pre>
    <div class="ai-fix-actions">
      <button id="ai-apply-fix" class="ai-fix-button ai-fix-button-primary"${a}>${t}</button>
    </div>
  `}function ye(e){document.querySelectorAll(".detail-tab").forEach(t=>{t.addEventListener("click",()=>{$=t.dataset.detailTab??"details",_(e),$==="ai-fix"&&be()})});let i=document.getElementById("ai-agent-select");i?.addEventListener("change",()=>{s.selectedAgentId=i.value}),document.getElementById("ai-generate-fix")?.addEventListener("click",()=>{$e(e)}),document.getElementById("ai-apply-fix")?.addEventListener("click",()=>{he()})}function G(){return{loadingAgents:!1,loadingPreview:!1,applying:!1,selectedAgentId:"",statusMessage:"",errorMessage:"",preview:null}}async function be(){if(v){!s.selectedAgentId&&v.length>0&&(s.selectedAgentId=v[0].id,m());return}s.loadingAgents=!0,s.errorMessage="",m();try{let e=await fetch("/api/ai/agents");if(!e.ok)throw new Error(`HTTP ${e.status}`);v=(await e.json()).agents??[],!s.selectedAgentId&&v.length>0&&(s.selectedAgentId=v[0].id)}catch(e){s.errorMessage=`Failed to load AI agents: ${String(e)}`,v=[]}finally{s.loadingAgents=!1,m()}}async function $e(e){if(!s.selectedAgentId){s.errorMessage="Choose an AI agent before generating a fix.",m();return}s.loadingPreview=!0,s.statusMessage="",s.errorMessage="",m();try{let i=await fetch("/api/ai/fixes/preview",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({agent_id:s.selectedAgentId,issue:e})}),t=await i.json();if(!i.ok||"error"in t)throw new Error("error"in t?t.error:`HTTP ${i.status}`);s.preview=t,s.statusMessage="Fix preview generated. Review the diff before applying it."}catch(i){s.errorMessage=`Failed to generate AI fix: ${String(i)}`,s.preview=null}finally{s.loadingPreview=!1,m()}}async function he(){if(s.preview){s.applying=!0,s.errorMessage="",m();try{let e=await fetch("/api/ai/fixes/apply",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({preview_id:s.preview.preview_id})}),i=await e.json();if(!e.ok||"error"in i)throw new Error("error"in i?i.error:`HTTP ${e.status}`);s.statusMessage=i.message}catch(e){s.errorMessage=`Failed to apply AI fix: ${String(e)}`}finally{s.applying=!1,m()}}}function m(){p&&_(p)}document.addEventListener("DOMContentLoaded",()=>{l("detail-close").addEventListener("click",D),l("detail-overlay").addEventListener("click",D)});function xe(){document.addEventListener("keydown",e=>{let i=e.target.tagName;if(!(i==="INPUT"||i==="SELECT"||i==="TEXTAREA")){if(e.key==="Escape"){D();return}U==="issues"&&(e.key==="j"||e.key==="ArrowDown"?(e.preventDefault(),b<g.length-1&&R(b+1),K()):e.key==="k"||e.key==="ArrowUp"?(e.preventDefault(),b>0&&R(b-1),K()):e.key==="Enter"&&p&&O(p))}})}function K(){document.querySelector(`.issue-row[data-idx="${b}"]`)?.scrollIntoView({behavior:"smooth",block:"nearest"})}function l(e){return document.getElementById(e)}function A(e,i){l(e).classList.add(i)}function M(e,i){let t={};for(let a of e){let o=i(a);t[o]=(t[o]??0)+1}return t}function N(e){let i=e.replaceAll("\\","/"),t=i.split("/").filter(Boolean);return t.length<=2?i:`${t.slice(-2).join("/")}`}})();
