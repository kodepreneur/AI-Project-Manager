<template>
  <div class="space-y-6">
    <!-- Top Welcome & Status Banner -->
    <div class="flex flex-wrap items-center justify-between gap-4 p-6 bg-gradient-to-r from-dark-900 via-dark-850 to-dark-900 border border-slate-800 rounded-2xl relative overflow-hidden shadow-xl">
      <div class="space-y-1 relative z-10">
        <div class="flex items-center gap-2">
          <h2 class="text-xl font-bold text-white tracking-tight">AI Project Control Center</h2>
          <router-link 
            to="/antigravity" 
            class="flex items-center gap-1.5 text-[10px] font-mono px-2.5 py-0.5 rounded-full transition border"
            :class="cliInfo?.connected ? 'bg-emerald-950/60 text-emerald-400 border-emerald-800/60 hover:border-emerald-500' : 'bg-rose-950/60 text-rose-400 border-rose-800/60 hover:border-rose-500'"
          >
            <span class="w-1.5 h-1.5 rounded-full" :class="cliInfo?.connected ? 'bg-emerald-400 animate-pulse' : 'bg-rose-400'"></span>
            <span>{{ cliInfo?.connected ? `Antigravity CLI v${cliInfo.version}` : 'CLI Offline — Test Connection' }}</span>
          </router-link>
        </div>
        <p class="text-xs text-slate-400 max-w-xl">
          Create software projects, generate PRDs, orchestrate Antigravity AI coding agents, and monitor VPS infrastructure.
        </p>
      </div>

      <div class="flex items-center gap-2 relative z-10">
        <button 
          @click="showWizard = true"
          class="px-4 py-2 rounded-xl text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow-lg shadow-brand-600/20 transition flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>New Project</span>
        </button>
      </div>

      <!-- Subtle background decoration -->
      <div class="absolute -right-10 -top-10 w-48 h-48 bg-brand-500/10 rounded-full blur-2xl pointer-events-none"></div>
    </div>

    <!-- Metric Stat Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <!-- Total Projects -->
      <div class="glass-panel p-4 rounded-xl border border-slate-800 flex items-center justify-between">
        <div>
          <div class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Total Projects</div>
          <div class="text-2xl font-bold text-white mt-1">{{ overview?.projects?.total || 0 }}</div>
          <div class="text-[10px] text-slate-500 mt-1">
            {{ overview?.projects?.running || 0 }} running • {{ overview?.projects?.development || 0 }} in dev
          </div>
        </div>
        <div class="p-3 rounded-xl bg-blue-500/10 text-blue-400 border border-blue-500/20">
          <FolderKanban class="w-5 h-5" />
        </div>
      </div>

      <!-- Total TODOs -->
      <div class="glass-panel p-4 rounded-xl border border-slate-800 flex items-center justify-between">
        <div>
          <div class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">TODOs / Tasks</div>
          <div class="text-2xl font-bold text-white mt-1">{{ overview?.todos?.total || 0 }}</div>
          <div class="text-[10px] text-slate-500 mt-1">
            {{ overview?.todos?.completed || 0 }} done • {{ overview?.todos?.in_progress || 0 }} in progress
          </div>
        </div>
        <div class="p-3 rounded-xl bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
          <CheckCircle2 class="w-5 h-5" />
        </div>
      </div>

      <!-- Active Antigravity Sessions -->
      <div class="glass-panel p-4 rounded-xl border border-slate-800 flex items-center justify-between">
        <div>
          <div class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">Antigravity Sessions</div>
          <div class="text-2xl font-bold text-white mt-1">{{ overview?.antigravity?.active_sessions || 0 }}</div>
          <div class="text-[10px] text-purple-400 mt-1 flex items-center gap-1">
            <span class="w-1.5 h-1.5 rounded-full bg-purple-400 animate-pulse"></span>
            Managed Child Processes
          </div>
        </div>
        <div class="p-3 rounded-xl bg-purple-500/10 text-purple-400 border border-purple-500/20">
          <Bot class="w-5 h-5" />
        </div>
      </div>

      <!-- System Health -->
      <div class="glass-panel p-4 rounded-xl border border-slate-800 flex items-center justify-between">
        <div>
          <div class="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">System CPU & RAM</div>
          <div class="text-2xl font-bold text-white mt-1">{{ Math.round(metrics?.cpu_percent || 0) }}% / {{ Math.round(metrics?.mem_percent || 0) }}%</div>
          <div class="text-[10px] text-slate-500 mt-1 font-mono">
            {{ metrics?.os || 'linux' }} • Load: {{ metrics?.load_1?.toFixed(2) || '0.00' }}
          </div>
        </div>
        <div class="p-3 rounded-xl bg-amber-500/10 text-amber-400 border border-amber-500/20">
          <Cpu class="w-5 h-5" />
        </div>
      </div>
    </div>

    <!-- Middle Section: System Resource Telemetry & Active Projects -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- System Resource Bars -->
      <div class="bg-dark-900 border border-slate-800 rounded-2xl p-5 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <h3 class="text-sm font-bold text-white flex items-center gap-2">
            <Activity class="w-4 h-4 text-brand-400" />
            <span>Resource Telemetry</span>
          </h3>
          <span class="text-[10px] font-mono text-slate-500">{{ metrics?.server_time || '' }}</span>
        </div>

        <!-- CPU Usage -->
        <div class="space-y-1.5">
          <div class="flex justify-between text-xs text-slate-300 font-mono">
            <span>CPU Load ({{ metrics?.num_cpu || 1 }} Cores)</span>
            <span>{{ metrics?.cpu_percent?.toFixed(1) || 0 }}%</span>
          </div>
          <div class="w-full h-2 bg-dark-850 rounded-full overflow-hidden border border-slate-800">
            <div 
              class="h-full bg-brand-500 transition-all duration-500 rounded-full"
              :style="{ width: `${metrics?.cpu_percent || 0}%` }"
            ></div>
          </div>
        </div>

        <!-- RAM Usage -->
        <div class="space-y-1.5">
          <div class="flex justify-between text-xs text-slate-300 font-mono">
            <span>Memory ({{ metrics?.mem_used_mb || 0 }} / {{ metrics?.mem_total_mb || 0 }} MB)</span>
            <span>{{ metrics?.mem_percent?.toFixed(1) || 0 }}%</span>
          </div>
          <div class="w-full h-2 bg-dark-850 rounded-full overflow-hidden border border-slate-800">
            <div 
              class="h-full bg-indigo-500 transition-all duration-500 rounded-full"
              :style="{ width: `${metrics?.mem_percent || 0}%` }"
            ></div>
          </div>
        </div>

        <!-- Disk Usage -->
        <div class="space-y-1.5">
          <div class="flex justify-between text-xs text-slate-300 font-mono">
            <span>Disk Space ({{ metrics?.disk_used_gb || 0 }} / {{ metrics?.disk_total_gb || 0 }} GB)</span>
            <span>{{ metrics?.disk_percent?.toFixed(1) || 0 }}%</span>
          </div>
          <div class="w-full h-2 bg-dark-850 rounded-full overflow-hidden border border-slate-800">
            <div 
              class="h-full bg-emerald-500 transition-all duration-500 rounded-full"
              :style="{ width: `${metrics?.disk_percent || 0}%` }"
            ></div>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-2 pt-2 text-[11px] text-slate-400 font-mono">
          <div class="p-2.5 rounded-lg bg-dark-850 border border-slate-800">
            <div class="text-slate-500">Processes</div>
            <div class="text-white font-bold text-sm mt-0.5">{{ metrics?.process_count || 0 }} Active</div>
          </div>
          <div class="p-2.5 rounded-lg bg-dark-850 border border-slate-800">
            <div class="text-slate-500">1m Load Avg</div>
            <div class="text-white font-bold text-sm mt-0.5">{{ metrics?.load_1?.toFixed(2) || '0.00' }}</div>
          </div>
        </div>
      </div>

      <!-- Recent Project Activity Feed -->
      <div class="lg:col-span-2 bg-dark-900 border border-slate-800 rounded-2xl p-5 flex flex-col">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800 mb-4">
          <h3 class="text-sm font-bold text-white flex items-center gap-2">
            <ScrollText class="w-4 h-4 text-purple-400" />
            <span>Recent Activity & Development Log</span>
          </h3>
          <router-link to="/logs" class="text-xs text-brand-400 hover:text-brand-300 font-medium">
            View All Logs →
          </router-link>
        </div>

        <div class="flex-1 space-y-2 overflow-y-auto max-h-64 pr-1">
          <div 
            v-for="log in overview?.recent_activity" 
            :key="log.id"
            class="p-2.5 rounded-lg bg-dark-850 border border-slate-800/80 flex items-start justify-between gap-3 text-xs"
          >
            <div class="flex items-start gap-2.5">
              <span 
                :class="[
                  'text-[9px] uppercase font-bold font-mono px-1.5 py-0.5 rounded mt-0.5',
                  log.level === 'ERROR' ? 'bg-rose-950 text-rose-400' :
                  log.level === 'WARN' ? 'bg-amber-950 text-amber-400' : 'bg-dark-750 text-blue-400'
                ]"
              >
                {{ log.category }}
              </span>
              <span class="text-slate-200">{{ log.message }}</span>
            </div>
            <span class="text-[10px] text-slate-500 font-mono whitespace-nowrap">
              {{ new Date(log.created_at).toLocaleTimeString() }}
            </span>
          </div>

          <div v-if="!overview?.recent_activity?.length" class="text-center py-10 text-xs text-slate-500">
            No recent activity recorded.
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Section: Active Projects Quick Grid -->
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-bold text-white flex items-center gap-2">
          <FolderKanban class="w-4 h-4 text-brand-400" />
          <span>Active Projects</span>
        </h3>
        <router-link to="/projects" class="text-xs text-brand-400 hover:text-brand-300 font-medium">
          View All Projects ({{ projectStore.projects.length }}) →
        </router-link>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <ProjectCard 
          v-for="proj in projectStore.projects.slice(0, 3)" 
          :key="proj.id" 
          :project="proj"
          @open="openProject"
          @terminal="openTerminal"
          @preview="openPreview"
        />
      </div>

      <div v-if="projectStore.projects.length === 0" class="p-8 text-center bg-dark-900 border border-slate-800 rounded-2xl">
        <FolderKanban class="w-10 h-10 text-slate-600 mx-auto mb-2" />
        <h4 class="text-sm font-bold text-white">No projects created yet</h4>
        <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
          Start by launching the project creation wizard to initialize your first software project and PRD.
        </p>
        <button 
          @click="showWizard = true"
          class="mt-4 px-4 py-2 rounded-xl text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white transition"
        >
          Create Project
        </button>
      </div>
    </div>

    <!-- Project Wizard Modal -->
    <ProjectWizardModal :is-open="showWizard" @close="showWizard = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { useProjectStore } from '@/stores/projects';
import { useSystemStore } from '@/stores/system';
import { api } from '@/services/api';
import ProjectCard from '@/components/projects/ProjectCard.vue';
import ProjectWizardModal from '@/components/projects/ProjectWizardModal.vue';
import { 
  Plus, FolderKanban, CheckCircle2, Bot, Cpu, 
  Activity, ScrollText 
} from 'lucide-vue-next';

const router = useRouter();
const projectStore = useProjectStore();
const systemStore = useSystemStore();

const showWizard = ref(false);
const overview = ref<any>(null);
const metrics = ref<any>(null);
const cliInfo = ref<any>(null);

let pollInterval: any = null;

onMounted(async () => {
  await projectStore.fetchProjects();
  overview.value = await systemStore.fetchOverview();
  metrics.value = await systemStore.fetchMetrics();

  try {
    const res = await api.get('/antigravity/test');
    cliInfo.value = res.data.data;
  } catch {
    cliInfo.value = null;
  }

  // Periodic metrics refresh
  pollInterval = setInterval(async () => {
    metrics.value = await systemStore.fetchMetrics();
  }, 5000);
});

onBeforeUnmount(() => {
  if (pollInterval) clearInterval(pollInterval);
});

function openProject(id: number) {
  router.push(`/projects/${id}`);
}

function openTerminal(project: any) {
  router.push(`/projects/${project.id}?tab=terminal`);
}

function openPreview(project: any) {
  router.push(`/projects/${project.id}?tab=preview`);
}
</script>
