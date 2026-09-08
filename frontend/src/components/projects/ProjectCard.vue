<template>
  <div class="glass-card rounded-xl p-5 flex flex-col justify-between group relative overflow-hidden border border-slate-800 hover:border-brand-500/50 transition duration-300">
    <!-- Top Meta & Status -->
    <div>
      <div class="flex items-start justify-between gap-3 mb-3">
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-slate-100 text-base group-hover:text-brand-400 transition tracking-tight">
              {{ project.name }}
            </h3>
            <span 
              :class="[
                'text-[10px] font-medium px-2 py-0.5 rounded-full border',
                project.status === 'Running' 
                  ? 'bg-emerald-950/60 text-emerald-400 border-emerald-800/60' 
                  : project.status === 'Development' 
                  ? 'bg-blue-950/60 text-blue-400 border-blue-800/60'
                  : 'bg-slate-800 text-slate-400 border-slate-700'
              ]"
            >
              ● {{ project.status }}
            </span>
          </div>
          <p class="text-xs text-slate-400 line-clamp-2 mt-1">{{ project.description || 'No description provided.' }}</p>
        </div>

        <button 
          @click="$emit('open', project.id)" 
          class="p-1.5 text-slate-400 hover:text-white rounded-lg hover:bg-slate-800 transition"
          title="Open Project"
        >
          <ExternalLink class="w-4 h-4" />
        </button>
      </div>

      <!-- Screenshot / Preview Frame -->
      <div 
        class="w-full h-32 rounded-lg bg-dark-950 border border-slate-800/80 mb-4 overflow-hidden relative cursor-pointer group/img"
        @click="$emit('preview', project)"
      >
        <img 
          v-if="project.screenshot_url" 
          :src="project.screenshot_url" 
          alt="Preview" 
          class="w-full h-full object-cover group-hover/img:scale-105 transition duration-300" 
        />
        <div v-else class="w-full h-full flex flex-col items-center justify-center text-slate-600">
          <Globe class="w-6 h-6 mb-1" />
          <span class="text-[11px] font-mono">Port :{{ project.preview_port }}</span>
        </div>
        <div class="absolute inset-0 bg-brand-600/10 opacity-0 group-hover/img:opacity-100 transition flex items-center justify-center">
          <span class="bg-dark-900/90 text-white text-[11px] px-2.5 py-1 rounded-md shadow border border-slate-700">Open Preview</span>
        </div>
      </div>

      <!-- Tech Stack Badges -->
      <div class="flex flex-wrap gap-1.5 mb-4">
        <span class="text-[11px] font-mono px-2 py-0.5 rounded bg-dark-850 text-slate-300 border border-slate-800">
          {{ project.framework }}
        </span>
        <span class="text-[11px] font-mono px-2 py-0.5 rounded bg-dark-850 text-slate-300 border border-slate-800">
          {{ project.frontend_tech }}
        </span>
        <span class="text-[11px] font-mono px-2 py-0.5 rounded bg-dark-850 text-slate-300 border border-slate-800">
          {{ project.database_type }}
        </span>
        <span class="text-[11px] font-mono px-2 py-0.5 rounded bg-dark-850 text-slate-400 border border-slate-800 flex items-center gap-1">
          <GitBranch class="w-3 h-3 text-slate-400" />
          {{ project.git_branch || 'main' }}
        </span>
      </div>

      <!-- TODO Progress -->
      <div class="mb-4">
        <div class="flex justify-between text-xs text-slate-400 mb-1 font-mono">
          <span>TODO Progress</span>
          <span class="text-slate-300 font-semibold">{{ todoPercent }}% ({{ completedCount }}/{{ totalCount }})</span>
        </div>
        <div class="w-full h-1.5 bg-dark-850 rounded-full overflow-hidden border border-slate-850">
          <div 
            class="h-full bg-gradient-to-r from-brand-600 to-indigo-500 rounded-full transition-all duration-500" 
            :style="{ width: `${todoPercent}%` }"
          ></div>
        </div>
      </div>
    </div>

    <!-- Bottom Actions -->
    <div class="flex items-center justify-between pt-3 border-t border-slate-800/80 gap-2">
      <div class="flex items-center gap-1.5">
        <router-link 
          :to="`/projects/${project.id}`" 
          class="px-2.5 py-1 rounded-md text-xs font-medium bg-brand-600/20 text-brand-400 hover:bg-brand-600 hover:text-white transition"
        >
          Open
        </router-link>
        <button 
          @click="$emit('terminal', project)"
          class="px-2.5 py-1 rounded-md text-xs font-medium bg-dark-800 text-slate-300 hover:bg-slate-700 hover:text-white transition flex items-center gap-1"
        >
          <Terminal class="w-3 h-3" />
          Terminal
        </button>
      </div>

      <button 
        @click="$emit('preview', project)"
        class="px-2.5 py-1 rounded-md text-xs font-medium bg-dark-800 text-slate-300 hover:bg-slate-700 hover:text-white transition flex items-center gap-1"
      >
        <Eye class="w-3 h-3" />
        Preview
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Project } from '@/types';
import { ExternalLink, Globe, GitBranch, Terminal, Eye } from 'lucide-vue-next';

const props = defineProps<{
  project: Project;
}>();

defineEmits(['open', 'terminal', 'preview']);

const totalCount = computed(() => props.project.tasks?.length || 0);
const completedCount = computed(() => {
  if (!props.project.tasks) return 0;
  return props.project.tasks.filter(t => t.status === 'done').length;
});

const todoPercent = computed(() => {
  if (totalCount.value === 0) return 0;
  return Math.round((completedCount.value / totalCount.value) * 100);
});
</script>
