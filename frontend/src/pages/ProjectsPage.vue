<template>
  <div class="space-y-6">
    <!-- Header with Search and Create -->
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-white tracking-tight">Software Projects</h2>
        <p class="text-xs text-slate-400">Manage codebases, runtimes, PRDs, and AI development lifecycles.</p>
      </div>

      <div class="flex items-center gap-3">
        <!-- Search -->
        <div class="relative">
          <Search class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input 
            v-model="search"
            type="text" 
            placeholder="Search projects..."
            class="pl-9 pr-4 py-1.5 bg-dark-900 border border-slate-800 rounded-lg text-xs text-white placeholder-slate-500 focus:outline-none focus:border-brand-500 w-48 sm:w-64"
            @input="filterProjects"
          />
        </div>

        <button 
          @click="showWizard = true"
          class="px-3.5 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow-md shadow-brand-600/20 transition flex items-center gap-1.5"
        >
          <Plus class="w-4 h-4" />
          <span>New Project</span>
        </button>
      </div>
    </div>

    <!-- Status Tabs -->
    <div class="flex items-center gap-2 border-b border-slate-800 pb-3 overflow-x-auto text-xs">
      <button 
        v-for="st in statusTabs" 
        :key="st.value"
        :class="[
          'px-3 py-1 rounded-md font-medium whitespace-nowrap transition',
          currentStatus === st.value ? 'bg-brand-600/20 text-brand-400 border border-brand-500/30' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'
        ]"
        @click="selectStatus(st.value)"
      >
        {{ st.label }}
      </button>
    </div>

    <!-- Project Cards Grid -->
    <div v-if="projectStore.projects.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <ProjectCard 
        v-for="project in projectStore.projects" 
        :key="project.id" 
        :project="project"
        @open="openProject"
        @terminal="openTerminal"
        @preview="openPreview"
      />
    </div>

    <!-- Empty State -->
    <div v-else class="p-12 text-center bg-dark-900 border border-slate-800 rounded-2xl">
      <FolderKanban class="w-12 h-12 text-slate-600 mx-auto mb-3" />
      <h3 class="text-base font-bold text-white">No projects found</h3>
      <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
        No projects match your current filter. Create a new software project using the wizard.
      </p>
      <button 
        @click="showWizard = true"
        class="mt-4 px-4 py-2 rounded-xl text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white transition"
      >
        Create Project
      </button>
    </div>

    <!-- Wizard Modal -->
    <ProjectWizardModal :is-open="showWizard" @close="showWizard = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProjectStore } from '@/stores/projects';
import ProjectCard from '@/components/projects/ProjectCard.vue';
import ProjectWizardModal from '@/components/projects/ProjectWizardModal.vue';
import { Search, Plus, FolderKanban } from 'lucide-vue-next';

const route = useRoute();
const router = useRouter();
const projectStore = useProjectStore();

const showWizard = ref(false);
const search = ref('');
const currentStatus = ref('');

const statusTabs = [
  { label: 'All Projects', value: '' },
  { label: 'Running', value: 'Running' },
  { label: 'Development', value: 'Development' },
  { label: 'Archived', value: 'Archived' },
];

onMounted(() => {
  if (route.query.status) {
    currentStatus.value = route.query.status as string;
  }
  if (route.query.create === '1') {
    showWizard.value = true;
  }
  loadProjects();
});

watch(() => route.query.status, (newStatus) => {
  currentStatus.value = (newStatus as string) || '';
  loadProjects();
});

function selectStatus(val: string) {
  currentStatus.value = val;
  router.push({ query: val ? { status: val } : {} });
}

function filterProjects() {
  projectStore.searchQuery = search.value;
  projectStore.fetchProjects();
}

function loadProjects() {
  projectStore.filterStatus = currentStatus.value;
  projectStore.fetchProjects();
}

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
