<template>
  <aside 
    :class="[
      'fixed inset-y-0 left-0 z-40 w-64 bg-dark-900 border-r border-slate-800/80 flex flex-col transition-transform duration-200 lg:static lg:translate-x-0',
      isOpen ? 'translate-x-0' : '-translate-x-full'
    ]"
  >
    <!-- Brand / Logo -->
    <div class="h-14 px-4 flex items-center gap-3 border-b border-slate-800/80 bg-dark-950/40">
      <div class="w-8 h-8 rounded-lg bg-gradient-to-tr from-brand-600 to-indigo-500 flex items-center justify-center shadow-lg shadow-brand-500/20">
        <Sparkles class="w-4 h-4 text-white" />
      </div>
      <div>
        <h1 class="text-sm font-bold text-white tracking-tight leading-none">AI Project Manager</h1>
        <span class="text-[10px] text-brand-400 font-mono">Antigravity Control</span>
      </div>
    </div>

    <!-- Navigation Links -->
    <div class="flex-1 overflow-y-auto px-3 py-4 space-y-6">
      <!-- Main Overview -->
      <div>
        <router-link 
          to="/dashboard" 
          class="flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition"
          :class="isActive('/dashboard') ? 'bg-brand-600 text-white shadow-md shadow-brand-600/20' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
        >
          <LayoutDashboard class="w-4 h-4" />
          <span>Dashboard</span>
        </router-link>
      </div>

      <!-- Projects Section -->
      <div>
        <div class="px-3 text-[11px] font-semibold text-slate-400 uppercase tracking-wider mb-1.5 flex items-center justify-between">
          <span>Projects</span>
          <router-link to="/projects?create=1" class="hover:text-white" title="New Project">
            <Plus class="w-3.5 h-3.5" />
          </router-link>
        </div>
        <div class="space-y-0.5">
          <router-link 
            to="/projects" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="route.path === '/projects' && !route.query.status ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <FolderKanban class="w-4 h-4 text-slate-400" />
            <span>All Projects</span>
          </router-link>
          <router-link 
            to="/projects?status=Running" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="route.query.status === 'Running' ? 'bg-dark-800 text-emerald-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <PlayCircle class="w-4 h-4 text-emerald-400" />
            <span>Running</span>
          </router-link>
          <router-link 
            to="/projects?status=Development" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="route.query.status === 'Development' ? 'bg-dark-800 text-blue-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Code2 class="w-4 h-4 text-blue-400" />
            <span>Development</span>
          </router-link>
          <router-link 
            to="/projects?status=Archived" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="route.query.status === 'Archived' ? 'bg-dark-800 text-slate-300' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Archive class="w-4 h-4 text-slate-400" />
            <span>Archived</span>
          </router-link>
        </div>
      </div>

      <!-- Development Section -->
      <div>
        <div class="px-3 text-[11px] font-semibold text-slate-400 uppercase tracking-wider mb-1.5">
          Development
        </div>
        <div class="space-y-0.5">
          <router-link 
            to="/kanban" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/kanban') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Columns3 class="w-4 h-4 text-slate-400" />
            <span>TODO Board</span>
          </router-link>
          <router-link 
            to="/prds" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/prds') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <FileText class="w-4 h-4 text-slate-400" />
            <span>PRDs & Specs</span>
          </router-link>
          <router-link 
            to="/antigravity" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/antigravity') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Bot class="w-4 h-4 text-purple-400" />
            <span>Antigravity AI</span>
          </router-link>
        </div>
      </div>

      <!-- Infrastructure Section -->
      <div>
        <div class="px-3 text-[11px] font-semibold text-slate-400 uppercase tracking-wider mb-1.5">
          Infrastructure
        </div>
        <div class="space-y-0.5">
          <router-link 
            to="/infrastructure/php" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/infrastructure/php') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Layers class="w-4 h-4 text-indigo-400" />
            <span>PHP Runtime</span>
          </router-link>
          <router-link 
            to="/infrastructure/mysql" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/infrastructure/mysql') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Database class="w-4 h-4 text-amber-400" />
            <span>MySQL</span>
          </router-link>
          <router-link 
            to="/infrastructure/postgres" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/infrastructure/postgres') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <HardDrive class="w-4 h-4 text-sky-400" />
            <span>PostgreSQL</span>
          </router-link>
          <router-link 
            to="/infrastructure/processes" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/infrastructure/processes') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Activity class="w-4 h-4 text-rose-400" />
            <span>Processes</span>
          </router-link>
        </div>
      </div>

      <!-- System Section -->
      <div>
        <div class="px-3 text-[11px] font-semibold text-slate-400 uppercase tracking-wider mb-1.5">
          System
        </div>
        <div class="space-y-0.5">
          <router-link 
            to="/logs" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/logs') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <ScrollText class="w-4 h-4 text-slate-400" />
            <span>Central Logs</span>
          </router-link>
          <router-link 
            to="/settings" 
            class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg text-xs font-medium transition"
            :class="isActive('/settings') ? 'bg-dark-800 text-brand-400' : 'text-slate-400 hover:text-slate-200 hover:bg-dark-850'"
          >
            <Settings class="w-4 h-4 text-slate-400" />
            <span>Settings</span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Footer System Status -->
    <div class="p-3 border-t border-slate-800/80 bg-dark-950/50">
      <div class="flex items-center justify-between text-[11px] text-slate-400">
        <span class="flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
          System Online
        </span>
        <span class="font-mono text-slate-400">v1.0.0</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { 
  Sparkles, LayoutDashboard, FolderKanban, PlayCircle, Code2, 
  Archive, Plus, Columns3, FileText, Bot, Layers, Database, 
  HardDrive, Activity, ScrollText, Settings 
} from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const route = useRoute();

function isActive(path: string) {
  return route.path === path || route.path.startsWith(path + '/');
}
</script>
