<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-white tracking-tight">Antigravity AI Sessions</h2>
        <p class="text-xs text-slate-400">Monitor active Antigravity CLI child processes, stream live execution logs, and review outputs.</p>
      </div>

      <div class="flex items-center gap-3">
        <button 
          @click="checkCliConnection"
          :disabled="checkingCli"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-purple-950/70 hover:bg-purple-900/90 border border-purple-800 text-purple-300 transition flex items-center gap-1.5 shadow"
        >
          <Activity class="w-3.5 h-3.5" :class="{ 'animate-spin': checkingCli }" />
          <span>{{ checkingCli ? 'Testing CLI...' : 'Test Connection' }}</span>
        </button>

        <div class="flex items-center gap-2">
          <label class="text-xs font-semibold text-slate-400">Project:</label>
          <select 
            v-model="selectedProjectId"
            class="px-3 py-1.5 bg-dark-900 border border-slate-800 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
            @change="loadSession"
          >
            <option v-for="p in projectStore.projects" :key="p.id" :value="p.id">
              {{ p.name }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Antigravity CLI Connectivity Status Bar -->
    <div 
      class="p-4 rounded-xl border flex flex-wrap items-center justify-between gap-3 text-xs transition"
      :class="cliInfo?.connected ? 'bg-emerald-950/20 border-emerald-800/40 text-emerald-200' : (cliChecked ? 'bg-rose-950/30 border-rose-800/50 text-rose-300' : 'bg-dark-900 border-slate-800 text-slate-300')"
    >
      <div class="flex items-center gap-2.5">
        <span 
          class="w-2.5 h-2.5 rounded-full"
          :class="cliInfo?.connected ? 'bg-emerald-400 animate-pulse' : (cliChecked ? 'bg-rose-500' : 'bg-slate-500')"
        ></span>
        <div>
          <span class="font-bold">Antigravity CLI: </span>
          <span v-if="cliInfo?.connected">
            Connected • {{ cliInfo.version }} ({{ cliInfo.arch }}) • Latency {{ cliInfo.latency_ms }}ms
          </span>
          <span v-else-if="cliChecked" class="text-rose-300">
            Disconnected or not found. Check <router-link to="/settings" class="underline font-semibold hover:text-white">Settings</router-link> to verify binary path.
          </span>
          <span v-else class="text-slate-400">Checking CLI connection...</span>
        </div>
      </div>

      <div v-if="cliInfo?.binary_path" class="font-mono text-[11px] text-slate-400">
        Binary: {{ cliInfo.binary_path }}
      </div>
    </div>

    <!-- Active Session Overview Card -->
    <div class="p-5 bg-dark-900 border border-purple-500/30 rounded-2xl flex flex-wrap items-center justify-between gap-4 shadow-xl">
      <div class="flex items-center gap-3">
        <div class="p-3 rounded-xl bg-purple-500/10 text-purple-400 border border-purple-500/20">
          <Bot class="w-6 h-6" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-sm text-white">Antigravity Autonomous Engine</h3>
            <span 
              :class="[
                'text-[10px] font-mono px-2 py-0.5 rounded-full border',
                session?.status === 'running' ? 'bg-purple-950 text-purple-400 border-purple-800' : 'bg-slate-800 text-slate-400 border-slate-700'
              ]"
            >
              ● {{ session?.status || 'idle' }}
            </span>
          </div>
          <div class="text-xs text-slate-400 mt-0.5 font-mono">
            PID: {{ session?.pid || '—' }} • Started: {{ session?.start_time ? new Date(session.start_time).toLocaleTimeString() : '—' }} • Task: {{ session?.current_task || 'None' }}
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button 
          v-if="session?.status === 'running'"
          @click="stopSession"
          class="px-3.5 py-1.5 rounded-lg text-xs font-semibold bg-rose-600 hover:bg-rose-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Square class="w-3.5 h-3.5" />
          <span>Terminate Task</span>
        </button>
      </div>
    </div>

    <!-- Live Execution Stream Window -->
    <div class="h-[520px] bg-[#090d16] border border-slate-800 rounded-2xl p-5 overflow-y-auto font-mono text-xs text-slate-300 leading-relaxed shadow-2xl space-y-1">
      <div v-for="(line, idx) in logs" :key="idx" class="flex items-start gap-2">
        <span class="text-purple-400 select-none">&gt;</span>
        <span class="text-slate-200">{{ line }}</span>
      </div>
      <div v-if="logs.length === 0" class="text-center py-24 text-slate-600">
        No active Antigravity output stream for this project. Dispatch a task from the Kanban board to view live progress.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useProjectStore } from '@/stores/projects';
import { api } from '@/services/api';
import { Bot, Square, Activity } from 'lucide-vue-next';

const projectStore = useProjectStore();
const selectedProjectId = ref<number | null>(null);
const session = ref<any>(null);
const logs = ref<string[]>([]);
const checkingCli = ref(false);
const cliChecked = ref(false);
const cliInfo = ref<any>(null);

onMounted(async () => {
  await projectStore.fetchProjects();
  if (projectStore.projects.length > 0) {
    selectedProjectId.value = projectStore.projects[0].id;
    loadSession();
  }
  checkCliConnection();
  setupWebSocket();
});

async function checkCliConnection() {
  checkingCli.value = true;
  try {
    const res = await api.get('/antigravity/test');
    cliInfo.value = res.data.data;
  } catch (err: any) {
    cliInfo.value = null;
  } finally {
    checkingCli.value = false;
    cliChecked.value = true;
  }
}

async function loadSession() {
  if (!selectedProjectId.value) return;
  const res = await api.get(`/projects/${selectedProjectId.value}/antigravity/status`);
  session.value = res.data.data;
  if (session.value?.stdout) {
    logs.value = session.value.stdout.split('\n').filter((l: string) => l.trim() !== '');
  } else {
    logs.value = [];
  }
}

async function stopSession() {
  if (!selectedProjectId.value) return;
  await api.post(`/projects/${selectedProjectId.value}/antigravity/stop`);
  logs.value.push('Session cancelled by user.');
  loadSession();
}

function setupWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const ws = new WebSocket(`${protocol}//${host}/ws/events`);

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.project_id && selectedProjectId.value && data.project_id !== selectedProjectId.value) {
        return;
      }
      if (data.type === 'antigravity_output') {
        logs.value.push(data.payload.line);
      } else if (data.type === 'antigravity_finished') {
        logs.value.push(`✔ Antigravity task finished (${data.payload.status}).`);
        loadSession();
      }
    } catch {}
  };
}
</script>
