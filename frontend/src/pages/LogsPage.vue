<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-white tracking-tight">Centralized Logging</h2>
        <p class="text-xs text-slate-400">Real-time structured logs streamed from Antigravity sessions, child processes, and system monitors.</p>
      </div>

      <div class="flex items-center gap-2">
        <button 
          @click="downloadLogs"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-dark-800 hover:bg-dark-750 text-slate-300 border border-slate-700 transition flex items-center gap-1.5"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Export Logs</span>
        </button>

        <button 
          @click="clearLogs"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-rose-950/60 hover:bg-rose-900 text-rose-300 border border-rose-800/60 transition flex items-center gap-1.5"
        >
          <Trash2 class="w-3.5 h-3.5" />
          <span>Clear Logs</span>
        </button>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="p-4 bg-dark-900 border border-slate-800 rounded-xl flex flex-wrap items-center gap-3 text-xs">
      <div class="flex-1 min-w-[200px]">
        <input 
          v-model="search"
          type="text" 
          placeholder="Filter logs by keyword..."
          class="w-full px-3 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-brand-500"
          @input="fetchLogs"
        />
      </div>

      <select 
        v-model="category"
        class="px-3 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-white"
        @change="fetchLogs"
      >
        <option value="">All Categories</option>
        <option value="Application">Application</option>
        <option value="Antigravity">Antigravity</option>
        <option value="Terminal">Terminal</option>
        <option value="PHP">PHP</option>
        <option value="Database">Database</option>
        <option value="Project Process">Project Process</option>
        <option value="System">System</option>
      </select>

      <select 
        v-model="level"
        class="px-3 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-white"
        @change="fetchLogs"
      >
        <option value="">All Levels</option>
        <option value="INFO">INFO</option>
        <option value="WARN">WARN</option>
        <option value="ERROR">ERROR</option>
      </select>
    </div>

    <!-- Logs Stream Window -->
    <div class="h-[550px] bg-[#090d16] border border-slate-800 rounded-2xl p-4 overflow-y-auto font-mono text-xs text-slate-300 leading-relaxed shadow-2xl space-y-1">
      <div 
        v-for="l in logs" 
        :key="l.id"
        class="py-1 px-2 rounded hover:bg-slate-900/50 flex items-start gap-3 transition"
      >
        <span class="text-slate-500 select-none whitespace-nowrap">{{ new Date(l.created_at).toLocaleTimeString() }}</span>
        <span 
          :class="[
            'px-1.5 py-0.2 rounded text-[10px] font-bold uppercase whitespace-nowrap',
            l.level === 'ERROR' ? 'bg-rose-950 text-rose-400 border border-rose-800/60' :
            l.level === 'WARN' ? 'bg-amber-950 text-amber-400 border border-amber-800/60' : 'bg-dark-800 text-blue-400'
          ]"
        >
          {{ l.level }}
        </span>
        <span class="text-purple-400 select-none whitespace-nowrap">[{{ l.category }}]</span>
        <span class="text-slate-200 flex-1 break-all">{{ l.message }}</span>
      </div>

      <div v-if="logs.length === 0" class="text-center py-24 text-slate-600">
        No log entries match the specified filters.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { api } from '@/services/api';
import type { ProjectLog } from '@/types';
import { Download, Trash2 } from 'lucide-vue-next';

const logs = ref<ProjectLog[]>([]);
const search = ref('');
const category = ref('');
const level = ref('');

onMounted(() => {
  fetchLogs();
  setupWebSocket();
});

async function fetchLogs() {
  const params: Record<string, string> = {};
  if (search.value) params.search = search.value;
  if (category.value) params.category = category.value;
  if (level.value) params.level = level.value;

  const res = await api.get('/logs', { params });
  logs.value = res.data.data;
}

function downloadLogs() {
  window.open('/api/logs/download', '_blank');
}

async function clearLogs() {
  if (confirm('Clear all logs?')) {
    await api.delete('/logs');
    logs.value = [];
  }
}

function setupWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const ws = new WebSocket(`${protocol}//${host}/ws/events`);

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.type === 'log') {
        logs.value.unshift(data.payload);
      }
    } catch {}
  };
}
</script>
