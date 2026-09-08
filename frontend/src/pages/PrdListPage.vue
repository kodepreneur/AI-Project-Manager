<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-white tracking-tight">Product Requirements & Specs</h2>
        <p class="text-xs text-slate-400">Manage 17-section structured PRDs, review architecture, and approve milestones.</p>
      </div>

      <div class="flex items-center gap-2">
        <label class="text-xs font-semibold text-slate-400">Select Project:</label>
        <select 
          v-model="selectedProjectId"
          class="px-3 py-1.5 bg-dark-900 border border-slate-800 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
        >
          <option v-for="p in projectStore.projects" :key="p.id" :value="p.id">
            {{ p.name }}
          </option>
        </select>
      </div>
    </div>

    <!-- Render PRD Editor for selected project -->
    <div v-if="selectedProjectId">
      <PrdEditor :project-id="selectedProjectId" @approved="onPrdApproved" />
    </div>

    <div v-else class="p-12 text-center bg-dark-900 border border-slate-800 rounded-2xl">
      <FileText class="w-10 h-10 text-slate-600 mx-auto mb-2" />
      <h3 class="text-sm font-bold text-white">No project selected</h3>
      <p class="text-xs text-slate-400 mt-1">Select or create a project to generate its PRD.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useProjectStore } from '@/stores/projects';
import PrdEditor from '@/components/prd/PrdEditor.vue';
import { FileText } from 'lucide-vue-next';

const projectStore = useProjectStore();
const selectedProjectId = ref<number | null>(null);

onMounted(async () => {
  await projectStore.fetchProjects();
  if (projectStore.projects.length > 0) {
    selectedProjectId.value = projectStore.projects[0].id;
  }
});

function onPrdApproved() {
  projectStore.fetchProjects();
}
</script>
