<template>
  <div v-if="project" class="space-y-6">
    <!-- Project Top Header -->
    <div class="p-6 bg-dark-900 border border-slate-800 rounded-2xl flex flex-wrap items-center justify-between gap-4">
      <div class="space-y-1">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-bold text-white tracking-tight">{{ project.name }}</h2>
          <span 
            :class="[
              'text-[11px] font-semibold px-2.5 py-0.5 rounded-full border',
              project.status === 'Running' 
                ? 'bg-emerald-950/70 text-emerald-400 border-emerald-800' 
                : 'bg-blue-950/70 text-blue-400 border-blue-800'
            ]"
          >
            ● {{ project.status }}
          </span>
          <span class="text-xs font-mono text-slate-400 bg-dark-850 px-2 py-0.5 rounded border border-slate-800">
            {{ project.framework }} • {{ project.database_type }}
          </span>
        </div>
        <p class="text-xs text-slate-400 font-mono">{{ project.path }}</p>
      </div>

      <!-- Quick Header Actions -->
      <div class="flex items-center gap-2">
        <button 
          @click="runHealthCheck"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-dark-800 hover:bg-dark-750 text-slate-300 border border-slate-700 transition flex items-center gap-1.5"
        >
          <Activity class="w-3.5 h-3.5 text-emerald-400" />
          <span>Health Check</span>
        </button>

        <button 
          v-if="project.status !== 'Running'"
          @click="startDevServer"
          class="px-3.5 py-1.5 rounded-lg text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Play class="w-3.5 h-3.5" />
          <span>Start Server</span>
        </button>
      </div>
    </div>

    <!-- 12-Tab Navigation Bar -->
    <div class="border-b border-slate-800 flex items-center gap-1 overflow-x-auto pb-1 text-xs">
      <button 
        v-for="t in tabs" 
        :key="t.id"
        :class="[
          'px-3.5 py-2 rounded-lg font-medium whitespace-nowrap transition flex items-center gap-1.5',
          activeTab === t.id ? 'bg-brand-600 text-white shadow-md' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'
        ]"
        @click="switchTab(t.id)"
      >
        <component :is="t.icon" class="w-3.5 h-3.5" />
        <span>{{ t.label }}</span>
      </button>
    </div>

    <!-- TAB 1: OVERVIEW -->
    <div v-if="activeTab === 'overview'" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-6">
        <!-- Key Metrics Cards -->
        <div class="grid grid-cols-3 gap-4">
          <div class="glass-panel p-4 rounded-xl border border-slate-800">
            <div class="text-[11px] font-semibold text-slate-400">TODO PROGRESS</div>
            <div class="text-xl font-bold text-white mt-1">{{ todoPercent }}%</div>
            <div class="text-[10px] text-slate-500 mt-1">{{ completedCount }} of {{ totalCount }} completed</div>
          </div>
          <div class="glass-panel p-4 rounded-xl border border-slate-800">
            <div class="text-[11px] font-semibold text-slate-400">PRD STATUS</div>
            <div class="text-xl font-bold text-white mt-1">{{ project.prd?.status || 'Draft' }}</div>
            <div class="text-[10px] text-brand-400 mt-1 font-mono">v{{ project.prd?.version || 1 }}</div>
          </div>
          <div class="glass-panel p-4 rounded-xl border border-slate-800">
            <div class="text-[11px] font-semibold text-slate-400">DEV PREVIEW</div>
            <div class="text-xl font-bold text-white mt-1">Port :{{ project.preview_port }}</div>
            <div class="text-[10px] text-emerald-400 mt-1 font-mono">http://localhost:{{ project.preview_port }}</div>
          </div>
        </div>

        <!-- Description & Information -->
        <div class="bg-dark-900 border border-slate-800 rounded-xl p-5 space-y-3">
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400">Project Specifications</h3>
          <p class="text-xs text-slate-300 leading-relaxed">{{ project.description || 'No description provided.' }}</p>
          
          <div class="grid grid-cols-2 gap-3 pt-3 border-t border-slate-800 text-xs font-mono">
            <div><span class="text-slate-500">Framework:</span> <span class="text-slate-200">{{ project.framework }}</span></div>
            <div><span class="text-slate-500">Frontend:</span> <span class="text-slate-200">{{ project.frontend_tech }}</span></div>
            <div><span class="text-slate-500">Database:</span> <span class="text-slate-200">{{ project.database_type }}</span></div>
            <div><span class="text-slate-500">PHP Runtime:</span> <span class="text-slate-200">{{ project.php_version }}</span></div>
            <div><span class="text-slate-500">Git Branch:</span> <span class="text-slate-200">{{ project.git_branch || 'main' }}</span></div>
            <div><span class="text-slate-500">Last Activity:</span> <span class="text-slate-200">{{ new Date(project.last_activity).toLocaleString() }}</span></div>
          </div>
        </div>
      </div>

      <!-- Overview Sidebar: Preview Snapshot & Quick Launch -->
      <div class="space-y-4">
        <div class="bg-dark-900 border border-slate-800 rounded-xl p-4">
          <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Preview Snapshot</h4>
          <div class="w-full h-40 rounded-lg bg-dark-950 border border-slate-800 overflow-hidden mb-3 relative">
            <img 
              v-if="project.screenshot_url" 
              :src="project.screenshot_url" 
              alt="Snapshot" 
              class="w-full h-full object-cover" 
            />
            <div v-else class="w-full h-full flex flex-col items-center justify-center text-slate-600 text-xs">
              <Globe class="w-6 h-6 mb-1" />
              <span>Standby</span>
            </div>
          </div>
          <button 
            @click="switchTab('preview')" 
            class="w-full py-2 rounded-lg text-xs font-semibold bg-dark-800 hover:bg-dark-750 text-slate-200 border border-slate-700 transition"
          >
            Open Interactive Preview
          </button>
        </div>
      </div>
    </div>

    <!-- TAB 2: PRD -->
    <div v-if="activeTab === 'prd'">
      <PrdEditor :project-id="project.id" @approved="refreshProject" />
    </div>

    <!-- TAB 3: TASKS (KANBAN) -->
    <div v-if="activeTab === 'tasks'">
      <KanbanBoard :project-id="project.id" @start-ai="onStartAI" />
    </div>

    <!-- TAB 4: ANTIGRAVITY AI -->
    <div v-if="activeTab === 'antigravity'" class="space-y-4">
      <div class="flex flex-wrap items-center justify-between gap-3 p-4 bg-dark-900 border border-purple-500/30 rounded-xl">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg bg-purple-500/10 text-purple-400 border border-purple-500/20">
            <Bot class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h3 class="font-bold text-sm text-white">Antigravity AI Session Runner</h3>
              <span 
                :class="[
                  'text-[10px] font-mono px-2 py-0.5 rounded-full border',
                  aiSession?.status === 'running' 
                    ? 'bg-purple-950 text-purple-400 border-purple-800'
                    : 'bg-slate-800 text-slate-400 border-slate-700'
                ]"
              >
                ● {{ aiSession?.status || 'idle' }}
              </span>
            </div>
            <p class="text-xs text-slate-400 font-mono">
              Current Task: {{ aiSession?.current_task || 'None' }}
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button 
            type="button"
            @click="testCli"
            :disabled="testingCli"
            class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-purple-950/70 hover:bg-purple-900 border border-purple-800 text-purple-300 transition flex items-center gap-1.5 shadow"
          >
            <Activity class="w-3.5 h-3.5" :class="{ 'animate-spin': testingCli }" />
            <span>{{ testingCli ? 'Testing...' : (cliInfo ? (cliInfo.connected ? '✔ CLI Online' : '✖ CLI Offline') : 'Test CLI') }}</span>
          </button>

          <button 
            v-if="aiSession?.status === 'running'"
            @click="stopAiSession"
            class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-rose-600 hover:bg-rose-500 text-white shadow transition flex items-center gap-1"
          >
            <Square class="w-3.5 h-3.5" />
            <span>Stop Session</span>
          </button>
        </div>
      </div>

      <!-- CLI status alert if test performed -->
      <div v-if="cliInfo" class="p-3 rounded-lg text-xs border" :class="cliInfo.connected ? 'bg-emerald-950/20 border-emerald-800/40 text-emerald-300' : 'bg-rose-950/30 border-rose-800 text-rose-300'">
        <div class="flex items-center justify-between font-semibold">
          <span>{{ cliInfo.connected ? `Antigravity CLI v${cliInfo.version} (${cliInfo.arch}) Ready` : 'Antigravity CLI Offline' }}</span>
          <span v-if="cliInfo.connected" class="font-mono text-[10px] text-emerald-400">{{ cliInfo.latency_ms }}ms latency</span>
        </div>
        <div class="text-[11px] font-mono text-slate-400 mt-0.5">
          {{ cliInfo.connected ? `Binary: ${cliInfo.binary_path}` : 'Please configure or verify binary path in Settings.' }}
        </div>
      </div>

      <!-- Live Terminal Output Stream -->
      <div class="h-[480px] bg-[#090d16] border border-slate-800 rounded-xl p-4 overflow-y-auto font-mono text-xs text-slate-300 leading-relaxed shadow-2xl">
        <div v-for="(line, idx) in aiOutputLines" :key="idx" class="py-0.5">
          <span class="text-purple-400 select-none mr-2">&gt;</span>
          <span>{{ line }}</span>
        </div>
        <div v-if="aiOutputLines.length === 0" class="text-center py-20 text-slate-600">
          No active session output. Start a task from the Kanban board to trigger Antigravity.
        </div>
      </div>
    </div>

    <!-- TAB 5: TERMINAL -->
    <div v-if="activeTab === 'terminal'">
      <WebTerminal :project-id="project.id" :project-path="project.path" />
    </div>

    <!-- TAB 6: PREVIEW -->
    <div v-if="activeTab === 'preview'">
      <PreviewPanel :project="project" />
    </div>

    <!-- TAB 7: GIT -->
    <div v-if="activeTab === 'git'">
      <GitPanel :project-id="project.id" />
    </div>

    <!-- TAB 8: ENVIRONMENT -->
    <div v-if="activeTab === 'env'" class="space-y-4">
      <div class="flex items-center justify-between p-4 bg-dark-900 border border-slate-800 rounded-xl">
        <div>
          <h3 class="text-sm font-bold text-white">Environment Variables (.env)</h3>
          <p class="text-xs text-slate-400">Configure application variables with automatic secret masking.</p>
        </div>
        <button 
          @click="saveEnv"
          :disabled="savingEnv"
          class="px-4 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Save class="w-3.5 h-3.5" />
          <span>{{ savingEnv ? 'Saving...' : 'Save .env' }}</span>
        </button>
      </div>

      <div class="bg-dark-900 border border-slate-800 rounded-xl p-4 space-y-2">
        <div v-for="(item, idx) in envItems" :key="idx" class="flex items-center gap-2">
          <input 
            v-if="!item.is_comment"
            v-model="item.key"
            type="text" 
            placeholder="KEY"
            class="w-1/3 px-3 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-brand-300 font-mono"
          />
          <input 
            v-if="!item.is_comment"
            v-model="item.value"
            :type="item.is_secret && !item.show ? 'password' : 'text'"
            placeholder="VALUE"
            class="flex-1 px-3 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
          />
          <span v-else class="flex-1 text-xs text-slate-500 font-mono py-1 px-3">{{ item.value }}</span>

          <button 
            v-if="item.is_secret" 
            @click="item.show = !item.show"
            class="text-xs text-slate-400 hover:text-white px-2 py-1"
          >
            {{ item.show ? 'Hide' : 'Show' }}
          </button>
          <button @click="removeEnvItem(idx)" class="text-rose-400 hover:text-rose-300 px-2 py-1">✕</button>
        </div>

        <button 
          @click="addEnvItem"
          class="mt-2 text-xs font-semibold text-brand-400 hover:text-brand-300 flex items-center gap-1"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Add Variable</span>
        </button>
      </div>
    </div>

    <!-- TAB 9: DATABASE -->
    <div v-if="activeTab === 'database'">
      <DatabaseManager :engine="project.database_type" :project-id="project.id" />
    </div>

    <!-- TAB 10: PROCESSES -->
    <div v-if="activeTab === 'processes'" class="space-y-4">
      <div class="flex items-center justify-between p-4 bg-dark-900 border border-slate-800 rounded-xl">
        <div>
          <h3 class="text-sm font-bold text-white">Project Processes & Background Daemons</h3>
          <p class="text-xs text-slate-400">Launch and track dev servers, workers, and build scripts.</p>
        </div>
        <button 
          @click="openLaunchProcess"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Launch Process</span>
        </button>
      </div>

      <div class="bg-dark-900 border border-slate-800 rounded-xl overflow-hidden">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-dark-950/60 border-b border-slate-800 text-slate-400">
              <th class="p-3">Process Name</th>
              <th class="p-3">Command</th>
              <th class="p-3">PID</th>
              <th class="p-3">Port</th>
              <th class="p-3">Status</th>
              <th class="p-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="proc in processes" :key="proc.id" class="border-b border-slate-850">
              <td class="p-3 font-semibold text-white">{{ proc.name }}</td>
              <td class="p-3 font-mono text-slate-300">{{ proc.command }}</td>
              <td class="p-3 font-mono text-slate-400">{{ proc.pid || '—' }}</td>
              <td class="p-3 font-mono text-slate-400">{{ proc.port || '—' }}</td>
              <td class="p-3">
                <span 
                  :class="[
                    'px-2 py-0.5 rounded-full text-[10px] font-mono border',
                    proc.status === 'running' ? 'bg-emerald-950 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-400 border-slate-700'
                  ]"
                >
                  ● {{ proc.status }}
                </span>
              </td>
              <td class="p-3 text-right">
                <button 
                  v-if="proc.status === 'running'"
                  @click="stopProcess(proc.id)"
                  class="px-2 py-1 rounded bg-rose-950 text-rose-300 hover:bg-rose-900 transition text-[11px]"
                >
                  Stop
                </button>
              </td>
            </tr>
            <tr v-if="processes.length === 0">
              <td colspan="6" class="text-center py-10 text-slate-500">No processes currently managed.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 11: LOGS -->
    <div v-if="activeTab === 'logs'" class="space-y-4">
      <div class="h-[500px] bg-[#090d16] border border-slate-800 rounded-xl p-4 overflow-y-auto font-mono text-xs text-slate-300 leading-relaxed shadow-2xl">
        <div v-for="l in projectLogs" :key="l.id" class="py-1 border-b border-slate-900 flex items-start gap-3">
          <span class="text-slate-500 select-none">{{ new Date(l.created_at).toLocaleTimeString() }}</span>
          <span 
            :class="[
              'px-1.5 py-0.2 rounded text-[10px] font-bold uppercase',
              l.level === 'ERROR' ? 'bg-rose-950 text-rose-400' :
              l.level === 'WARN' ? 'bg-amber-950 text-amber-400' : 'bg-dark-800 text-blue-400'
            ]"
          >
            {{ l.category }}
          </span>
          <span class="text-slate-200">{{ l.message }}</span>
        </div>
        <div v-if="projectLogs.length === 0" class="text-center py-20 text-slate-600">
          No logs captured for this project yet.
        </div>
      </div>
    </div>

    <!-- TAB 12: SETTINGS & HEALTH CHECK -->
    <div v-if="activeTab === 'settings'" class="bg-dark-900 border border-slate-800 rounded-xl p-6 space-y-4 max-w-2xl">
      <h3 class="text-sm font-bold text-white">Project Settings</h3>

      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1">Preview Port</label>
        <input 
          v-model="project.preview_port"
          type="number" 
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
        />
      </div>

      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1">Preview Command</label>
        <input 
          v-model="project.preview_command"
          type="text" 
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
        />
      </div>

      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1">Test Command</label>
        <input 
          v-model="project.test_command"
          type="text" 
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
        />
      </div>

      <div class="pt-4 flex items-center justify-between border-t border-slate-800">
        <button 
          @click="archiveProject"
          class="text-xs font-semibold text-amber-400 hover:text-amber-300"
        >
          Archive Project
        </button>

        <button 
          @click="saveSettings"
          class="px-4 py-2 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow"
        >
          Update Settings
        </button>
      </div>
    </div>

    <!-- Health Check Modal -->
    <div 
      v-if="showHealthModal" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
      @click.self="showHealthModal = false"
    >
      <div class="w-full max-w-md bg-dark-900 border border-slate-700 rounded-2xl shadow-2xl overflow-hidden animate-scale-up">
        <div class="px-6 py-4 border-b border-slate-800 bg-dark-950/40 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Activity class="w-5 h-5 text-emerald-400" />
            <h3 class="text-sm font-bold text-white">Project Health Diagnostics</h3>
          </div>
          <button @click="showHealthModal = false" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="p-6 space-y-3">
          <div 
            v-for="(h, idx) in healthResults" 
            :key="idx"
            class="p-3 bg-dark-850 border border-slate-800 rounded-lg flex items-center justify-between text-xs"
          >
            <div>
              <div class="font-semibold text-white">{{ h.name }}</div>
              <div class="text-[10px] text-slate-400">{{ h.description }}</div>
            </div>
            <span 
              :class="[
                'px-2 py-0.5 rounded text-[10px] font-bold',
                h.passed ? 'bg-emerald-950 text-emerald-400 border border-emerald-800' : 'bg-rose-950 text-rose-400 border border-rose-800'
              ]"
            >
              {{ h.passed ? 'PASSED' : 'CHECK' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProjectStore } from '@/stores/projects';
import { api } from '@/services/api';
import type { Project } from '@/types';
import PrdEditor from '@/components/prd/PrdEditor.vue';
import KanbanBoard from '@/components/kanban/KanbanBoard.vue';
import WebTerminal from '@/components/terminal/WebTerminal.vue';
import PreviewPanel from '@/components/preview/PreviewPanel.vue';
import GitPanel from '@/components/git/GitPanel.vue';
import DatabaseManager from '@/components/database/DatabaseManager.vue';
import { 
  FolderKanban, FileText, Columns3, Bot, Terminal, 
  Eye, GitBranch, Key, Database, Activity, ScrollText, 
  Settings, Play, Globe, Square, Save, Plus 
} from 'lucide-vue-next';

const route = useRoute();
const router = useRouter();
const projectStore = useProjectStore();

const project = ref<Project | null>(null);
const activeTab = ref('overview');

const aiSession = ref<any>(null);
const aiOutputLines = ref<string[]>([]);
const testingCli = ref(false);
const cliInfo = ref<any>(null);
const processes = ref<any[]>([]);
const envItems = ref<any[]>([]);
const savingEnv = ref(false);
const projectLogs = ref<any[]>([]);
const showHealthModal = ref(false);
const healthResults = ref<any[]>([]);

const tabs = [
  { id: 'overview', label: 'Overview', icon: FolderKanban },
  { id: 'prd', label: 'PRD', icon: FileText },
  { id: 'tasks', label: 'Tasks', icon: Columns3 },
  { id: 'antigravity', label: 'Antigravity', icon: Bot },
  { id: 'terminal', label: 'Terminal', icon: Terminal },
  { id: 'preview', label: 'Preview', icon: Eye },
  { id: 'git', label: 'Git', icon: GitBranch },
  { id: 'env', label: 'Environment', icon: Key },
  { id: 'database', label: 'Database', icon: Database },
  { id: 'processes', label: 'Processes', icon: Activity },
  { id: 'logs', label: 'Logs', icon: ScrollText },
  { id: 'settings', label: 'Settings', icon: Settings },
];

onMounted(async () => {
  if (route.query.tab) {
    activeTab.value = route.query.tab as string;
  }
  await loadProjectData();
  setupWebSocketEvents();
});

watch(() => route.params.id, () => {
  loadProjectData();
});

async function loadProjectData() {
  const p = await projectStore.fetchProject(route.params.id as string);
  project.value = p;
  if (p) {
    fetchProcesses();
    fetchEnv();
    fetchLogs();
    fetchAiStatus();
  }
}

function switchTab(tabId: string) {
  activeTab.value = tabId;
  router.push({ query: { tab: tabId } });
}

const totalCount = computed(() => project.value?.tasks?.length || 0);
const completedCount = computed(() => {
  if (!project.value?.tasks) return 0;
  return project.value.tasks.filter(t => t.status === 'done').length;
});
const todoPercent = computed(() => {
  if (totalCount.value === 0) return 0;
  return Math.round((completedCount.value / totalCount.value) * 100);
});

async function fetchProcesses() {
  const res = await api.get(`/projects/${project.value?.id}/processes`);
  processes.value = res.data.data;
}

async function fetchEnv() {
  const res = await api.get(`/projects/${project.value?.id}/env`);
  envItems.value = res.data.data;
}

async function saveEnv() {
  savingEnv.value = true;
  try {
    await api.post(`/projects/${project.value?.id}/env`, { items: envItems.value });
    alert('Environment configuration saved.');
  } finally {
    savingEnv.value = false;
  }
}

function addEnvItem() {
  envItems.value.push({ key: '', value: '', is_secret: false });
}

function removeEnvItem(idx: number) {
  envItems.value.splice(idx, 1);
}

async function fetchLogs() {
  const res = await api.get(`/logs?project_id=${project.value?.id}`);
  projectLogs.value = res.data.data;
}

async function fetchAiStatus() {
  const res = await api.get(`/projects/${project.value?.id}/antigravity/status`);
  aiSession.value = res.data.data;
}

function onStartAI(task: any) {
  switchTab('antigravity');
  aiOutputLines.value = [`Launching task context for '${task.title}'...`];
  fetchAiStatus();
}

async function stopAiSession() {
  await api.post(`/projects/${project.value?.id}/antigravity/stop`);
  aiOutputLines.value.push('Session cancelled by user.');
  fetchAiStatus();
}

async function testCli() {
  testingCli.value = true;
  cliInfo.value = null;
  try {
    const res = await api.post(`/projects/${project.value?.id}/antigravity/test`);
    cliInfo.value = res.data.data;
  } catch (err: any) {
    cliInfo.value = {
      connected: false,
      error: err.response?.data?.error?.message || 'Failed to connect to Antigravity CLI'
    };
  } finally {
    testingCli.value = false;
  }
}

async function startDevServer() {
  await api.post(`/projects/${project.value?.id}/processes/start`, {
    name: 'Dev Server',
    command: project.value?.preview_command || 'npm run dev',
    port: project.value?.preview_port || 8001,
  });
  await loadProjectData();
}

async function stopProcess(id: number) {
  await api.post(`/processes/${id}/stop`);
  await fetchProcesses();
}

function openLaunchProcess() {
  const cmd = prompt('Enter command to run in project directory:', 'npm run dev');
  if (cmd) {
    api.post(`/projects/${project.value?.id}/processes/start`, {
      name: 'Custom Process',
      command: cmd,
      port: 0,
    }).then(() => fetchProcesses());
  }
}

async function runHealthCheck() {
  const res = await api.get(`/projects/${project.value?.id}/health`);
  healthResults.value = res.data.data;
  showHealthModal.value = true;
}

async function saveSettings() {
  await projectStore.updateProject(project.value!.id, {
    preview_port: project.value!.preview_port,
    preview_command: project.value!.preview_command,
    test_command: project.value!.test_command,
  });
  alert('Settings saved.');
}

async function archiveProject() {
  if (confirm('Archive this project?')) {
    await projectStore.archiveProject(project.value!.id);
    router.push('/projects');
  }
}

function refreshProject() {
  loadProjectData();
}

function setupWebSocketEvents() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const host = window.location.host;
  const ws = new WebSocket(`${protocol}//${host}/ws/events`);

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.project_id && project.value && data.project_id !== project.value.id) {
        return;
      }
      if (data.type === 'antigravity_output') {
        aiOutputLines.value.push(data.payload.line);
      } else if (data.type === 'antigravity_finished') {
        aiOutputLines.value.push(`✔ Task finished with status: ${data.payload.status}`);
        fetchAiStatus();
        projectStore.fetchProjects();
      } else if (data.type === 'log') {
        projectLogs.value.unshift(data.payload);
      }
    } catch {}
  };
}
</script>
