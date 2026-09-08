<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-white tracking-tight">Infrastructure & Runtimes</h2>
        <p class="text-xs text-slate-400">Inspect host PHP binaries, MySQL/PostgreSQL databases, and manage background daemons.</p>
      </div>

      <!-- Infra Subtabs -->
      <div class="flex items-center gap-1 bg-dark-900 border border-slate-800 p-1 rounded-xl text-xs">
        <button 
          v-for="st in subtabs" 
          :key="st.id"
          :class="[
            'px-3 py-1.5 rounded-lg font-medium transition flex items-center gap-1.5',
            activeSubtab === st.id ? 'bg-brand-600 text-white shadow' : 'text-slate-400 hover:text-slate-200'
          ]"
          @click="activeSubtab = st.id"
        >
          <component :is="st.icon" class="w-3.5 h-3.5" />
          <span>{{ st.label }}</span>
        </button>
      </div>
    </div>

    <!-- SUBTAB 1: PHP RUNTIME -->
    <div v-if="activeSubtab === 'php'" class="space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div 
          v-for="ver in phpData?.versions" 
          :key="ver.version"
          class="glass-panel p-4 rounded-xl border border-slate-800"
        >
          <div class="flex items-center justify-between">
            <span class="text-sm font-bold text-white">PHP {{ ver.version }}</span>
            <span 
              :class="[
                'text-[10px] font-semibold px-2 py-0.5 rounded-full border',
                ver.installed ? 'bg-emerald-950/70 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-500 border-slate-700'
              ]"
            >
              {{ ver.installed ? 'Installed' : 'Not Detected' }}
            </span>
          </div>
          <div class="text-[11px] text-slate-400 font-mono mt-2 truncate">
            {{ ver.path || 'Unavailable' }}
          </div>
          <div v-if="ver.is_default" class="mt-2 text-[10px] text-brand-400 font-semibold uppercase">
            ★ Default CLI Binary
          </div>
        </div>
      </div>

      <!-- Loaded PHP Modules -->
      <div class="bg-dark-900 border border-slate-800 rounded-2xl p-5 space-y-3">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center justify-between">
          <span>Loaded PHP Extensions & Modules</span>
          <span class="font-mono text-slate-500">{{ phpData?.modules?.length || 0 }} Active</span>
        </h3>
        <div class="flex flex-wrap gap-1.5 max-h-48 overflow-y-auto">
          <span 
            v-for="mod in phpData?.modules" 
            :key="mod"
            class="px-2 py-0.5 rounded bg-dark-850 border border-slate-800 text-[11px] font-mono text-slate-300"
          >
            {{ mod }}
          </span>
          <span v-if="!phpData?.modules?.length" class="text-xs text-slate-500">
            No PHP extensions detected.
          </span>
        </div>
      </div>
    </div>

    <!-- SUBTAB 2: MYSQL -->
    <div v-if="activeSubtab === 'mysql'">
      <DatabaseManager engine="MySQL" />
    </div>

    <!-- SUBTAB 3: POSTGRESQL -->
    <div v-if="activeSubtab === 'postgres'">
      <DatabaseManager engine="PostgreSQL" />
    </div>

    <!-- SUBTAB 4: PROCESSES -->
    <div v-if="activeSubtab === 'processes'" class="bg-dark-900 border border-slate-800 rounded-2xl p-5 space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-bold text-white">System Processes</h3>
        <span class="text-xs font-mono text-slate-400">Total: {{ systemMetrics?.process_count || 0 }}</span>
      </div>
      <p class="text-xs text-slate-400">
        Active child processes and server workers are monitored safely. To view project-specific dev servers, open the project detail page.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { api } from '@/services/api';
import { useSystemStore } from '@/stores/system';
import DatabaseManager from '@/components/database/DatabaseManager.vue';
import { Layers, Database, HardDrive, Activity } from 'lucide-vue-next';

const route = useRoute();
const systemStore = useSystemStore();

const activeSubtab = ref('php');
const phpData = ref<any>(null);
const systemMetrics = ref<any>(null);

const subtabs = [
  { id: 'php', label: 'PHP Runtime', icon: Layers },
  { id: 'mysql', label: 'MySQL', icon: Database },
  { id: 'postgres', label: 'PostgreSQL', icon: HardDrive },
  { id: 'processes', label: 'Processes', icon: Activity },
];

onMounted(async () => {
  if (route.params.type) {
    activeSubtab.value = route.params.type as string;
  }
  loadPHP();
  systemMetrics.value = await systemStore.fetchMetrics();
});

async function loadPHP() {
  try {
    const res = await api.get('/php/diagnostics');
    phpData.value = res.data.data;
  } catch (e) {
    console.error(e);
  }
}
</script>
