"use strict";(()=>{var b,c=[],p=[],d=[],m=null,v=-1,O="overview",f="all",M="all",_="all",H="",q={blocker:0,critical:1,major:2,minor:3,info:4},L={blocker:"#ef4444",critical:"#f97316",major:"#eab308",minor:"#22c55e",info:"#64748b"},C={bug:"Bug",code_smell:"Code Smell",vulnerability:"Vulnerability",security_hotspot:"Hotspot"};async function B(){try{let e=await fetch("/report.json");if(!e.ok)throw new Error(`HTTP ${e.status}`);b=await e.json(),c=b.issues??[],z(),U(),X(),Y(),J(),Q(),Z(),y(),te(),P(),W(),ne(),a("tab-issue-count").textContent=String(c.length),a("tab-file-count").textContent=String(new Set(c.map(t=>t.component_path)).size)}catch(e){a("app").innerHTML=`<div class="error">Failed to load report: ${String(e)}</div>`}}document.addEventListener("DOMContentLoaded",B);function z(){let e=b.metadata,t=new Date(e.analysis_date).toLocaleString();a("project-key").textContent=e.project_key,a("scan-date").textContent=t,a("scan-version").textContent=`v${e.version}`,a("elapsed").textContent=`${e.elapsed_ms}ms`}function K(){let e=b.measures,t=[{metric:"Bugs",operator:"=",threshold:0,value:e.bugs,passed:e.bugs===0},{metric:"Vulnerabilities",operator:"=",threshold:0,value:e.vulnerabilities,passed:e.vulnerabilities===0}];return{status:t.every(n=>n.passed)?"passed":"failed",conditions:t}}function U(){let e=K(),t=a("gate-hero");t.classList.remove("gate-loading"),t.classList.add(e.status==="passed"?"gate-passed":"gate-failed"),a("gate-icon").textContent=e.status==="passed"?"\u2713":"\u2717",a("gate-status").textContent=e.status==="passed"?"Passed":"Failed";let s=e.conditions.map(n=>{let l=n.passed?"cond-pass":"cond-fail",o=n.passed?"\u2713":"\u2717";return`<div class="gate-cond ${l}">
      <span class="gate-cond-icon">${o}</span>
      <span class="gate-cond-metric">${i(n.metric)}</span>
      <span class="gate-cond-value">${n.value}</span>
    </div>`}).join("");a("gate-conditions").innerHTML=s}function X(){let e=b.measures;g("m-bugs",e.bugs),g("m-vulns",e.vulnerabilities),g("m-smells",e.code_smells),g("m-ncloc",e.ncloc),g("m-files",e.files),g("m-comments",e.comments),S("card-bugs",e.bugs,[0,1,5]),S("card-vulns",e.vulnerabilities,[0,1,3]),S("card-smells",e.code_smells,[0,10,50]),h("card-ncloc","card-neutral"),h("card-files","card-neutral"),h("card-comments","card-neutral")}function g(e,t){a(e).textContent=t.toLocaleString()}function S(e,t,s){let n=a(e);t<=s[0]?h(e,"card-green"):t<=s[1]?h(e,"card-yellow"):h(e,"card-red")}function Y(){let e=k(c,r=>r.severity),t=Math.max(1,...Object.values(e)),s=["blocker","critical","major","minor","info"],n="",l="",o=c.length||1;for(let r of s){let u=e[r]??0,E=u/t*100,T=L[r]??"#64748b";n+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${r}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${E}%;background:${T}"></div></div>
      <span class="sev-bar-count">${u}</span>
    </div>`,u>0&&(l+=`<div class="sev-segment" style="width:${u/o*100}%;background:${T}" title="${r}: ${u}"></div>`)}a("sev-bars").innerHTML=n,a("sev-proportional").innerHTML=l;let $=k(c,r=>r.type),D=Math.max(1,...Object.values($)),N={bug:"#ef4444",vulnerability:"#f97316",code_smell:"#22c55e",security_hotspot:"#eab308"},R="";for(let[r,u]of Object.entries(C)){let E=$[r]??0,T=E/D*100,V=N[r]??"#64748b";R+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${u}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${T}%;background:${V}"></div></div>
      <span class="sev-bar-count">${E}</span>
    </div>`}a("type-bars").innerHTML=R}function J(){let e=k(c,s=>s.component_path),t=Object.entries(e).sort((s,n)=>n[1]-s[1]).slice(0,10);if(!t.length){a("hotspot-files").innerHTML='<div class="empty-state">No issues found</div>';return}a("hotspot-files").innerHTML=t.map(([s,n])=>{let l=I(s);return`<div class="hotspot-row" data-path="${i(s)}">
      <span class="hotspot-file" title="${i(s)}">${i(l)}</span>
      <span class="hotspot-count">${n}</span>
    </div>`}).join(""),a("hotspot-files").querySelectorAll(".hotspot-row").forEach(s=>{s.addEventListener("click",()=>{let n=s.dataset.path;F("files"),se(n)})})}function Q(){let e=Object.entries(b.measures.by_language).sort((s,n)=>n[1]-s[1]),t=Math.max(1,e[0]?.[1]??1);if(!e.length){a("by-lang").innerHTML='<span class="empty-state">No language data</span>';return}a("by-lang").innerHTML=e.map(([s,n])=>`<div class="lang-row">
      <span class="lang-name">${i(s)}</span>
      <div class="lang-bar-track"><div class="lang-bar-fill" style="width:${n/t*100}%"></div></div>
      <span class="lang-count">${n} files</span>
    </div>`).join("")}function W(){document.querySelectorAll(".tab").forEach(e=>{e.addEventListener("click",()=>{let t=e.dataset.tab;F(t)})})}function F(e){O=e,document.querySelectorAll(".tab").forEach(t=>t.classList.remove("active")),document.querySelector(`.tab[data-tab="${e}"]`)?.classList.add("active"),document.querySelectorAll(".panel").forEach(t=>t.classList.add("hidden")),a(`panel-${e}`).classList.remove("hidden")}function Z(){let e=[...new Set(c.map(s=>s.rule_key))].sort(),t=a("filter-rule");e.forEach(s=>{let n=document.createElement("option");n.value=s,n.textContent=s,t.appendChild(n)}),a("filter-severity").addEventListener("change",s=>{f=s.target.value,y()}),a("filter-type").addEventListener("change",s=>{M=s.target.value,y()}),t.addEventListener("change",s=>{_=s.target.value,y()}),a("search").addEventListener("input",s=>{H=s.target.value.toLowerCase(),y()}),G()}function G(){let e=k(c,s=>s.severity),t=["blocker","critical","major","minor","info"];a("sev-chips").innerHTML=t.map(s=>{let n=e[s]??0,l=L[s];return`<div class="sev-chip${f===s?" active":""}" data-sev="${s}"
      style="--chip-color:${l};--chip-bg:${l}15">
      <span class="chip-dot" style="background:${l}"></span>
      ${s}
      <span class="chip-count">${n}</span>
    </div>`}).join(""),a("sev-chips").querySelectorAll(".sev-chip").forEach(s=>{s.addEventListener("click",()=>{let n=s.dataset.sev;f=f===n?"all":n,a("filter-severity").value=f,y(),G()})})}function y(){p=c.filter(e=>!(f!=="all"&&e.severity!==f||M!=="all"&&e.type!==M||_!=="all"&&e.rule_key!==_||H&&!`${e.component_path} ${e.message} ${e.rule_key}`.toLowerCase().includes(H))),p.sort((e,t)=>{let s=q[e.severity]??99,n=q[t.severity]??99;return s-n}),v=-1,ee()}function ee(){let e=a("issue-list");if(a("issue-count").textContent=`${p.length} issue${p.length!==1?"s":""}`,!p.length){e.innerHTML='<div class="empty-state">No issues match the current filters.</div>';return}e.innerHTML=p.map((t,s)=>{let n=L[t.severity]??"#64748b",l=I(t.component_path),o=t.end_line&&t.end_line!==t.line?`L${t.line}\u2013${t.end_line}`:`L${t.line}`,$=C[t.type]??t.type;return`<div class="issue-row" data-idx="${s}">
      <span class="issue-sev">
        <span class="issue-sev-dot" style="background:${n}"></span>
        ${i(t.severity)}
      </span>
      <span class="issue-type">${i($)}</span>
      <div class="issue-main">
        <span class="issue-msg">${i(t.message)}</span>
        <span class="issue-file" title="${i(t.component_path)}">${i(l)}:${o}</span>
      </div>
      <span class="issue-rule">${i(t.rule_key)}</span>
    </div>`}).join(""),e.querySelectorAll(".issue-row").forEach(t=>{t.addEventListener("click",()=>{let s=parseInt(t.dataset.idx,10);x(s)})})}function te(){let e=new Map;for(let t of c){let s=t.component_path;e.has(s)||e.set(s,[]),e.get(s).push(t)}d=[...e.entries()].sort((t,s)=>s[1].length-t[1].length).map(([t,s])=>({path:t,shortPath:I(t),issues:s.sort((n,l)=>n.line-l.line),expanded:!1}))}function P(){let e=a("file-tree");if(!d.length){e.innerHTML='<div class="empty-state">No issues found</div>';return}e.innerHTML=d.map((t,s)=>`<div class="file-group${t.expanded?" expanded":""}" data-gi="${s}">
      <div class="file-group-header">
        <span class="file-group-chevron">\u25B6</span>
        <span class="file-group-name" title="${i(t.path)}">${i(t.shortPath)}</span>
        <span class="file-group-count">${t.issues.length}</span>
      </div>
      <div class="file-group-issues" style="${t.expanded?"":"display:none"}">
        ${t.issues.map((n,l)=>{let o=L[n.severity]??"#64748b";return`<div class="file-issue" data-gi="${s}" data-ii="${l}">
            <span class="issue-sev">
              <span class="issue-sev-dot" style="background:${o}"></span>
              ${i(n.severity)}
            </span>
            <span class="issue-msg">${i(n.message)}</span>
            <span class="file-issue-line">L${n.line}</span>
          </div>`}).join("")}
      </div>
    </div>`).join(""),e.querySelectorAll(".file-group-header").forEach(t=>{t.addEventListener("click",()=>{let s=t.closest(".file-group"),n=parseInt(s.dataset.gi,10);d[n].expanded=!d[n].expanded,s.classList.toggle("expanded");let l=s.querySelector(".file-group-issues");l.style.display=d[n].expanded?"":"none"})}),e.querySelectorAll(".file-issue").forEach(t=>{t.addEventListener("click",s=>{s.stopPropagation();let n=parseInt(t.dataset.gi,10),l=parseInt(t.dataset.ii,10),o=d[n].issues[l];j(o)})})}function se(e){let t=d.findIndex(n=>n.path===e);if(t<0)return;d[t].expanded=!0,P(),document.querySelector(`.file-group[data-gi="${t}"]`)?.scrollIntoView({behavior:"smooth",block:"start"})}function x(e){v=e,m=p[e]??null,document.querySelectorAll(".issue-row").forEach(t=>t.classList.remove("selected")),document.querySelector(`.issue-row[data-idx="${e}"]`)?.classList.add("selected"),m&&j(m)}function j(e){m=e;let t=L[e.severity]??"#64748b",s=C[e.type]??e.type,n=e.end_line&&e.end_line!==e.line?`${e.line}:${e.column} \u2013 ${e.end_line}:${e.end_column}`:`${e.line}:${e.column}`,l=`
    <div class="detail-section">
      <div class="detail-msg">${i(e.message)}</div>
    </div>
    <div class="detail-section">
      <div class="detail-section-title">Properties</div>
      <div class="detail-field">
        <span class="detail-field-label">Severity</span>
        <span class="detail-field-value"><span class="issue-sev-dot" style="background:${t};display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:6px"></span>${i(e.severity)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Type</span>
        <span class="detail-field-value">${i(s)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Rule</span>
        <span class="detail-field-value" style="font-family:var(--font-mono);color:var(--accent)">${i(e.rule_key)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Status</span>
        <span class="detail-field-value">${i(e.status)}</span>
      </div>
      ${e.engine_id?`<div class="detail-field">
        <span class="detail-field-label">Engine</span>
        <span class="detail-field-value">${i(e.engine_id)}</span>
      </div>`:""}
      ${e.tags?.length?`<div class="detail-field">
        <span class="detail-field-label">Tags</span>
        <span class="detail-field-value">${e.tags.map(o=>i(o)).join(", ")}</span>
      </div>`:""}
    </div>
    <div class="detail-section">
      <div class="detail-section-title">Location</div>
      <div class="detail-field">
        <span class="detail-field-label">File</span>
        <span class="detail-field-value" style="font-family:var(--font-mono);font-size:12px;word-break:break-all">${i(e.component_path)}</span>
      </div>
      <div class="detail-field">
        <span class="detail-field-label">Lines</span>
        <span class="detail-field-value" style="font-family:var(--font-mono)">${n}</span>
      </div>
    </div>`;e.secondary_locations?.length&&(l+=`<div class="detail-section">
      <div class="detail-section-title">Related Locations (${e.secondary_locations.length})</div>
      <div class="detail-loc-list">
        ${e.secondary_locations.map((o,$)=>`
          <div class="detail-loc-item">
            <div class="detail-loc-file">${i(o.file_path||e.component_path)}:${o.start_line}</div>
            ${o.message?`<div class="detail-loc-msg">${i(o.message)}</div>`:""}
          </div>
        `).join("")}
      </div>
    </div>`),a("detail-title").textContent=e.rule_key,a("detail-body").innerHTML=l,a("detail-panel").classList.add("open"),a("detail-overlay").classList.add("open")}function w(){a("detail-panel").classList.remove("open"),a("detail-overlay").classList.remove("open"),m=null,document.querySelectorAll(".issue-row").forEach(e=>e.classList.remove("selected"))}document.addEventListener("DOMContentLoaded",()=>{a("detail-close").addEventListener("click",w),a("detail-overlay").addEventListener("click",w)});function ne(){document.addEventListener("keydown",e=>{let t=e.target.tagName;if(!(t==="INPUT"||t==="SELECT"||t==="TEXTAREA")){if(e.key==="Escape"){w();return}O==="issues"&&(e.key==="j"||e.key==="ArrowDown"?(e.preventDefault(),v<p.length-1&&x(v+1),A()):e.key==="k"||e.key==="ArrowUp"?(e.preventDefault(),v>0&&x(v-1),A()):e.key==="Enter"&&m&&j(m))}})}function A(){document.querySelector(`.issue-row[data-idx="${v}"]`)?.scrollIntoView({behavior:"smooth",block:"nearest"})}function a(e){return document.getElementById(e)}function h(e,t){a(e).classList.add(t)}function i(e){return e.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/>/g,"&gt;").replace(/"/g,"&quot;")}function I(e){let t=e.replace(/\\/g,"/").split("/");return t.length>3?t.slice(-3).join("/"):t.join("/")}function k(e,t){let s={};for(let n of e){let l=t(n);s[l]=(s[l]??0)+1}return s}})();
