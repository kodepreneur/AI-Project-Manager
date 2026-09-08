<template>
  <div class="space-y-4">
    <!-- Top Preview Control Bar -->
    <div class="flex flex-wrap items-center justify-between gap-3 p-4 bg-dark-900 border border-slate-800 rounded-xl">
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
          <Globe class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-sm text-white">Application Preview</h3>
            <span 
              :class="[
                'text-[10px] font-mono px-2 py-0.5 rounded-full border',
                previewInfo?.is_responding
                  ? 'bg-emerald-950/70 text-emerald-400 border-emerald-800'
                  : 'bg-amber-950/70 text-amber-400 border-amber-800'
              ]"
            >
              {{ previewInfo?.is_responding ? '● Live Responding' : '○ Standby / Starting' }}
            </span>
          </div>
          <p class="text-xs text-slate-400 font-mono">
            {{ previewInfo?.preview_url || `http://localhost:${project.preview_port}` }}
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button 
          @click="captureScreenshot"
          :disabled="capturing"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-dark-800 hover:bg-dark-750 text-slate-300 border border-slate-700 transition flex items-center gap-1.5"
        >
          <Camera class="w-3.5 h-3.5" />
          <span>{{ capturing ? 'Capturing...' : 'Capture Screenshot' }}</span>
        </button>

        <a 
          :href="previewInfo?.preview_url || `http://localhost:${project.preview_port}`" 
          target="_blank"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Open in New Tab</span>
        </a>
      </div>
    </div>

    <!-- Live Iframe or Screenshot Preview -->
    <div class="h-[550px] bg-dark-950 border border-slate-800 rounded-xl overflow-hidden shadow-2xl relative flex flex-col">
      <!-- Browser Top Address bar simulation -->
      <div class="h-9 px-4 bg-dark-900 border-b border-slate-800 flex items-center gap-3">
        <div class="flex items-center gap-1.5">
          <span class="w-2.5 h-2.5 rounded-full bg-slate-700"></span>
          <span class="w-2.5 h-2.5 rounded-full bg-slate-700"></span>
          <span class="w-2.5 h-2.5 rounded-full bg-slate-700"></span>
        </div>
        <div class="flex-1 max-w-md mx-auto bg-dark-950 border border-slate-800 rounded px-3 py-1 text-[11px] font-mono text-slate-400 text-center truncate">
          {{ previewInfo?.preview_url || `http://localhost:${project.preview_port}` }}
        </div>
        <button @click="fetchPreviewInfo" class="text-slate-400 hover:text-white" title="Refresh">
          <RotateCcw class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Iframe -->
      <iframe 
        v-if="previewInfo?.is_responding"
        :src="previewInfo.preview_url" 
        class="w-full flex-1 border-0 bg-white"
      ></iframe>

      <!-- Standby fallback view with Screenshot -->
      <div v-else class="flex-1 flex flex-col items-center justify-center p-8 text-center bg-dark-950 relative overflow-hidden">
        <img 
          v-if="previewInfo?.screenshot_url" 
          :src="previewInfo.screenshot_url" 
          alt="Screenshot Preview" 
          class="max-w-xl max-h-72 rounded-lg border border-slate-800 shadow-xl mb-4 object-cover" 
        />
        <div v-else class="w-16 h-16 rounded-2xl bg-dark-900 border border-slate-800 flex items-center justify-center text-slate-500 mb-3">
          <Globe class="w-8 h-8" />
        </div>
        <h4 class="text-sm font-bold text-white mb-1">Development Server Not Responding</h4>
        <p class="text-xs text-slate-400 max-w-md mb-4">
          Start the project runtime process (e.g. <code class="text-brand-400 font-mono">{{ project.preview_command }}</code>) in the Processes or Terminal tab to view live preview.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { api } from '@/services/api';
import type { Project } from '@/types';
import { Globe, Camera, ExternalLink, RotateCcw } from 'lucide-vue-next';

const props = defineProps<{
  project: Project;
}>();

const previewInfo = ref<any>(null);
const capturing = ref(false);

onMounted(() => {
  fetchPreviewInfo();
});

async function fetchPreviewInfo() {
  try {
    const res = await api.get(`/projects/${props.project.id}/preview`);
    previewInfo.value = res.data.data;
  } catch {
    previewInfo.value = null;
  }
}

async function captureScreenshot() {
  capturing.value = true;
  try {
    const res = await api.post(`/projects/${props.project.id}/preview/screenshot`);
    if (previewInfo.value) {
      previewInfo.value.screenshot_url = res.data.data.screenshot_url;
    }
  } finally {
    capturing.value = false;
  }
}
</script>
