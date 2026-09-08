<template>
  <div 
    v-if="isOpen" 
    class="fixed inset-0 z-50 flex items-start justify-center pt-20 px-4 bg-black/70 backdrop-blur-sm animate-fade-in"
    @click.self="close"
  >
    <div class="w-full max-w-xl bg-dark-900 border border-slate-700/80 rounded-xl shadow-2xl overflow-hidden flex flex-col">
      <!-- Search Input Header -->
      <div class="flex items-center px-4 py-3 border-b border-slate-800 gap-3">
        <Search class="w-4 h-4 text-slate-400" />
        <input 
          ref="searchInput"
          v-model="query"
          type="text" 
          placeholder="Search commands, projects, actions..."
          class="w-full bg-transparent text-sm text-white placeholder-slate-500 focus:outline-none"
          @keydown.esc="close"
          @keydown.down.prevent="navigate(1)"
          @keydown.up.prevent="navigate(-1)"
          @keydown.enter.prevent="executeSelected"
        />
        <kbd class="bg-dark-800 border border-slate-700 text-slate-400 px-1.5 py-0.5 rounded text-[10px] font-mono">ESC</kbd>
      </div>

      <!-- Action / Result Items List -->
      <div class="max-h-80 overflow-y-auto p-2 space-y-1">
        <div 
          v-for="(item, idx) in filteredItems" 
          :key="idx"
          :class="[
            'px-3 py-2 rounded-lg text-xs flex items-center justify-between cursor-pointer transition',
            selectedIndex === idx ? 'bg-brand-600 text-white' : 'text-slate-300 hover:bg-dark-800 hover:text-white'
          ]"
          @click="item.action"
          @mouseenter="selectedIndex = idx"
        >
          <div class="flex items-center gap-2.5">
            <component :is="item.icon" class="w-4 h-4 opacity-80" />
            <span class="font-medium">{{ item.title }}</span>
          </div>
          <span 
            :class="[
              'text-[10px] uppercase font-mono px-1.5 py-0.5 rounded',
              selectedIndex === idx ? 'bg-brand-700 text-brand-100' : 'bg-dark-800 text-slate-400'
            ]"
          >
            {{ item.category }}
          </span>
        </div>

        <div v-if="filteredItems.length === 0" class="text-center py-6 text-xs text-slate-500">
          No matching commands or projects found.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useProjectStore } from '@/stores/projects';
import { 
  Search, Plus, LayoutDashboard, FolderKanban, Columns3, 
  FileText, Bot, Layers, Database, HardDrive, ScrollText, Settings 
} from 'lucide-vue-next';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits(['close']);

const router = useRouter();
const projectStore = useProjectStore();

const query = ref('');
const selectedIndex = ref(0);
const searchInput = ref<HTMLInputElement | null>(null);

watch(() => props.isOpen, (open) => {
  if (open) {
    query.value = '';
    selectedIndex.value = 0;
    nextTick(() => {
      searchInput.value?.focus();
    });
  }
});

const defaultCommands = [
  { title: 'Dashboard Overview', category: 'Navigation', icon: LayoutDashboard, action: () => navigateTo('/dashboard') },
  { title: 'Create New Project', category: 'Action', icon: Plus, action: () => navigateTo('/projects?create=1') },
  { title: 'All Projects', category: 'Projects', icon: FolderKanban, action: () => navigateTo('/projects') },
  { title: 'TODO Board / Kanban', category: 'Development', icon: Columns3, action: () => navigateTo('/kanban') },
  { title: 'PRDs & Specifications', category: 'Development', icon: FileText, action: () => navigateTo('/prds') },
  { title: 'Antigravity AI Sessions', category: 'AI', icon: Bot, action: () => navigateTo('/antigravity') },
  { title: 'PHP Runtime Diagnostics', category: 'Infrastructure', icon: Layers, action: () => navigateTo('/infrastructure/php') },
  { title: 'MySQL Management', category: 'Infrastructure', icon: Database, action: () => navigateTo('/infrastructure/mysql') },
  { title: 'PostgreSQL Management', category: 'Infrastructure', icon: HardDrive, action: () => navigateTo('/infrastructure/postgres') },
  { title: 'Centralized Logs', category: 'System', icon: ScrollText, action: () => navigateTo('/logs') },
  { title: 'System Settings', category: 'System', icon: Settings, action: () => navigateTo('/settings') },
];

const filteredItems = computed(() => {
  const q = query.value.toLowerCase().trim();
  const projectItems = projectStore.projects.map(p => ({
    title: `Project: ${p.name}`,
    category: 'Project',
    icon: FolderKanban,
    action: () => navigateTo(`/projects/${p.id}`),
  }));

  const all = [...defaultCommands, ...projectItems];
  if (!q) return all;

  return all.filter(item => 
    item.title.toLowerCase().includes(q) || 
    item.category.toLowerCase().includes(q)
  );
});

function navigate(dir: number) {
  if (filteredItems.value.length === 0) return;
  selectedIndex.value = (selectedIndex.value + dir + filteredItems.value.length) % filteredItems.value.length;
}

function executeSelected() {
  const item = filteredItems.value[selectedIndex.value];
  if (item) {
    item.action();
  }
}

function navigateTo(path: string) {
  router.push(path);
  close();
}

function close() {
  emit('close');
}
</script>
