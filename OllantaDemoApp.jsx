import React, { useState, useEffect } from 'react';
import { Search, LogOut, User, LayoutGrid, AlertCircle, ShieldAlert, Code2, Layers, FileCode2, ChevronLeft, CheckCircle2, XCircle } from 'lucide-react';

// Componente de Animação Sequencial (GIF)
const SequenceAnimation = ({ images = ['login.png', 'projects.png', 'server-dash.png'] }) => {
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setActiveIndex((prev) => (prev + 1) % images.length);
    }, 3000); // 3 segundos por imagem
    return () => clearInterval(interval);
  }, [images.length]);

  return (
    <div className="relative w-full rounded-xl overflow-hidden border border-gray-800 bg-[#161b22]">
      {/* Container de imagens com aspect ratio 16:9 */}
      <div className="relative w-full" style={{ paddingBottom: '56.25%' }}>
        {images.map((img, idx) => (
          <img
            key={idx}
            src={`docs/imgs/${img}`}
            alt={`Frame ${idx + 1}`}
            className="absolute inset-0 w-full h-full object-cover transition-opacity duration-800"
            style={{ opacity: activeIndex === idx ? 1 : 0 }}
          />
        ))}
      </div>

      {/* Indicadores (dots) */}
      <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-2 z-10">
        {images.map((_, idx) => (
          <div
            key={idx}
            className="w-2.5 h-2.5 rounded-full transition-colors duration-300"
            style={{
              backgroundColor: activeIndex === idx ? '#34d399' : '#475569'
            }}
          />
        ))}
      </div>
    </div>
  );
};

const App = () => {
  const [step, setStep] = useState('login'); // login, loading, projects, dashboard
  const [isAutoPlaying, setIsAutoPlaying] = useState(true);

  // Efeito para simular a navegação automática (estilo GIF)
  useEffect(() => {
    if (!isAutoPlaying) return;

    let timer;
    if (step === 'login') {
      timer = setTimeout(() => setStep('loading'), 2000);
    } else if (step === 'loading') {
      timer = setTimeout(() => setStep('projects'), 1000);
    } else if (step === 'projects') {
      timer = setTimeout(() => setStep('dashboard'), 2000);
    } else if (step === 'dashboard') {
      timer = setTimeout(() => setStep('login'), 5000); // Reinicia o loop
    }

    return () => clearTimeout(timer);
  }, [step, isAutoPlaying]);

  const Navbar = () => (
    <div className="flex items-center justify-between px-6 py-3 bg-[#0d1117] border-b border-gray-800 text-gray-300">
      <div className="flex items-center gap-2 text-[#5865f2] font-bold text-lg">
        <Search size={20} />
        <span>Ollanta</span>
      </div>
      <div className="flex items-center gap-4 text-sm">
        <div className="flex items-center gap-1 cursor-pointer hover:text-white transition-colors">
          <span>Administrator</span>
        </div>
        <button 
          onClick={() => setStep('login')}
          className="flex items-center gap-1 border border-gray-700 px-3 py-1 rounded hover:bg-gray-800 transition-colors"
        >
          Sign out
        </button>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-[#0d1117] text-gray-100 font-sans selection:bg-indigo-500/30">
      
      {/* Botão de Controle (Opcional) */}
      <div className="fixed bottom-4 right-4 z-50 flex gap-2">
        <button 
          onClick={() => setIsAutoPlaying(!isAutoPlaying)}
          className="bg-[#5865f2] hover:bg-[#4752c4] text-white px-4 py-2 rounded-full text-sm font-medium shadow-lg transition-all"
        >
          {isAutoPlaying ? "Pausar Loop" : "Retomar Auto-play"}
        </button>
      </div>

      {step === 'login' && (
        <div className="flex items-center justify-center min-h-screen animate-in fade-in duration-700">
          <div className="w-full max-w-md p-8 bg-[#161b22] rounded-xl border border-gray-800 shadow-2xl">
            <div className="mb-8 space-y-2">
              <div className="flex items-center gap-2 text-[#5865f2] font-bold text-2xl">
                <Search size={28} />
                <span>Ollanta</span>
              </div>
              <p className="text-gray-400 text-sm">Static analysis platform</p>
            </div>

            <div className="space-y-6">
              <div className="space-y-2">
                <label className="text-xs font-bold text-gray-400 uppercase tracking-wider">Username</label>
                <div className="w-full bg-[#0d1117] border border-gray-700 p-3 rounded-lg text-gray-200">
                  {isAutoPlaying ? <span className="animate-pulse">admin</span> : "admin"}
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-xs font-bold text-gray-400 uppercase tracking-wider">Password</label>
                <div className="w-full bg-[#0d1117] border border-gray-700 p-3 rounded-lg text-gray-200">
                   ••••••
                </div>
              </div>

              <button 
                className={`w-full py-3 rounded-lg font-bold transition-all ${isAutoPlaying ? 'bg-[#5865f2] scale-95 opacity-80' : 'bg-[#5865f2] hover:bg-[#4752c4]'}`}
              >
                Sign in
              </button>
            </div>
          </div>
        </div>
      )}

      {step === 'loading' && (
        <div className="flex items-center justify-center min-h-screen">
          <div className="flex flex-col items-center gap-4">
            <div className="w-10 h-10 border-4 border-[#5865f2] border-t-transparent rounded-full animate-spin"></div>
            <p className="text-gray-400 animate-pulse">Autenticando...</p>
          </div>
        </div>
      )}

      {step === 'projects' && (
        <div className="animate-in slide-in-from-bottom-4 duration-500">
          <Navbar />
          <main className="max-w-7xl mx-auto p-8">
            <div className="mb-10">
              <h1 className="text-2xl font-bold mb-1 flex items-center gap-2">
                Projects <span className="text-gray-500 font-normal text-lg">(1)</span>
              </h1>
              <p className="text-gray-400 text-sm">All projects registered on this platform</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div 
                onClick={() => setStep('dashboard')}
                className={`group cursor-pointer p-6 bg-[#161b22] border border-gray-800 rounded-xl hover:border-gray-600 transition-all ${isAutoPlaying ? 'ring-2 ring-[#5865f2] ring-offset-4 ring-offset-[#0d1117]' : ''}`}
              >
                <span className="text-[#5865f2] text-xs font-medium mb-2 block">project</span>
                <h3 className="text-xl font-bold mb-4 group-hover:text-[#5865f2] transition-colors">project</h3>
                <div className="flex items-center justify-between mt-8">
                  <span className="text-gray-500 text-xs italic">Updated 16h ago</span>
                </div>
              </div>
            </div>
          </main>
        </div>
      )}

      {step === 'dashboard' && (
        <div className="animate-in fade-in duration-500">
          <Navbar />
          <main className="max-w-7xl mx-auto p-6 space-y-6 pb-20">
            {/* Header */}
            <div className="flex flex-col gap-2">
              <button 
                onClick={() => setStep('projects')}
                className="flex items-center gap-1 text-[#5865f2] text-sm hover:underline w-fit"
              >
                <ChevronLeft size={16} /> Projects
              </button>
              <div className="flex items-center gap-3">
                <h1 className="text-3xl font-bold">project</h1>
                <span className="bg-red-500/20 text-red-500 text-[10px] font-bold px-2 py-0.5 rounded border border-red-500/30">ERROR</span>
              </div>
              <p className="text-gray-500 text-sm">project</p>
            </div>

            {/* Tabs */}
            <div className="flex border-b border-gray-800 gap-8">
              {['Overview', 'Issues (505)', 'Activity', 'Quality Gate', 'Webhooks', 'Profiles'].map((tab, idx) => (
                <button 
                  key={tab}
                  className={`pb-3 text-sm transition-colors ${idx === 0 ? 'text-[#5865f2] border-b-2 border-[#5865f2]' : 'text-gray-500 hover:text-gray-300'}`}
                >
                  {tab}
                </button>
              ))}
            </div>

            {/* Animação Sequencial - Server Screens */}
            <div>
              <div className="mb-4">
                <h3 className="text-sm font-bold text-gray-400 uppercase tracking-widest">Server Dashboard Preview</h3>
              </div>
              <SequenceAnimation images={['login.png', 'projects.png', 'server-dash.png']} />
              <p className="text-gray-500 text-xs mt-2">
                Visualização sequencial: Login → Projects → Dashboard
              </p>
            </div>

            {/* Main Warning Card */}
            <div className="p-6 bg-[#1c1414] border border-red-900/50 rounded-xl flex items-center gap-6">
              <div className="w-12 h-12 bg-red-600 rounded-full flex items-center justify-center flex-shrink-0">
                <span className="text-2xl font-bold">X</span>
              </div>
              <div className="flex-grow">
                <p className="text-gray-400 text-xs font-bold uppercase tracking-widest mb-1">QUALITY GATE</p>
                <h2 className="text-2xl font-bold text-red-500 leading-tight">Failed</h2>
              </div>
              <div className="bg-[#161b22] px-4 py-2 rounded-lg border border-gray-800">
                <p className="text-white text-sm font-bold">12 bugs found</p>
              </div>
            </div>

            {/* Quick Metrics Bar */}
            <div className="bg-[#161b22] p-4 rounded-xl border border-gray-800 flex items-center gap-6">
              <div className="bg-[#0d1117] px-3 py-1.5 rounded text-xs font-bold text-[#5865f2] border border-gray-700">NEW CODE</div>
              <div className="text-sm font-medium">
                <span className="text-[#f1c40f]">505</span> <span className="text-gray-400">new issues</span>
              </div>
              <div className="text-sm font-medium border-l border-gray-700 pl-6">
                <span className="text-green-500">0</span> <span className="text-gray-400">closed</span>
              </div>
            </div>

            {/* Grid Metrics */}
            <div className="grid grid-cols-2 md:grid-cols-6 gap-4">
              {[
                { label: 'BUGS', value: '12', color: 'text-red-500', icon: <AlertCircle size={16} /> },
                { label: 'VULNERABILITIES', value: '0', color: 'text-green-500', icon: <ShieldAlert size={16} /> },
                { label: 'CODE SMELLS', value: '493', color: 'text-[#f1c40f]', icon: <Code2 size={16} /> },
                { label: 'COVERAGE', value: '—', color: 'text-gray-500', icon: <Layers size={16} /> },
                { label: 'DUPLICATION', value: '—', color: 'text-gray-500', icon: <FileCode2 size={16} /> },
                { label: 'LINES OF CODE', value: '23.1k', color: 'text-white', icon: null },
              ].map((item) => (
                <div key={item.label} className="bg-[#161b22] border border-gray-800 p-4 rounded-xl flex flex-col gap-2">
                  <div className={`text-2xl font-bold ${item.color}`}>{item.value}</div>
                  <div className="text-[10px] text-gray-500 font-bold tracking-wider">{item.label}</div>
                </div>
              ))}
            </div>

            {/* Charts Section */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Severity Distribution */}
              <div className="bg-[#161b22] border border-gray-800 p-6 rounded-xl space-y-4">
                <h4 className="text-xs font-bold text-gray-500 uppercase tracking-widest">SEVERITY DISTRIBUTION</h4>
                <div className="space-y-3">
                  {[
                    { label: 'Blocker', count: 0, width: '0%', color: 'bg-red-700' },
                    { label: 'Critical', count: 25, width: '8%', color: 'bg-orange-600' },
                    { label: 'Major', count: 94, width: '25%', color: 'bg-[#f1c40f]' },
                    { label: 'Minor', count: 383, width: '80%', color: 'bg-green-500' },
                    { label: 'Info', count: 3, width: '2%', color: 'bg-gray-500' },
                  ].map((row) => (
                    <div key={row.label} className="flex items-center gap-4 text-xs">
                      <span className="w-16 text-gray-400">{row.label}</span>
                      <div className="flex-grow h-3 bg-gray-800 rounded-full overflow-hidden">
                        <div className={`h-full ${row.color}`} style={{ width: row.width }}></div>
                      </div>
                      <span className="w-8 text-right font-bold">{row.count}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Type Distribution */}
              <div className="bg-[#161b22] border border-gray-800 p-6 rounded-xl space-y-4">
                <h4 className="text-xs font-bold text-gray-500 uppercase tracking-widest">TYPE DISTRIBUTION</h4>
                <div className="space-y-4">
                  {[
                    { label: 'Bug', count: 12, width: '5%', color: 'bg-red-500' },
                    { label: 'Code Smell', count: 493, width: '95%', color: 'bg-green-500' },
                    { label: 'Vulnerability', count: 0, width: '0%', color: 'bg-[#f1c40f]' },
                  ].map((row) => (
                    <div key={row.label} className="flex items-center gap-4 text-xs">
                      <span className="w-20 text-gray-400">{row.label}</span>
                      <div className="flex-grow h-3 bg-gray-800 rounded-full overflow-hidden">
                        <div className={`h-full ${row.color}`} style={{ width: row.width }}></div>
                      </div>
                      <span className="w-8 text-right font-bold">{row.count}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Hotspot Files */}
            <div className="bg-[#161b22] border border-gray-800 rounded-xl overflow-hidden">
              <div className="p-4 border-b border-gray-800">
                <h4 className="text-xs font-bold text-gray-500 uppercase tracking-widest">HOTSPOT FILES</h4>
              </div>
              <div className="divide-y divide-gray-800">
                {[
                  { path: 'ollantaengine/summarizer/cumsum_test.go', count: 32 },
                  { path: 'ollantaengine/newcode/resolver_test.go', count: 21 },
                  { path: 'api/static/app.js', count: 19 },
                ].map((file) => (
                  <div key={file.path} className="flex items-center justify-between p-4 text-sm hover:bg-gray-800/50 transition-colors cursor-default">
                    <span className="text-gray-300 font-mono text-xs">{file.path}</span>
                    <span className="text-[#f1c40f] font-bold">{file.count}</span>
                  </div>
                ))}
              </div>
            </div>
          </main>
        </div>
      )}
    </div>
  );
};

export default App;
