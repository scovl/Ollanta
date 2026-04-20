"use strict";(()=>{var b,r=[],p=[],d=[],m=null,u=-1,O="overview",f="all",M="all",S="all",H="",q={blocker:0,critical:1,major:2,minor:3,info:4},L={blocker:"#ef4444",critical:"#f97316",major:"#eab308",minor:"#22c55e",info:"#64748b"},C={bug:"Bug",code_smell:"Code Smell",vulnerability:"Vulnerability",security_hotspot:"Hotspot"};async function V(){try{let e=await fetch("/report.json");if(!e.ok)throw new Error(`HTTP ${e.status}`);b=await e.json(),r=b.issues??[],U(),K(),W(),X(),Y(),J(),Z(),y(),te(),G(),Q(),ae(),a("tab-issue-count").textContent=String(r.length),a("tab-file-count").textContent=String(new Set(r.map(t=>t.component_path)).size)}catch(e){a("app").innerHTML=`<div class="error">Failed to load report: ${String(e)}</div>`}}document.addEventListener("DOMContentLoaded",V);function U(){let e=b.metadata,t=new Date(e.analysis_date).toLocaleString();a("project-key").textContent=e.project_key,a("scan-date").textContent=t,a("scan-version").textContent=`v${e.version}`,a("elapsed").textContent=`${e.elapsed_ms}ms`}function z(){let e=b.measures,t=[{metric:"Bugs",operator:"=",threshold:0,value:e.bugs,passed:e.bugs===0},{metric:"Vulnerabilities",operator:"=",threshold:0,value:e.vulnerabilities,passed:e.vulnerabilities===0}];return{status:t.every(n=>n.passed)?"passed":"failed",conditions:t}}function K(){let e=z(),t=a("gate-hero");t.classList.remove("gate-loading"),t.classList.add(e.status==="passed"?"gate-passed":"gate-failed"),a("gate-icon").textContent=e.status==="passed"?"\u2713":"\u2717",a("gate-status").textContent=e.status==="passed"?"Passed":"Failed";let s=e.conditions.map(n=>{let i=n.passed?"cond-pass":"cond-fail",o=n.passed?"\u2713":"\u2717";return`<div class="gate-cond ${i}">
      <span class="gate-cond-icon">${o}</span>
      <span class="gate-cond-metric">${l(n.metric)}</span>
      <span class="gate-cond-value">${n.value}</span>
    </div>`}).join("");a("gate-conditions").innerHTML=s}function W(){let e=b.measures;g("m-bugs",e.bugs),g("m-vulns",e.vulnerabilities),g("m-smells",e.code_smells),g("m-ncloc",e.ncloc),g("m-files",e.files),g("m-comments",e.comments),_("card-bugs",e.bugs,[0,1,5]),_("card-vulns",e.vulnerabilities,[0,1,3]),_("card-smells",e.code_smells,[0,10,50]),h("card-ncloc","card-neutral"),h("card-files","card-neutral"),h("card-comments","card-neutral")}function g(e,t){a(e).textContent=t.toLocaleString()}function _(e,t,s){let n=a(e);t<=s[0]?h(e,"card-green"):t<=s[1]?h(e,"card-yellow"):h(e,"card-red")}function X(){let e=k(r,c=>c.severity),t=Math.max(1,...Object.values(e)),s=["blocker","critical","major","minor","info"],n="",i="",o=r.length||1;for(let c of s){let v=e[c]??0,E=v/t*100,T=L[c]??"#64748b";n+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${c}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${E}%;background:${T}"></div></div>
      <span class="sev-bar-count">${v}</span>
    </div>`,v>0&&(i+=`<div class="sev-segment" style="width:${v/o*100}%;background:${T}" title="${c}: ${v}"></div>`)}a("sev-bars").innerHTML=n,a("sev-proportional").innerHTML=i;let $=k(r,c=>c.type),P=Math.max(1,...Object.values($)),B={bug:"#ef4444",vulnerability:"#f97316",code_smell:"#22c55e",security_hotspot:"#eab308"},R="";for(let[c,v]of Object.entries(C)){let E=$[c]??0,T=E/P*100,N=B[c]??"#64748b";R+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${v}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${T}%;background:${N}"></div></div>
      <span class="sev-bar-count">${E}</span>
    </div>`}a("type-bars").innerHTML=R}function Y(){let e=k(r,s=>s.component_path),t=Object.entries(e).sort((s,n)=>n[1]-s[1]).slice(0,10);if(!t.length){a("hotspot-files").innerHTML='<div class="empty-state">No issues found</div>';return}a("hotspot-files").innerHTML=t.map(([s,n])=>{let i=j(s);return`<div class="hotspot-row" data-path="${l(s)}">
      <span class="hotspot-file" title="${l(s)}">${l(i)}</span>
      <span class="hotspot-count">${n}</span>
    </div>`}).join(""),a("hotspot-files").querySelectorAll(".hotspot-row").forEach(s=>{s.addEventListener("click",()=>{let n=s.dataset.path;D("files"),se(n)})})}function J(){let e=Object.entries(b.measures.by_language).sort((s,n)=>n[1]-s[1]),t=Math.max(1,e[0]?.[1]??1);if(!e.length){a("by-lang").innerHTML='<span class="empty-state">No language data</span>';return}a("by-lang").innerHTML=e.map(([s,n])=>`<div class="lang-row">
      <span class="lang-name">${l(s)}</span>
      <div class="lang-bar-track"><div class="lang-bar-fill" style="width:${n/t*100}%"></div></div>
      <span class="lang-count">${n} files</span>
    </div>`).join("")}function Q(){document.querySelectorAll(".tab").forEach(e=>{e.addEventListener("click",()=>{let t=e.dataset.tab;D(t)})})}function D(e){O=e,document.querySelectorAll(".tab").forEach(t=>t.classList.remove("active")),document.querySelector(`.tab[data-tab="${e}"]`)?.classList.add("active"),document.querySelectorAll(".panel").forEach(t=>t.classList.add("hidden")),a(`panel-${e}`).classList.remove("hidden")}function Z(){let e=[...new Set(r.map(s=>s.rule_key))].sort(),t=a("filter-rule");e.forEach(s=>{let n=document.createElement("option");n.value=s,n.textContent=s,t.appendChild(n)}),a("filter-severity").addEventListener("change",s=>{f=s.target.value,y()}),a("filter-type").addEventListener("change",s=>{M=s.target.value,y()}),t.addEventListener("change",s=>{S=s.target.value,y()}),a("search").addEventListener("input",s=>{H=s.target.value.toLowerCase(),y()}),F()}function F(){let e=k(r,s=>s.severity),t=["blocker","critical","major","minor","info"];a("sev-chips").innerHTML=t.map(s=>{let n=e[s]??0,i=L[s];return`<div class="sev-chip${f===s?" active":""}" data-sev="${s}"
      style="--chip-color:${i};--chip-bg:${i}15">
      <span class="chip-dot" style="background:${i}"></span>
      ${s}
      <span class="chip-count">${n}</span>
    </div>`}).join(""),a("sev-chips").querySelectorAll(".sev-chip").forEach(s=>{s.addEventListener("click",()=>{let n=s.dataset.sev;f=f===n?"all":n,a("filter-severity").value=f,y(),F()})})}function y(){p=r.filter(e=>!(f!=="all"&&e.severity!==f||M!=="all"&&e.type!==M||S!=="all"&&e.rule_key!==S||H&&!`${e.component_path} ${e.message} ${e.rule_key}`.toLowerCase().includes(H))),p.sort((e,t)=>{let s=q[e.severity]??99,n=q[t.severity]??99;return s-n}),u=-1,ee()}function ee(){let e=a("issue-list");if(a("issue-count").textContent=`${p.length} issue${p.length!==1?"s":""}`,!p.length){e.innerHTML='<div class="empty-state">No issues match the current filters.</div>';return}e.innerHTML=p.map((t,s)=>{let n=L[t.severity]??"#64748b",i=j(t.component_path),o=t.end_line&&t.end_line!==t.line?`L${t.line}\u2013${t.end_line}`:`L${t.line}`,$=C[t.type]??t.type;return`<div class="issue-row" data-idx="${s}">
      <span class="issue-sev">
        <span class="issue-sev-dot" style="background:${n}"></span>
        ${l(t.severity)}
      </span>
      <span class="issue-type">${l($)}</span>
      <div class="issue-main">
        <span class="issue-msg">${l(t.message)}</span>
        <span class="issue-file" title="${l(t.component_path)}">${l(i)}:${o}</span>
      </div>
      <span class="issue-rule">${l(t.rule_key)}</span>
    </div>`}).join(""),e.querySelectorAll(".issue-row").forEach(t=>{t.addEventListener("click",()=>{let s=parseInt(t.dataset.idx,10);w(s)})})}function te(){let e=new Map;for(let t of r){let s=t.component_path;e.has(s)||e.set(s,[]),e.get(s).push(t)}d=[...e.entries()].sort((t,s)=>s[1].length-t[1].length).map(([t,s])=>({path:t,shortPath:j(t),issues:s.sort((n,i)=>n.line-i.line),expanded:!1}))}function G(){let e=a("file-tree");if(!d.length){e.innerHTML='<div class="empty-state">No issues found</div>';return}e.innerHTML=d.map((t,s)=>`<div class="file-group${t.expanded?" expanded":""}" data-gi="${s}">
      <div class="file-group-header">
        <span class="file-group-chevron">\u25B6</span>
        <span class="file-group-name" title="${l(t.path)}">${l(t.shortPath)}</span>
        <span class="file-group-count">${t.issues.length}</span>
      </div>
      <div class="file-group-issues" style="${t.expanded?"":"display:none"}">
        ${t.issues.map((n,i)=>{let o=L[n.severity]??"#64748b";return`<div class="file-issue" data-gi="${s}" data-ii="${i}">
            <span class="issue-sev">
              <span class="issue-sev-dot" style="background:${o}"></span>
              ${l(n.severity)}
            </span>
            <span class="issue-msg">${l(n.message)}</span>
            <span class="file-issue-line">L${n.line}</span>
          </div>`}).join("")}
      </div>
    </div>`).join(""),e.querySelectorAll(".file-group-header").forEach(t=>{t.addEventListener("click",()=>{let s=t.closest(".file-group"),n=parseInt(s.dataset.gi,10);d[n].expanded=!d[n].expanded,s.classList.toggle("expanded");let i=s.querySelector(".file-group-issues");i.style.display=d[n].expanded?"":"none"})}),e.querySelectorAll(".file-issue").forEach(t=>{t.addEventListener("click",s=>{s.stopPropagation();let n=parseInt(t.dataset.gi,10),i=parseInt(t.dataset.ii,10),o=d[n].issues[i];I(o)})})}function se(e){let t=d.findIndex(n=>n.path===e);if(t<0)return;d[t].expanded=!0,G(),document.querySelector(`.file-group[data-gi="${t}"]`)?.scrollIntoView({behavior:"smooth",block:"start"})}function w(e){u=e,m=p[e]??null,document.querySelectorAll(".issue-row").forEach(t=>t.classList.remove("selected")),document.querySelector(`.issue-row[data-idx="${e}"]`)?.classList.add("selected"),m&&I(m)}function I(e){m=e;let t=L[e.severity]??"#64748b",s=C[e.type]??e.type,n=e.end_line&&e.end_line!==e.line?`${e.line}:${e.column} \u2013 ${e.end_line}:${e.end_column}`:`${e.line}:${e.column}`,i=`
    <div class="detail-section">
      <div class="detail-msg">${l(e.message)}</div>
    </div>
    <div class="detail-section">
      <div class="detail-section-title">Properties</div>
      <div class="detail-field">
        <span class="detail-field-label">Severity</span>
        <span class="detail-field-value"><span class="issue-sev-dot" style="background:${t};display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:6px"></span>${l(e.severity)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Type</span>
        <span class="detail-field-value">${l(s)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Rule</span>
        <span class="detail-field-value" style="font-family:var(--font-mono);color:var(--accent)">${l(e.rule_key)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Status</span>
        <span class="detail-field-value">${l(e.status)}</span>
      </div>
      ${e.engine_id?`<div class="detail-field">
        <span class="detail-field-label">Engine</span>
        <span class="detail-field-value">${l(e.engine_id)}</span>
      </div>`:""}
      ${e.tags?.length?`<div class="detail-field">
        <span class="detail-field-label">Tags</span>
        <span class="detail-field-value">${e.tags.map(o=>l(o)).join(", ")}</span>
      </div>`:""}
    </div>
    <div class="detail-section">
      <div class="detail-section-title">Location</div>
      <div class="detail-field">
        <span class="detail-field-label">File</span>
        <span class="detail-field-value" style="font-family:var(--font-mono);font-size:12px;word-break:break-all">${l(e.component_path)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Lines</span>
        <span class="detail-field-value" style="font-family:var(--font-mono)">${n}</span>
      </div>
    </div>`;e.secondary_locations?.length&&(i+=`<div class="detail-section">
      <div class="detail-section-title">Related Locations (${e.secondary_locations.length})</div>
      <div class="detail-loc-list">
        ${e.secondary_locations.map((o,$)=>`
          <div class="detail-loc-item">
            <div class="detail-loc-file">${l(o.file_path||e.component_path)}:${o.start_line}</div>
            ${o.message?`<div class="detail-loc-msg">${l(o.message)}</div>`:""}
          </div>
        `).join("")}
      </div>
    </div>`),a("detail-title").textContent=e.rule_key,a("detail-body").innerHTML=i,a("detail-panel").classList.add("open"),a("detail-overlay").classList.add("open"),ne(e.rule_key)}async function ne(e){if(!document.getElementById("rule-detail-section")){let s=a("detail-body"),n=document.createElement("div");n.id="rule-detail-section",n.innerHTML='<div class="detail-section"><div class="detail-section-title">Loading rule details\u2026</div></div>',s.appendChild(n)}try{let s=await fetch(`/rules/${encodeURIComponent(e)}`);if(!s.ok)throw new Error("not found");let n=await s.json(),i=document.getElementById("rule-detail-section");if(!i)return;let o="";n.rationale&&(o+=`<div class="detail-section">
        <div class="detail-section-title">Why is this a problem?</div>
        <div class="rule-rationale">${l(n.rationale)}</div>
      </div>`),n.description&&n.description!==n.rationale&&(o+=`<div class="detail-section">
        <div class="detail-section-title">Description</div>
        <div class="rule-rationale">${l(n.description)}</div>
      </div>`),n.noncompliant_code&&(o+=`<div class="detail-section">
        <div class="detail-section-title">\u2718 Noncompliant Code</div>
        <pre class="rule-code noncompliant"><code>${l(n.noncompliant_code)}</code></pre>
      </div>`),n.compliant_code&&(o+=`<div class="detail-section">
        <div class="detail-section-title">\u2714 Compliant Code</div>
        <pre class="rule-code compliant"><code>${l(n.compliant_code)}</code></pre>
      </div>`),i.innerHTML=o}catch{let s=document.getElementById("rule-detail-section");s&&(s.innerHTML="")}}function x(){a("detail-panel").classList.remove("open"),a("detail-overlay").classList.remove("open"),m=null,document.querySelectorAll(".issue-row").forEach(e=>e.classList.remove("selected"))}document.addEventListener("DOMContentLoaded",()=>{a("detail-close").addEventListener("click",x),a("detail-overlay").addEventListener("click",x)});function ae(){document.addEventListener("keydown",e=>{let t=e.target.tagName;if(!(t==="INPUT"||t==="SELECT"||t==="TEXTAREA")){if(e.key==="Escape"){x();return}O==="issues"&&(e.key==="j"||e.key==="ArrowDown"?(e.preventDefault(),u<p.length-1&&w(u+1),A()):e.key==="k"||e.key==="ArrowUp"?(e.preventDefault(),u>0&&w(u-1),A()):e.key==="Enter"&&m&&I(m))}})}function A(){document.querySelector(`.issue-row[data-idx="${u}"]`)?.scrollIntoView({behavior:"smooth",block:"nearest"})}function a(e){return document.getElementById(e)}function h(e,t){a(e).classList.add(t)}function l(e){return e.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;")}function j(e){let t=e.replace(/\\/g,"/").split("/");return t.length>3?t.slice(-3).join("/"):t.join("/")}function k(e,t){let s={};for(let n of e){let i=t(n);s[i]=(s[i]??0)+1}return s}})();
