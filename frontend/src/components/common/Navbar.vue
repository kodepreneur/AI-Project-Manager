<template>
  <header class="h-14 border-b border-slate-800/80 bg-dark-900/90 backdrop-blur-md px-4 flex items-center justify-between sticky top-0 z-30">
    <div class="flex items-center gap-3">
      <button 
        @click="$emit('toggle-sidebar')" 
        class="p-1.5 text-slate-400 hover:text-white rounded-lg hover:bg-slate-800 lg:hidden"
      >
        <Menu class="w-5 h-5" />
      </button>

      <div class="flex items-center gap-2 text-sm">
        <span class="text-slate-400 font-medium">Control Panel</span>
        <span class="text-slate-600">/</span>
        <span class="text-slate-200 font-semibold">{{ currentTitle }}</span>
      </div>
    </div>

    <!-- Center/Right Search & Actions -->
    <div class="flex items-center gap-3">
      <!-- Command Palette Trigger -->
      <button 
        @click="$emit('open-palette')"
        class="hidden sm:flex items-center gap-2 bg-dark-850 hover:bg-dark-800 border border-slate-800 hover:border-slate-700 px-3 py-1.5 rounded-lg text-xs text-slate-400 hover:text-slate-200 transition"
      >
        <Search class="w-3.5 h-3.5" />
        <span>Quick search...</span>
        <kbd class="bg-dark-750 border border-slate-700 text-slate-400 px-1.5 py-0.5 rounded text-[10px] font-mono">⌘K</kbd>
      </button>

      <!-- Antigravity Global Pulse -->
      <div class="flex items-center gap-2 px-2.5 py-1 rounded-full bg-blue-950/40 border border-blue-800/40 text-blue-400 text-xs">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
        </span>
        <span class="font-medium hidden md:inline">Antigravity Engine Active</span>
      </div>

      <!-- User Profile Menu -->
      <div class="flex items-center gap-2 border-l border-slate-800 pl-3">
        <div class="w-7 h-7 rounded-full bg-brand-600 text-white flex items-center justify-center text-xs font-bold shadow">
          {{ userInitials }}
        </div>
        <span class="text-xs font-medium text-slate-300 hidden sm:inline">{{ authStore.user?.username || 'Admin' }}</span>
        <button 
          @click="handleLogout"
          class="p-1.5 text-slate-400 hover:text-rose-400 rounded-lg hover:bg-slate-800 transition"
          title="Logout"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { Menu, Search, LogOut } from 'lucide-vue-next';

defineEmits(['toggle-sidebar', 'open-palette']);

const route = useRoute();
const authStore = useAuthStore();

const currentTitle = computed(() => {
  return (route.meta?.title as string) || 'Dashboard';
});

const userInitials = computed(() => {
  const name = authStore.user?.username || 'Admin';
  return name.slice(0, 2).toUpperCase();
});

function handleLogout() {
  if (confirm('Are you sure you want to log out?')) {
    authStore.logout();
  }
}
</script>
