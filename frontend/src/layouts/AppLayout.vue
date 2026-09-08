<template>
  <div class="min-h-screen bg-dark-950 flex">
    <!-- Sidebar -->
    <Sidebar :is-open="sidebarOpen" />

    <!-- Backdrop for mobile drawer -->
    <div 
      v-if="sidebarOpen" 
      class="fixed inset-0 bg-black/60 z-30 lg:hidden"
      @click="sidebarOpen = false"
    ></div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0">
      <Navbar 
        @toggle-sidebar="sidebarOpen = !sidebarOpen" 
        @open-palette="paletteOpen = true"
      />

      <main class="flex-1 p-4 lg:p-6 overflow-y-auto max-w-7xl w-full mx-auto">
        <router-view />
      </main>
    </div>

    <!-- Global Command Palette -->
    <CommandPalette :is-open="paletteOpen" @close="paletteOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import Navbar from '@/components/common/Navbar.vue';
import Sidebar from '@/components/common/Sidebar.vue';
import CommandPalette from '@/components/common/CommandPalette.vue';

const sidebarOpen = ref(false);
const paletteOpen = ref(false);

function handleKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault();
    paletteOpen.value = !paletteOpen.value;
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown);
});
</script>
