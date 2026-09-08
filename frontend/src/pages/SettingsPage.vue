<template>
  <div class="space-y-6 max-w-3xl">
    <div>
      <h2 class="text-xl font-bold text-white tracking-tight">System & AI Settings</h2>
      <p class="text-xs text-slate-400">Configure global directories, Antigravity CLI binary paths, and daemon defaults.</p>
    </div>

    <div class="bg-dark-900 border border-slate-800 rounded-2xl p-6 space-y-5">
      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1.5">Application Name</label>
        <input 
          v-model="settings.app_name"
          type="text" 
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
        />
      </div>

      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1.5">Antigravity CLI Binary Path</label>
        <div class="flex gap-2">
          <input 
            v-model="settings.antigravity_binary"
            type="text" 
            placeholder="antigravity"
            class="flex-1 px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white font-mono focus:outline-none focus:border-brand-500"
          />
          <button 
            type="button"
            @click="testCliConnection"
            :disabled="testingCli"
            class="px-3.5 py-2 rounded-lg text-xs font-semibold bg-purple-950/80 hover:bg-purple-900 border border-purple-700/60 text-purple-300 transition flex items-center gap-1.5"
          >
            <Bot class="w-3.5 h-3.5" />
            <span>{{ testingCli ? 'Testing...' : 'Test CLI' }}</span>
          </button>
        </div>
        <p class="text-[11px] text-slate-500 mt-1">Command or binary path (e.g. 'agy', 'antigravity-ide', or full path).</p>

        <!-- Live CLI Test Result -->
        <div v-if="cliStatus" class="mt-3 p-3 rounded-lg border text-xs" :class="cliStatus.success ? 'bg-purple-950/30 border-purple-700/50 text-purple-200' : 'bg-rose-950/40 border-rose-800 text-rose-300'">
          <div class="flex items-center justify-between font-semibold">
            <span>{{ cliStatus.success ? '✔ ' + cliStatus.data.status_summary : '✖ Connection Failed' }}</span>
            <span v-if="cliStatus.success" class="font-mono text-[10px] text-purple-400">{{ cliStatus.data.latency_ms }}ms</span>
          </div>
          <div v-if="cliStatus.success" class="mt-1 text-[11px] font-mono text-slate-400">
            Binary: {{ cliStatus.data.binary_path }}
          </div>
          <div v-else class="mt-1 text-[11px]">
            {{ cliStatus.error?.message }} ({{ cliStatus.error?.hint }})
          </div>
        </div>
      </div>

      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1.5">Default Project Storage Root</label>
        <input 
          v-model="settings.project_root"
          type="text" 
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white font-mono focus:outline-none focus:border-brand-500"
        />
      </div>

      <div class="space-y-3 pt-2">
        <label class="flex items-center gap-3 cursor-pointer">
          <input 
            type="checkbox" 
            v-model="settings.terminal_enabled"
            class="rounded bg-dark-850 border-slate-700 text-brand-600 focus:ring-0" 
          />
          <div>
            <div class="text-xs font-semibold text-white">Enable Web PTY Terminal</div>
            <p class="text-[11px] text-slate-400">Allows interactive pseudo-terminal sessions directly from the browser.</p>
          </div>
        </label>

        <label class="flex items-center gap-3 cursor-pointer">
          <input 
            type="checkbox" 
            v-model="settings.screenshot_enabled"
            class="rounded bg-dark-850 border-slate-700 text-brand-600 focus:ring-0" 
          />
          <div>
            <div class="text-xs font-semibold text-white">Enable Automated Preview Screenshots</div>
            <p class="text-[11px] text-slate-400">Captures visual snapshots of running project web ports using Chromium or fallback cards.</p>
          </div>
        </label>
      </div>

      <div class="pt-4 border-t border-slate-800 flex justify-end">
        <button 
          @click="saveSettings"
          :disabled="saving"
          class="px-5 py-2 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Save class="w-3.5 h-3.5" />
          <span>{{ saving ? 'Saving...' : 'Save Configuration' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { api } from '@/services/api';
import { Save, Bot } from 'lucide-vue-next';

const settings = reactive({
  app_name: 'AI Project Manager',
  antigravity_binary: 'antigravity',
  project_root: '',
  terminal_enabled: true,
  screenshot_enabled: true,
});

const saving = ref(false);
const testingCli = ref(false);
const cliStatus = ref<any>(null);

onMounted(async () => {
  try {
    const res = await api.get('/settings');
    Object.assign(settings, res.data.data);
    testCliConnection();
  } catch (e) {
    console.error(e);
  }
});

async function testCliConnection() {
  testingCli.value = true;
  cliStatus.value = null;
  try {
    const res = await api.post('/antigravity/test', { binary: settings.antigravity_binary });
    cliStatus.value = res.data;
  } catch (err: any) {
    cliStatus.value = {
      success: false,
      error: err.response?.data?.error || { message: 'Failed to communicate with Antigravity CLI' },
    };
  } finally {
    testingCli.value = false;
  }
}

async function saveSettings() {
  saving.value = true;
  try {
    await api.put('/settings', settings);
    alert('Settings updated successfully.');
  } finally {
    saving.value = false;
  }
}
</script>
