"use strict";(()=>{function i(e){return e.replaceAll("&","&amp;").replaceAll("<","&lt;").replaceAll(">","&gt;").replaceAll('"',"&quot;")}function he(e){return[{key:"details",label:"Details"},{key:"rule",label:"Rule"},{key:"ai-fix",label:"Fix with AI"}].map(s=>`<button class="detail-tab${e===s.key?" active":""}" data-detail-tab="${s.key}">${s.label}</button>`).join("")}function ye(e,t,s){let a=e.end_line&&e.end_line!==e.line?`-${e.end_line}`:"",n=Ie(t,s),l=Pe(t);return`
    <div class="detail-section">
      <div class="detail-section-title">Fix with AI</div>
      <div class="detail-msg ai-fix-callout">Ollanta prepares the issue context, sends only the relevant snippet to the selected agent, and shows a preview before writing any changes to your code.</div>
    </div>

    <div class="detail-section">
      <div class="detail-field detail-field-stack">
        <span class="detail-field-label">Target</span>
        <span class="detail-field-value detail-mono-block">${i(e.component_path)}:${e.line}${a}</span>
      </div>
      <div class="detail-field detail-field-stack">
        <span class="detail-field-label">Issue</span>
        <span class="detail-field-value">${i(e.message)}</span>
      </div>
    </div>

    <div class="detail-section">
      <div class="detail-section-title">Model</div>
      ${n}
      ${t.statusMessage?`<div class="ai-fix-status ai-fix-status-ok">${i(t.statusMessage)}</div>`:""}
      ${t.errorMessage?`<div class="ai-fix-status ai-fix-status-error">${i(t.errorMessage)}</div>`:""}
    </div>

    <div class="detail-section">
      <div class="detail-section-title">Preview</div>
      ${l}
    </div>
  `}function Ie(e,t){if(e.loadingOptions)return'<div class="detail-loading">Loading AI models\u2026</div>';if(t.length===0)return'<div class="detail-empty">No AI provider is available for the local scanner.</div>';let s=t.find(m=>m.id===e.selectedProviderId)??t[0],a=t.map(m=>`<option value="${i(m.id)}"${e.selectedProviderId===m.id?" selected":""}>${i(m.label)}</option>`).join(""),l=(s?.models??[]).map(m=>`<option value="${i(m)}"></option>`).join(""),u='<div class="ai-fix-helper">This provider can generate local previews without an API key.</div>',p="Required for this provider";s?.requires_api_key&&(s.configured?(u=`<div class="ai-fix-helper">Using the scanner's configured API key. Paste another key below to override it for this session.</div>`,p="Optional override"):u='<div class="ai-fix-helper">Paste an API key for the selected provider to generate the fix.</div>');let c=s?.requires_api_key?`<div class="ai-fix-control-group">
          <label class="ai-fix-control-label" for="ai-api-key-input">API key</label>
          <input id="ai-api-key-input" class="ai-fix-select ai-fix-input" type="password" value="${i(e.apiKey)}" placeholder="${p}" autocomplete="off">
        </div>`:"",f=e.loadingPreview?"Generating\u2026":"Generate fix",v=e.loadingPreview?" disabled":"";return`<div class="ai-fix-controls">
      <div class="ai-fix-control-group">
        <label class="ai-fix-control-label" for="ai-provider-select">Provider</label>
        <select id="ai-provider-select" class="ai-fix-select">${a}</select>
      </div>
      <div class="ai-fix-control-group">
        <label class="ai-fix-control-label" for="ai-model-input">Model</label>
        <input id="ai-model-input" class="ai-fix-select ai-fix-input" list="ai-model-options" value="${i(e.selectedModel)}" placeholder="${i(s?.default_model||"gpt-5.5")}" autocomplete="off">
        <datalist id="ai-model-options">${l}</datalist>
      </div>
      ${c}
      ${u}
      <button id="ai-generate-fix" class="ai-fix-button"${v}>${f}</button>
    </div>`}function Pe(e){if(!e.preview)return'<div class="detail-empty">Generate a fix preview to inspect the patch before Ollanta edits your local file.</div>';let t=e.preview.summary||"Generated fix preview",s=e.preview.explanation?`<div class="rule-rationale">${i(e.preview.explanation)}</div>`:"",a=e.applying?"Applying\u2026":"Apply to file",n=e.applying?" disabled":"";return`
    <div class="ai-fix-preview-meta">
      <div><strong>Provider:</strong> ${i(e.preview.agent.label)}</div>
      <div><strong>Model:</strong> ${i(e.preview.agent.model)}</div>
      <div><strong>Summary:</strong> ${i(t)}</div>
    </div>
    ${s}
    <pre class="rule-code ai-fix-diff"><code>${i(e.preview.diff)}</code></pre>
    <div class="ai-fix-actions">
      <button id="ai-apply-fix" class="ai-fix-button ai-fix-button-primary"${n}>${a}</button>
    </div>
  `}var d,y=[],M=[],S=[],L=[],j="",K=!1,W="",ie=new Map,k=null,C=-1,Le="overview",G=null,I="details",D="",ee=!1,b=null,r=fe(),P="all",le="all",oe="all",re="all",ce="all",de="",U="all",E="file",z="asc",X={blocker:0,critical:1,major:2,minor:3,info:4},N={blocker:"#ef4444",critical:"#f97316",major:"#eab308",minor:"#22c55e",info:"#64748b"},te={bug:"Bug",code_smell:"Code Smell",vulnerability:"Vulnerability",security_hotspot:"Hotspot"},we={security:"Security",reliability:"Reliability",maintainability:"Maintainability",testability:"Testability"};function T(e,t){return`<span class="icon-${i(e)}" role="img" aria-label="${i(t)}"></span>`}async function He(){try{let e=await fetch("/report.json");if(!e.ok)throw new Error(`HTTP ${e.status}`);d=await e.json(),y=d.issues??[],Ae(),je(),qe(),Ke(),Ge(),Ue(),ze(),Qe(),Oe(),Ve(),Q(),J(),L.length&&me(j||L[0].path),tt(),at(),x(),it(),Me(),st(),gt(),_e(),o("tab-issue-count").textContent=String(y.length),o("tab-file-count").textContent=String(new Set(y.map(t=>t.component_path)).size),o("tab-coverage-count").textContent=$(d.measures.coverage??d.test_signals?.summary?.coverage),o("tab-mutant-count").textContent=String(pe().survived)}catch(e){o("app").innerHTML=`<div class="error">Failed to load report: ${String(e)}</div>`}}document.addEventListener("DOMContentLoaded",He);function Ae(){let e=d.metadata,t=new Date(e.analysis_date).toLocaleString();o("project-key").textContent=e.project_key,o("scan-date").textContent=t,o("scan-version").textContent=`v${e.version}`,o("elapsed").textContent=`${e.elapsed_ms}ms`}function Fe(){let e=d.measures,t=d.test_signals?.summary,s=[{metric:"Bugs",operator:"=",threshold:0,value:e.bugs,passed:e.bugs===0},{metric:"Vulnerabilities",operator:"=",threshold:0,value:e.vulnerabilities,passed:e.vulnerabilities===0},{metric:"Code Smells",operator:"\u2264",threshold:10,value:e.code_smells,passed:e.code_smells<=10,severity:e.code_smells<=10?void 0:e.code_smells<=20?"warning":void 0}];return e.coverage!=null?s.push({metric:"Coverage",operator:"\u2265",threshold:70,value:e.coverage,passed:e.coverage>=70,severity:e.coverage>=70?void 0:e.coverage>=60?"warning":void 0}):s.push({metric:"Coverage",operator:"\u2265",threshold:70,value:0,passed:!1,severity:"missing"}),t&&(t.tests!=null&&s.push({metric:"Test Failures",operator:"=",threshold:0,value:t.test_failures??0,passed:(t.test_failures??0)===0}),t.mutation_score!=null&&s.push({metric:"Mutation Score",operator:"\u2265",threshold:60,value:t.mutation_score,passed:t.mutation_score>=60,severity:t.mutation_score>=60?void 0:t.mutation_score>=40?"warning":void 0}),t.changed_mutation_score!=null&&s.push({metric:"Changed Mutation",operator:"\u2265",threshold:60,value:t.changed_mutation_score,passed:t.changed_mutation_score>=60,severity:t.changed_mutation_score>=60?void 0:t.changed_mutation_score>=40?"warning":void 0})),{status:s.filter(l=>!l.passed&&l.severity!=="warning"&&l.severity!=="missing").length===0?"passed":"failed",conditions:s}}function je(){let e=Fe(),t=o("gate-hero");t.classList.remove("gate-loading"),t.classList.add(e.status==="passed"?"gate-passed":"gate-failed");let s=o("gate-icon");if(s.className=`gate-icon icon-${e.status==="passed"?"pass":"fail"}`,s.setAttribute("aria-label",e.status==="passed"?"Passed":"Failed"),o("gate-status").textContent=e.status==="passed"?"Passed":"Failed",e.status==="passed"){let n=e.conditions.filter(u=>!u.passed&&u.severity!=="warning");e.conditions.filter(u=>!u.passed&&u.severity==="warning").length&&!n.length&&(o("gate-status").textContent="Passed with warnings",t.classList.add("gate-warn"))}let a=e.conditions.map(n=>{let l=n.passed?"cond-pass":n.severity==="warning"?"cond-warn":"cond-fail",u=n.passed?T("pass","Passed"):T("fail","Failed");return`<div class="gate-cond ${l}">
      <span class="gate-cond-icon">${u}</span>
      <span class="gate-cond-metric">${i(n.metric)} ${i(n.operator)} ${n.threshold}</span>
      <span class="gate-cond-value">${n.value}</span>
    </div>`}).join("");o("gate-conditions").innerHTML=a}function qe(){let e=d.measures,t=d.test_signals?.summary;h("m-bugs",e.bugs),h("m-vulns",e.vulnerabilities),h("m-smells",e.code_smells),V("m-coverage",$(e.coverage)),h("m-ncloc",e.ncloc),h("m-files",e.files),h("m-comments",e.comments),t?(h("m-tests",t.tests??0),h("m-test-failures",t.test_failures??0),h("m-test-skipped",t.test_skipped??0),h("m-mutants-skipped",t.mutants_skipped??e.mutants_skipped??0),h("m-mutants-error",t.mutants_error??e.mutants_error??0)):(h("m-tests",e.tests??0),h("m-test-failures",e.test_failures??0),h("m-test-skipped",e.test_skipped??0),h("m-mutants-skipped",e.mutants_skipped??0),h("m-mutants-error",e.mutants_error??0)),O("card-bugs",e.bugs,[0,1,5]),O("card-vulns",e.vulnerabilities,[0,1,3]),O("card-smells",e.code_smells,[0,10,50]),Ne("card-coverage",e.coverage),Re("card-tests",t?.tests??e.tests,[50,20,0]),O("card-test-failures",t?.test_failures??e.test_failures??0,[0,1,5]),g("card-ncloc","card-neutral"),g("card-files","card-neutral"),g("card-comments","card-neutral"),g("card-test-skipped","card-neutral"),g("card-mutants-skipped","card-neutral"),g("card-mutants-error",(t?.mutants_error??e.mutants_error??0)>0?"card-red":"card-neutral");let s=e.duplicated_lines_density;V("m-duplication",$(s)),O("card-duplication",s??0,[3,10,20]),g("card-duplication",s==null?"card-neutral":"");let a=d.test_signals?.health;if(a){V("m-test-health",`${a.status} \xB7 ${a.score}`);let l=o("card-test-health");l.classList.remove("card-neutral","card-green","card-yellow","card-red"),a.status==="healthy"?l.classList.add("card-green"):a.status==="at_risk"?l.classList.add("card-red"):a.status==="partial"?l.classList.add("card-yellow"):l.classList.add("card-neutral")}else V("m-test-health","\u2014"),g("card-test-health","card-neutral");let n=o("card-coverage");n.classList.add("clickable"),n.addEventListener("click",()=>{F("coverage")})}function h(e,t){o(e).textContent=t.toLocaleString()}function V(e,t){o(e).textContent=t}function $(e){return e==null?"\u2014":`${e.toFixed(1)}%`}function O(e,t,s){t<=s[0]?g(e,"card-green"):t<=s[1]?g(e,"card-yellow"):g(e,"card-red")}function Ne(e,t){t==null?g(e,"card-neutral"):t>=80?g(e,"card-green"):t>=60?g(e,"card-yellow"):g(e,"card-red")}function Re(e,t,s){if(t==null){g(e,"card-neutral");return}t>=s[0]?g(e,"card-green"):t>=s[1]?g(e,"card-yellow"):g(e,"card-red")}function Oe(){let e=o("mutation-summary"),t=pe();if(!t.hasSignal){e.innerHTML='<div class="empty-state compact">No mutation report was collected for this scan. Add a supported report such as <span class="mono">ollanta-mutations.json</span>, Stryker JSON, or PIT XML to see mutation score and survived mutants.</div>';return}e.innerHTML=`<div class="mutation-kpis">
      <div><span class="mutation-kpi-value ${se(t.score)}">${$(t.score)}</span><span class="mutation-kpi-label">mutation score</span></div>
      <div><span class="mutation-kpi-value">${t.killed.toLocaleString()}/${t.total.toLocaleString()}</span><span class="mutation-kpi-label">killed mutants</span></div>
      <div><span class="mutation-kpi-value ${t.survived>0?"mut-warning":"mut-success"}">${t.survived.toLocaleString()}</span><span class="mutation-kpi-label">survived mutants</span></div>
    </div>
    ${De(t.modules)}
    ${Be(t.survivedMutants)}`}function pe(){let e=d.test_signals?.summary,t=(d.test_signals?.modules??[]).filter(p=>p.mutation),s=e?.changed_mutants_total||e?.mutants_total||d.measures.changed_mutants_total||d.measures.mutants_total||0,a=e?.changed_mutants_killed||e?.mutants_killed||d.measures.changed_mutants_killed||d.measures.mutants_killed||0,n=e?.changed_mutants_survived||e?.mutants_survived||d.measures.changed_mutants_survived||d.measures.mutants_survived||0,l=e?.changed_mutation_score??e?.mutation_score??d.measures.changed_mutation_score??d.measures.mutation_score,u=t.flatMap(p=>p.mutation?.survived_mutants??[]).slice(0,8);return{hasSignal:t.length>0||s>0||l!=null,score:l,total:s,killed:a,survived:n,modules:t,survivedMutants:u}}function De(e){return e.length?`<div class="mutation-module-list">
    ${e.slice(0,5).map(t=>{let s=t.mutation,a=s.changed_code_score??s.score,n=s.changed_survived??s.survived??0,l=s.changed_total??s.total??0;return`<div class="mutation-module-row">
        <span class="mutation-module-main"><span class="mutation-module-name">${i(t.name||t.root)}</span><span class="mutation-module-meta">${i(s.tool||"mutation")} \xB7 ${l.toLocaleString()} mutants</span></span>
        <span class="mutation-pill ${se(a)}">${$(a)}</span>
        <span class="mutation-survived ${n>0?"mut-warning":"mut-success"}">${n.toLocaleString()} survived</span>
      </div>`}).join("")}
  </div>`:""}function Be(e){return e.length?`<div class="mutation-survivors">
    ${e.map(t=>`<div class="mutation-survivor-row">
      <span class="mutation-survivor-file">${i(H(t.file||""))}${t.line?`:L${t.line}`:""}</span>
      <span class="mutation-survivor-meta">${i(t.mutator||t.description||"survived mutant")}</span>
    </div>`).join("")}
  </div>`:""}function se(e){return e==null?"card-neutral":e>=80?"card-green":e>=60?"card-yellow":"card-red"}function J(){let e=o("mutants-page"),t=pe();if(!t.hasSignal){e.innerHTML='<div class="empty-state">No mutation data collected. Run with <span class="mono">-with-mutations</span> to see survived mutants.</div>';return}let s=t.modules.flatMap(c=>(c.mutation?.survived_mutants??[]).map(f=>({...f,moduleName:c.name||c.root}))),a=[...new Set(s.map(c=>c.moduleName))].sort(),n=s;U!=="all"&&(n=n.filter(c=>c.moduleName===U)),n.sort((c,f)=>{let v=0;return E==="file"?v=(c.file||"").localeCompare(f.file||""):E==="line"?v=(c.line??0)-(f.line??0):E==="module"&&(v=c.moduleName.localeCompare(f.moduleName)),z==="asc"?v:-v});let l=`
    <div class="mutants-toolbar">
      <div class="toolbar-left">
        <select id="mutant-filter-module">
          <option value="all">All modules</option>
          ${a.map(c=>`<option value="${i(c)}"${c===U?" selected":""}>${i(c)}</option>`).join("")}
        </select>
        <select id="mutant-sort-field">
          <option value="file"${E==="file"?" selected":""}>Sort by file</option>
          <option value="line"${E==="line"?" selected":""}>Sort by line</option>
          <option value="module"${E==="module"?" selected":""}>Sort by module</option>
        </select>
        <select id="mutant-sort-dir">
          <option value="asc"${z==="asc"?" selected":""}>Ascending</option>
          <option value="desc"${z==="desc"?" selected":""}>Descending</option>
        </select>
      </div>
      <div class="toolbar-right">
        <span class="result-count">${n.length.toLocaleString()} survived</span>
      </div>
    </div>
  `,u=`
    <div class="mutation-kpis">
      <div><span class="mutation-kpi-value ${se(t.score)}">${$(t.score)}</span><span class="mutation-kpi-label">mutation score</span></div>
      <div><span class="mutation-kpi-value">${t.killed.toLocaleString()}/${t.total.toLocaleString()}</span><span class="mutation-kpi-label">killed mutants</span></div>
      <div><span class="mutation-kpi-value ${t.survived>0?"mut-warning":"mut-success"}">${t.survived.toLocaleString()}</span><span class="mutation-kpi-label">survived mutants</span></div>
    </div>
  `,p="";n.length?p=`
      <table class="mutants-table">
        <thead>
          <tr>
            <th>File</th>
            <th>Line</th>
            <th>Module</th>
            <th>Mutator</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          ${n.map(c=>`
            <tr class="mutant-row">
              <td class="mutant-file">${i(H(c.file||""))}</td>
              <td class="mutant-line">${c.line??"\u2014"}</td>
              <td class="mutant-module">${i(c.moduleName)}</td>
              <td class="mutant-mutator">${i(c.mutator||"\u2014")}</td>
              <td class="mutant-desc">${i(c.description||c.replacement||"survived mutant")}</td>
            </tr>
          `).join("")}
        </tbody>
      </table>
    `:p='<div class="empty-state compact">No survived mutants match the current filter.</div>',e.innerHTML=l+u+p,_e()}function _e(){let e=document.getElementById("mutant-filter-module"),t=document.getElementById("mutant-sort-field"),s=document.getElementById("mutant-sort-dir");e?.addEventListener("change",a=>{U=a.target.value,J()}),t?.addEventListener("change",a=>{E=a.target.value,J()}),s?.addEventListener("change",a=>{z=a.target.value,J()})}function Ve(){let e=d.test_signals?.modules??[];if(!e.length)return;let t=o("ts-modules"),s="";for(let l of e){let u=l.health,p=l.coverage,c=l.mutation,f=l.suites??[],v=f.reduce((R,B)=>R+(B.failures??0)+(B.errors??0),0),m=f.reduce((R,B)=>R+(B.tests??0),0),_=l.architecture_role??"",A=_?`<span class="ts-role-badge">${i(_)}</span>`:"",ne=$e(u?.status),xe=$(p?.coverage),Te=q(p?.coverage??void 0),Ee=$(c?.changed_code_score??c?.score),Ce=se(c?.changed_code_score??c?.score);s+=`<div class="ts-module">
      <div class="ts-module-head">
        <span class="ts-module-name">${i(l.name)}</span>
        ${A}
        <span class="ts-health-badge ${ne}">${u?.status??"no data"}</span>
      </div>
      <div class="ts-module-meta">
        ${p?`<span class="ts-metric ${Te}">Coverage ${xe}</span>`:""}
        ${c?.score!=null||c?.changed_code_score!=null?`<span class="ts-metric ${Ce}">Mutation ${Ee}</span>`:""}
        ${m>0?`<span class="ts-metric">${m} test${m===1?"":"s"}</span>`:""}
        ${v>0?`<span class="ts-metric ts-fail">${v} failed</span>`:""}
      </div>
      ${u?.recommendations?.length?`<div class="ts-recommendations">${u.recommendations.map(R=>`<div class="ts-rec">${i(R)}</div>`).join("")}</div>`:""}
    </div>`}let a=d.test_signals?.health,n=a?`<span class="ts-health-badge ${$e(a.status)}">${a.status} \xB7 score ${a.score}</span>`:"";t.innerHTML=`<div class="ts-header">
      <h3>Test Signals</h3>
      ${n}
    </div>
    <div class="ts-module-list">${s||'<div class="empty-state compact">No module-level test data was collected.</div>'}</div>`}function $e(e){return e==="healthy"?"card-green":e==="at_risk"?"card-red":e==="partial"?"card-yellow":"card-neutral"}function Ke(){let e=Z(y,v=>v.severity),t=Math.max(1,...Object.values(e)),s=["blocker","critical","major","minor","info"],a="",n="",l=y.length||1;for(let v of s){let m=e[v]??0,_=m/t*100,A=N[v]??"#64748b";a+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${v}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${_}%;background:${A}"></div></div>
      <span class="sev-bar-count">${m}</span>
    </div>`,m>0&&(n+=`<div class="sev-segment" style="width:${m/l*100}%;background:${A}" title="${v}: ${m}"></div>`)}o("sev-bars").innerHTML=a,o("sev-proportional").innerHTML=n;let u=Z(y,v=>v.type),p=Math.max(1,...Object.values(u)),c={bug:"#ef4444",vulnerability:"#f97316",code_smell:"#22c55e",security_hotspot:"#eab308"},f="";for(let[v,m]of Object.entries(te)){let _=u[v]??0,A=_/p*100,ne=c[v]??"#64748b";f+=`<div class="sev-bar-row">
      <span class="sev-bar-label">${m}</span>
      <div class="sev-bar-track"><div class="sev-bar-fill" style="width:${A}%;background:${ne}"></div></div>
      <span class="sev-bar-count">${_}</span>
    </div>`}o("type-bars").innerHTML=f}function Ge(){let e=[...y].sort((t,s)=>{let a=(X[t.severity]??99)-(X[s.severity]??99);return a!==0?a:t.component_path.localeCompare(s.component_path)||t.line-s.line}).slice(0,6);if(!e.length){o("priority-issues").innerHTML='<div class="empty-state compact">No issues found</div>';return}o("priority-issues").innerHTML=e.map((t,s)=>{let a=N[t.severity]??"#64748b",n=H(t.component_path);return`<button class="priority-row" data-idx="${s}">
      <span class="issue-sev-dot" style="background:${a}"></span>
      <span class="priority-main">
        <span class="priority-title">${i(t.message)}</span>
        <span class="priority-meta" title="${i(t.component_path)}">${i(n)}:L${t.line} \xB7 ${i(t.rule_key)}</span>
      </span>
      <span class="priority-severity">${i(t.severity)}</span>
    </button>`}).join(""),o("priority-issues").querySelectorAll(".priority-row").forEach(t=>{t.addEventListener("click",()=>{let s=Number.parseInt(t.dataset.idx,10);ge(e[s])})})}function Ue(){let e=Z(y,s=>s.component_path),t=Object.entries(e).sort((s,a)=>a[1]-s[1]).slice(0,10);if(!t.length){o("hotspot-files").innerHTML='<div class="empty-state">No issues found</div>';return}o("hotspot-files").innerHTML=t.map(([s,a])=>{let n=H(s);return`<div class="hotspot-row" data-path="${i(s)}">
      <span class="hotspot-file" title="${i(s)}">${i(n)}</span>
      <span class="hotspot-count">${a}</span>
    </div>`}).join(""),o("hotspot-files").querySelectorAll(".hotspot-row").forEach(s=>{s.addEventListener("click",()=>{let a=s.dataset.path;F("files"),lt(a)})})}function ze(){L=(d.test_signals?.modules??[]).flatMap(t=>(t.files??[]).map(s=>Je(t.name,t.root,s))).filter(t=>t.linesToCover>0).sort((t,s)=>(t.coverage??101)-(s.coverage??101)||s.uncoveredLines.length-t.uncoveredLines.length||t.path.localeCompare(s.path)),!j&&L.length&&(j=L[0].path)}function Je(e,t,s){let a=s.lines_to_cover??0,n=s.covered_lines??0,l=a>0?n*100/a:null;return{moduleName:e,moduleRoot:t,path:s.path,linesToCover:a,coveredLines:n,coveredLineNumbers:s.covered_line_numbers??[],uncoveredLines:s.uncovered_lines??[],coverage:l}}function Qe(){let e=o("coverage-summary");if(!L.length){e.innerHTML='<div class="empty-state compact">Run with <span class="mono">-with-tests</span> and provide a coverage report to see file-level details.</div>';return}let t=d.test_signals?.summary,s=L.slice(0,5);e.innerHTML=`<div class="coverage-kpis">
      <div><span class="coverage-kpi-value">${$(t?.coverage??d.measures.coverage)}</span><span class="coverage-kpi-label">overall</span></div>
      <div><span class="coverage-kpi-value">${(t?.covered_lines??0).toLocaleString()}/${(t?.lines_to_cover??0).toLocaleString()}</span><span class="coverage-kpi-label">covered lines</span></div>
      <div><span class="coverage-kpi-value">${(t?.modules_with_coverage??0).toLocaleString()}</span><span class="coverage-kpi-label">modules</span></div>
    </div>
    <div class="coverage-file-list">
      ${s.map(We).join("")}
    </div>`,e.querySelectorAll(".coverage-mini-row").forEach(a=>{a.addEventListener("click",()=>{let n=a.dataset.coveragePath;n&&(F("coverage"),me(n))})})}function We(e){return`<button class="coverage-mini-row" data-coverage-path="${i(e.path)}">
    <span class="coverage-mini-main">
      <span class="coverage-file-name" title="${i(e.path)}">${i(H(e.path))}</span>
      <span class="coverage-file-meta">${i(e.moduleName)} \xB7 ${e.uncoveredLines.length.toLocaleString()} uncovered lines</span>
    </span>
    <span class="coverage-pill ${q(e.coverage)}">${$(e.coverage??void 0)}</span>
  </button>`}function Q(){let e=o("coverage-details");if(!L.length){e.innerHTML='<div class="empty-state">No file-level coverage was collected for this scan.</div>';return}e.innerHTML=`<div class="coverage-toolbar">
      <div>
        <h3>Coverage Files</h3>
        <p>${L.length.toLocaleString()} files with line-level detail. Overall coverage includes all measured files, not only those listed here.</p>
      </div>
      <span class="coverage-pill ${q(d.measures.coverage??null)}">${$(d.measures.coverage)}</span>
    </div>
    <div class="coverage-browser">
      <aside class="coverage-sidebar">
        <div class="coverage-file-list coverage-file-list-large">
          ${L.map(Xe).join("")}
        </div>
      </aside>
      <section class="coverage-code-viewer">
        ${Ye()}
      </section>
    </div>`,e.querySelectorAll(".coverage-row").forEach(t=>{t.addEventListener("click",()=>{let s=t.dataset.coveragePath;s&&me(s)})})}function Xe(e){return`<button class="coverage-row${e.path===j?" active":""}" data-coverage-path="${i(e.path)}">
    <div class="coverage-row-main">
      <div class="coverage-row-title" title="${i(e.path)}">${i(e.path)}</div>
      <div class="coverage-row-subtitle">${i(e.moduleName)} \xB7 ${i(e.moduleRoot)} \xB7 ${e.coveredLines.toLocaleString()}/${e.linesToCover.toLocaleString()} lines covered</div>
    </div>
    <div class="coverage-row-meter">
      <span class="coverage-pill ${q(e.coverage)}">${$(e.coverage??void 0)}</span>
      <div class="coverage-track"><div class="coverage-fill ${q(e.coverage)}" style="width:${e.coverage??0}%"></div></div>
    </div>
  </button>`}async function me(e){if(L.some(t=>t.path===e)){if(j=e,W="",ie.has(e)){K=!1,Q();return}K=!0,Q();try{let t=await fetch(`/api/files/source?path=${encodeURIComponent(e)}`);if(!t.ok)throw new Error(`HTTP ${t.status}`);let s=await t.json();ie.set(e,s.file)}catch(t){W=`Could not load source for ${e}: ${String(t)}`}finally{K=!1,Q()}}}function Ye(){let e=L.find(p=>p.path===j);if(!e)return'<div class="code-empty"><p>Select a file to inspect coverage.</p></div>';if(K)return'<div class="code-empty"><div class="spinner"></div></div>';if(W)return`<div class="code-empty"><p>${i(W)}</p></div>`;let t=ie.get(e.path);if(!t)return'<div class="code-empty"><p>Select a coverage file to inspect source lines.</p></div>';let s=new Set(e.coveredLineNumbers),a=new Set(e.uncoveredLines),l=t.content.split(`
`).map((p,c)=>{let f=c+1,v=a.has(f),m=!v&&s.has(f),_=Ze(m,v);return`<div class="code-line${_.stateClass}">
      <span class="code-gutter">${f}</span>
      <code class="code-text">${p.length?i(p):"&nbsp;"}</code>
      <span class="code-markers">${et(_)}</span>
    </div>`}).join(""),u=e.coveredLineNumbers.length?"covered and uncovered lines":"uncovered lines only";return`<div class="code-viewer-shell coverage-source-shell">
    <div class="code-viewer-head">
      <div>
        <div class="code-viewer-path mono">${i(t.path)}</div>
        <div class="code-viewer-meta">${i(t.language||"plain text")} \xB7 ${t.line_count.toLocaleString()} lines \xB7 ${u}</div>
      </div>
      <div class="code-viewer-stats"><span class="coverage-pill ${q(e.coverage)}">${$(e.coverage??void 0)}</span></div>
    </div>
    <div class="coverage-legend">
      <span><span class="legend-dot legend-covered"></span>Covered</span>
      <span><span class="legend-dot legend-uncovered"></span>Not covered</span>
    </div>
    <div class="code-surface">${l}</div>
  </div>`}function Ze(e,t){return t?{stateClass:" is-uncovered",marker:"not covered",chipClass:" chip-uncovered"}:e?{stateClass:" is-covered",marker:"covered",chipClass:" chip-covered"}:{stateClass:"",marker:"",chipClass:""}}function et(e){return e.marker?`<span class="coverage-line-chip${e.chipClass}">${e.marker}</span>`:""}function q(e){return e==null?"card-neutral":e>=80?"card-green":e>=60?"card-yellow":"card-red"}function tt(){let e=Object.entries(d.measures.by_language).sort((s,a)=>a[1]-s[1]),t=Math.max(1,e[0]?.[1]??1);if(!e.length){o("by-lang").innerHTML='<span class="empty-state">No language data</span>';return}o("by-lang").innerHTML=e.map(([s,a])=>`<div class="lang-row">
      <span class="lang-name">${i(s)}</span>
      <div class="lang-bar-track"><div class="lang-bar-fill" style="width:${a/t*100}%"></div></div>
      <span class="lang-count">${a} files</span>
    </div>`).join("")}function st(){document.querySelectorAll(".tab").forEach(t=>{t.addEventListener("click",()=>{let s=t.dataset.tab;F(s)})}),document.querySelector(".tabs").addEventListener("keydown",t=>{let s=Array.from(document.querySelectorAll(".tab[role='tab']")),a=s.findIndex(n=>n.getAttribute("aria-selected")==="true");if(t.key==="ArrowRight"){t.preventDefault();let n=s[(a+1)%s.length];n.focus(),F(n.dataset.tab)}else if(t.key==="ArrowLeft"){t.preventDefault();let n=s[(a-1+s.length)%s.length];n.focus(),F(n.dataset.tab)}})}function F(e){Le=e,document.querySelectorAll(".tab").forEach(s=>{s.classList.remove("active"),s.setAttribute("aria-selected","false")});let t=document.querySelector(`.tab[data-tab="${e}"]`);t?.classList.add("active"),t?.setAttribute("aria-selected","true"),document.querySelectorAll(".panel").forEach(s=>s.classList.add("hidden")),o(`panel-${e}`).classList.remove("hidden")}function at(){let e=[...new Set(y.map(n=>n.rule_key))].sort((n,l)=>n.localeCompare(l)),t=o("filter-rule");e.forEach(n=>{let l=document.createElement("option");l.value=n,l.textContent=n,t.appendChild(l)});let s=new Set;for(let n of y)for(let l of n.tags??[])s.add(l);let a=o("filter-tag");[...s].sort().forEach(n=>{let l=document.createElement("option");l.value=n,l.textContent=n,a.appendChild(l)}),o("filter-severity").addEventListener("change",n=>{P=n.target.value,x()}),o("filter-type").addEventListener("change",n=>{le=n.target.value,x()}),t.addEventListener("change",n=>{oe=n.target.value,x()}),o("filter-quality").addEventListener("change",n=>{re=n.target.value,x()}),a.addEventListener("change",n=>{ce=n.target.value,x()}),o("search").addEventListener("input",n=>{de=n.target.value.toLowerCase(),x()}),Se()}function Se(){let e=Z(y,s=>s.severity),t=["blocker","critical","major","minor","info"];o("sev-chips").innerHTML=t.map(s=>{let a=e[s]??0,n=N[s]??"#64748b",l=P===s?" active":"";return`<button class="sev-chip${l}" data-sev="${s}"
      style="--chip-color:${n};--chip-bg:${n}15" aria-pressed="${l?"true":"false"}">
      <span class="chip-dot" style="background:${n}"></span>
      ${s}
      <span class="chip-count">${a}</span>
    </button>`}).join(""),o("sev-chips").querySelectorAll(".sev-chip").forEach(s=>{s.addEventListener("click",()=>{let a=s.dataset.sev;P=P===a?"all":a,o("filter-severity").value=P,x(),Se()})})}function x(){M=y.filter(s=>!(P!=="all"&&s.severity!==P||le!=="all"&&s.type!==le||oe!=="all"&&s.rule_key!==oe||re!=="all"&&s.quality!==re||ce!=="all"&&!(s.tags??[]).includes(ce)||de&&!`${s.component_path} ${s.message} ${s.rule_key}`.toLowerCase().includes(de))),M.sort((s,a)=>{let n=X[s.severity]??99,l=X[a.severity]??99;return n-l}),C=-1,nt();let e=M.length,t=document.getElementById("filter-announcer");t&&(t.textContent=`${e} issue${e===1?"":"s"} match the current filters`)}function nt(){let e=o("issue-list"),t=M.length===1?"issue":"issues";if(o("issue-count").textContent=`${M.length} ${t}`,!M.length){e.innerHTML='<div class="empty-state">No issues match the current filters.</div>';return}e.innerHTML=M.map((s,a)=>{let n=N[s.severity]??"#64748b",l=H(s.component_path),u=s.end_line&&s.end_line!==s.line?`L${s.line}\u2013${s.end_line}`:`L${s.line}`,p=te[s.type]??s.type,c=s.quality?`<span class="quality-badge quality-${i(s.quality)}">${i(we[s.quality]??s.quality)}</span>`:"";return`<div class="issue-row" role="button" tabindex="0" aria-label="${i(s.severity)} issue: ${i(s.message)}" data-idx="${a}">
      <span class="issue-sev">
        <span class="issue-sev-dot" style="background:${n}"></span>
        ${i(s.severity)}
      </span>
      <span class="issue-type">${i(p)}</span>
      <div class="issue-main">
        <span class="issue-msg">${i(s.message)}</span>
        <span class="issue-file" title="${i(s.component_path)}">${i(l)}:${u}</span>
      </div>
      ${c}
      <span class="issue-rule">${i(s.rule_key)}</span>
    </div>`}).join(""),e.querySelectorAll(".issue-row").forEach(s=>{s.addEventListener("click",()=>{let a=Number.parseInt(s.dataset.idx,10);Y(a)}),s.addEventListener("keydown",a=>{let n=a;if(n.key==="Enter"||n.key===" "){n.preventDefault();let l=Number.parseInt(s.dataset.idx,10);Y(l)}})})}function it(){let e=new Map;for(let t of y){let s=t.component_path;e.has(s)||e.set(s,[]),e.get(s).push(t)}S=[...e.entries()].sort((t,s)=>s[1].length-t[1].length).map(([t,s])=>({path:t,shortPath:H(t),issues:[...s].sort((a,n)=>a.line-n.line),expanded:!1}))}function Me(){let e=o("file-tree");if(!S.length){e.innerHTML='<div class="empty-state">No issues found</div>';return}e.innerHTML=S.map((t,s)=>`<div class="file-group${t.expanded?" expanded":""}" data-gi="${s}">
      <div class="file-group-header">
        <span class="file-group-chevron icon-chevron" role="img" aria-label="Expand"></span>
        <span class="file-group-name" title="${i(t.path)}">${i(t.shortPath)}</span>
        <span class="file-group-count">${t.issues.length}</span>
      </div>
      <div class="file-group-issues" style="${t.expanded?"":"display:none"}">
        ${t.issues.map((a,n)=>{let l=N[a.severity]??"#64748b";return`<div class="file-issue" data-gi="${s}" data-ii="${n}">
            <span class="issue-sev">
              <span class="issue-sev-dot" style="background:${l}"></span>
              ${i(a.severity)}
            </span>
            <span class="issue-msg">${i(a.message)}</span>
            <span class="file-issue-line">L${a.line}</span>
          </div>`}).join("")}
      </div>
    </div>`).join(""),e.querySelectorAll(".file-group-header").forEach(t=>{t.addEventListener("click",()=>{let s=t.closest(".file-group"),a=Number.parseInt(s.dataset.gi,10);S[a].expanded=!S[a].expanded,s.classList.toggle("expanded");let n=s.querySelector(".file-group-issues");n.style.display=S[a].expanded?"":"none"})}),e.querySelectorAll(".file-issue").forEach(t=>{t.addEventListener("click",s=>{s.stopPropagation();let a=Number.parseInt(t.dataset.gi,10),n=Number.parseInt(t.dataset.ii,10),l=S[a].issues[n];ge(l)})})}function lt(e){let t=S.findIndex(a=>a.path===e);if(t<0)return;S[t].expanded=!0,Me(),document.querySelector(`.file-group[data-gi="${t}"]`)?.scrollIntoView({behavior:"smooth",block:"start"})}function Y(e,t=!0){C=e,k=M[e]??null,document.querySelectorAll(".issue-row").forEach(a=>a.classList.remove("selected"));let s=document.querySelector(`.issue-row[data-idx="${e}"]`);s?.classList.add("selected"),s?.focus(),t&&k&&ge(k)}function ge(e){G=document.activeElement,k=e,I="details",D="",ee=!0,r=fe(),o("detail-title").textContent=e.message||e.rule_key,ae(e),o("detail-panel").classList.add("open"),o("detail-overlay").classList.add("open"),o("detail-panel").querySelector("button, [href], input, select, textarea, [tabindex]:not([tabindex='-1'])")?.focus(),ot(e.rule_key)}async function ot(e){try{let t=await fetch(`/rules/${encodeURIComponent(e)}`);if(!t.ok)throw new Error("not found");let s=await t.json(),a="";s.rationale&&(a+=`<div class="detail-section">
        <div class="detail-section-title">Why is this a problem?</div>
        <div class="rule-rationale">${i(s.rationale)}</div>
      </div>`),s.description&&s.description!==s.rationale&&(a+=`<div class="detail-section">
        <div class="detail-section-title">Description</div>
        <div class="rule-rationale">${i(s.description)}</div>
      </div>`),s.noncompliant_code&&(a+=`<div class="detail-section">
        <div class="detail-section-title">${T("cross","Noncompliant")} Noncompliant Code</div>
        <pre class="rule-code noncompliant"><code>${i(s.noncompliant_code)}</code></pre>
      </div>`),s.compliant_code&&(a+=`<div class="detail-section">
        <div class="detail-section-title">${T("check","Compliant")} Compliant Code</div>
        <pre class="rule-code compliant"><code>${i(s.compliant_code)}</code></pre>
      </div>`),D=a||'<div class="detail-empty">No additional rule details available.</div>'}catch{D='<div class="detail-empty">Rule details are not available for this issue.</div>'}finally{ee=!1,k?.rule_key===e&&ae(k)}}function rt(e){document.getElementById("detailCopy")?.addEventListener("click",()=>{ct(e)})}async function ct(e){let t=[];t.push(`Issue: ${e.message||""}`),t.push(`Severity: ${e.severity}`),t.push(`Type: ${te[e.type]??e.type}`),t.push(`Rule: ${e.rule_key}`),e.engine_id&&t.push(`Engine: ${e.engine_id}`),t.push(`File: ${e.component_path}`);let s=e.end_line&&e.end_line!==e.line?`lines ${e.line}\u2013${e.end_line}`:`line ${e.line}`;t.push(`Location: ${s}${e.column?", column "+e.column:""}`),t.push(`Status: ${e.status}`),e.tags?.length&&t.push(`Tags: ${e.tags.join(", ")}`);try{let a=await fetch(`/rules/${encodeURIComponent(e.rule_key)}`);if(a.ok){let n=await a.json();n.rationale&&t.push(`
Why is this a problem?
${n.rationale}`),n.noncompliant_code&&t.push(`
Noncompliant code:
${n.noncompliant_code}`),n.compliant_code&&t.push(`
Compliant code:
${n.compliant_code}`)}}catch{}try{await navigator.clipboard.writeText(t.join(`
`));let a=document.getElementById("detailCopy");a&&(a.innerHTML=`${T("check","Copied")} Copied`,setTimeout(()=>{a.innerHTML=`${T("copy","Copy")} Copy`},2e3))}catch{let a=document.getElementById("detailCopy");a&&(a.innerHTML=`${T("warn","Failed")} Failed`,setTimeout(()=>{a.innerHTML=`${T("copy","Copy")} Copy`},2e3))}}function ue(){o("detail-panel").classList.remove("open"),o("detail-overlay").classList.remove("open"),k=null,D="",ee=!1,r=fe(),document.querySelectorAll(".issue-row").forEach(e=>e.classList.remove("selected")),G&&(G.focus(),G=null)}function ae(e){let t=`
    <div class="detail-tabs">
      ${he(I)}
    </div>
    <div class="detail-tab-panel${I==="details"?"":" hidden"}" data-detail-panel="details">
      ${dt(e)}
    </div>
    <div class="detail-tab-panel${I==="rule"?"":" hidden"}" data-detail-panel="rule">
      ${ee?'<div class="detail-loading">Loading rule details\u2026</div>':D}
    </div>
    <div class="detail-tab-panel${I==="ai-fix"?"":" hidden"}" data-detail-panel="ai-fix">
      ${ye(e,r,b??[])}
    </div>
  `;o("detail-body").innerHTML=t,ut(e),rt(e)}function dt(e){let t=N[e.severity]??"#64748b",s=te[e.type]??e.type,a=e.end_line&&e.end_line!==e.line?`${e.line}:${e.column} \u2013 ${e.end_line}:${e.end_column}`:`${e.line}:${e.column}`,n=`
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
      ${e.quality?`<div class="detail-field">
        <span class="detail-field-label">Quality</span>
        <span class="detail-field-value"><span class="quality-badge quality-${i(e.quality)}">${i(we[e.quality]??e.quality)}</span></span>
      </div>`:""}
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
        <span class="detail-field-value">${e.tags.map(l=>i(l)).join(", ")}</span>
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
        <span class="detail-field-value" style="font-family:var(--font-mono)">${a}</span>
      </div>
    </div>`;return e.secondary_locations?.length&&(n+=`<div class="detail-section">
      <div class="detail-section-title">Related Locations (${e.secondary_locations.length})</div>
      <div class="detail-loc-list">
        ${e.secondary_locations.map(l=>`
          <div class="detail-loc-item">
            <div class="detail-loc-file">${i(l.file_path||e.component_path)}:${l.start_line}</div>
            ${l.message?`<div class="detail-loc-msg">${i(l.message)}</div>`:""}
          </div>
        `).join("")}
      </div>
    </div>`),n}function ut(e){document.querySelectorAll(".detail-tab").forEach(n=>{n.addEventListener("click",()=>{I=n.dataset.detailTab??"details",ae(e),I==="ai-fix"&&vt()})});let t=document.getElementById("ai-provider-select");t?.addEventListener("change",()=>{r.selectedProviderId=t.value,r.selectedModel="",ve(),r.preview=null,r.statusMessage="",r.errorMessage="",w()});let s=document.getElementById("ai-model-input");s?.addEventListener("input",()=>{r.selectedModel=s.value});let a=document.getElementById("ai-api-key-input");a?.addEventListener("input",()=>{r.apiKey=a.value}),document.getElementById("ai-generate-fix")?.addEventListener("click",()=>{pt(e)}),document.getElementById("ai-apply-fix")?.addEventListener("click",()=>{mt()})}function fe(){return{loadingOptions:!1,loadingPreview:!1,applying:!1,selectedProviderId:"",selectedModel:"",apiKey:"",statusMessage:"",errorMessage:"",preview:null}}function ke(){return!b||b.length===0?null:b.find(e=>e.id===r.selectedProviderId)??b[0]}function ve(){if(!b||b.length===0){r.selectedProviderId="",r.selectedModel="";return}b.some(t=>t.id===r.selectedProviderId)||(r.selectedProviderId=b[0].id);let e=ke();if(!e){r.selectedModel="";return}r.selectedModel||(r.selectedModel=e.default_model||e.models[0]||"")}async function vt(){if(b){ve(),w();return}r.loadingOptions=!0,r.errorMessage="",w();try{let e=await fetch("/api/ai/providers");if(!e.ok)throw new Error(`HTTP ${e.status}`);b=(await e.json()).providers??[],ve()}catch(e){r.errorMessage=`Failed to load AI models: ${String(e)}`,b=[]}finally{r.loadingOptions=!1,w()}}async function pt(e){let t=ke(),s=r.selectedModel.trim();if(!t||!r.selectedProviderId){r.errorMessage="Choose an AI provider before generating a fix.",w();return}if(!s){r.errorMessage="Choose a model before generating a fix.",w();return}if(t.requires_api_key&&!t.configured&&!r.apiKey.trim()){r.errorMessage="Provide an API key for the selected provider before generating a fix.",w();return}r.selectedModel=s,r.loadingPreview=!0,r.statusMessage="",r.errorMessage="",w();try{let a={provider:r.selectedProviderId,model:s,api_key:r.apiKey.trim()||void 0,issue:e},n=await fetch("/api/ai/fixes/preview",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(a)}),l=await n.json();if(!n.ok||"error"in l)throw new Error("error"in l?l.error:`HTTP ${n.status}`);r.preview=l,r.statusMessage="Fix preview generated. Review the diff before applying it."}catch(a){r.errorMessage=`Failed to generate AI fix: ${String(a)}`,r.preview=null}finally{r.loadingPreview=!1,w()}}async function mt(){if(r.preview){r.applying=!0,r.errorMessage="",w();try{let e=await fetch("/api/ai/fixes/apply",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({preview_id:r.preview.preview_id})}),t=await e.json();if(!e.ok||"error"in t)throw new Error("error"in t?t.error:`HTTP ${e.status}`);r.statusMessage=t.message}catch(e){r.errorMessage=`Failed to apply AI fix: ${String(e)}`}finally{r.applying=!1,w()}}}function w(){k&&ae(k)}document.addEventListener("DOMContentLoaded",()=>{o("detail-close").addEventListener("click",ue),o("detail-overlay").addEventListener("click",ue)});function gt(){document.addEventListener("keydown",e=>{let t=e.target.tagName;if(!(t==="INPUT"||t==="SELECT"||t==="TEXTAREA")){if(e.key==="Escape"){ue();return}Le==="issues"&&(e.key==="j"||e.key==="ArrowDown"?(e.preventDefault(),C<M.length-1&&Y(C+1,!1),be()):(e.key==="k"||e.key==="ArrowUp")&&(e.preventDefault(),C>0&&Y(C-1,!1),be()))}})}function be(){document.querySelector(`.issue-row[data-idx="${C}"]`)?.scrollIntoView({behavior:"smooth",block:"nearest"})}function o(e){return document.getElementById(e)}function g(e,t){o(e).classList.add(t)}function Z(e,t){let s={};for(let a of e){let n=t(a);s[n]=(s[n]??0)+1}return s}function H(e){let t=e.replaceAll("\\","/"),s=t.split("/").filter(Boolean);return s.length<=2?t:`${s.slice(-2).join("/")}`}})();
