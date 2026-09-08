<template>
  <div class="flex flex-col h-[520px] bg-dark-950 border border-slate-800 rounded-xl overflow-hidden shadow-2xl">
    <!-- Terminal Header Bar -->
    <div class="h-10 px-4 bg-dark-900 border-b border-slate-800 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <!-- Traffic light indicator -->
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-rose-500/80"></span>
          <span class="w-2.5 h-2.5 rounded-full bg-amber-500/80"></span>
          <span class="w-2.5 h-2.5 rounded-full bg-emerald-500/80"></span>
        </div>

        <div class="flex items-center gap-2 text-xs font-mono text-slate-300">
          <Terminal class="w-3.5 h-3.5 text-brand-400" />
          <span>{{ projectPath || 'Terminal' }}</span>
        </div>

        <!-- Connection Status Badge -->
        <span 
          :class="[
            'text-[10px] font-mono px-2 py-0.5 rounded-full border',
            connected 
              ? 'bg-emerald-950/60 text-emerald-400 border-emerald-800/60' 
              : 'bg-rose-950/60 text-rose-400 border-rose-800/60'
          ]"
        >
          {{ connected ? 'Connected' : 'Disconnected' }}
        </span>
      </div>

      <div class="flex items-center gap-2">
        <button 
          @click="reconnect"
          class="p-1 rounded text-slate-400 hover:text-white hover:bg-slate-800 transition"
          title="Reconnect Terminal"
        >
          <RotateCcw class="w-3.5 h-3.5" />
        </button>
        <button 
          @click="clear"
          class="p-1 rounded text-slate-400 hover:text-white hover:bg-slate-800 transition"
          title="Clear Terminal"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>

    <!-- xterm.js DOM Mount Element -->
    <div ref="terminalContainer" class="flex-1 p-2 bg-[#090d16] overflow-hidden"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import { Terminal, RotateCcw, Trash2 } from 'lucide-vue-next';

const props = defineProps<{
  projectId: number | string;
  projectPath?: string;
}>();

const terminalContainer = ref<HTMLDivElement | null>(null);
const connected = ref(false);

let term: XTerm | null = null;
let fitAddon: FitAddon | null = null;
let socket: WebSocket | null = null;

onMounted(() => {
  initTerminal();
  connectWebSocket();
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
  closeConnection();
  if (term) term.dispose();
});

function initTerminal() {
  if (!terminalContainer.value) return;

  term = new XTerm({
    cursorBlink: true,
    fontFamily: 'JetBrains Mono, Menlo, Monaco, Courier New, monospace',
    fontSize: 12,
    lineHeight: 1.3,
    theme: {
      background: '#090d16',
      foreground: '#cbd5e1',
      cursor: '#3b82f6',
      selectionBackground: '#1d4ed840',
      black: '#0f172a',
      red: '#ef4444',
      green: '#10b981',
      yellow: '#f59e0b',
      blue: '#3b82f6',
      magenta: '#a855f7',
      cyan: '#06b6d4',
      white: '#f8fafc',
    },
  });

  fitAddon = new FitAddon();
  term.loadAddon(fitAddon);
  term.loadAddon(new WebLinksAddon());

  term.open(terminalContainer.value);
  fitAddon.fit();

  term.onData((data) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(data);
    }
  });
}

function connectWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const token = localStorage.getItem('token') || '';
  const url = `${protocol}//${host}/ws/projects/${props.projectId}/terminal?token=${token}`;

  socket = new WebSocket(url);
  socket.binaryType = 'arraybuffer';

  socket.onopen = () => {
    connected.value = true;
    term?.writeln('\x1b[32m✔ Project terminal session established.\x1b[0m\r\n');
    sendResize();
  };

  socket.onmessage = (event) => {
    if (event.data instanceof ArrayBuffer) {
      term?.write(new Uint8Array(event.data));
    } else {
      term?.write(event.data);
    }
  };

  socket.onclose = () => {
    connected.value = false;
    term?.writeln('\r\n\x1b[31m✖ Terminal session disconnected.\x1b[0m\r\n');
  };

  socket.onerror = () => {
    connected.value = false;
  };
}

function sendResize() {
  if (term && fitAddon && socket && socket.readyState === WebSocket.OPEN) {
    fitAddon.fit();
    socket.send(JSON.stringify({
      type: 'resize',
      cols: term.cols,
      rows: term.rows,
    }));
  }
}

function handleResize() {
  sendResize();
}

function reconnect() {
  closeConnection();
  term?.reset();
  connectWebSocket();
}

function clear() {
  term?.clear();
}

function closeConnection() {
  if (socket) {
    socket.close();
    socket = null;
  }
}
</script>
