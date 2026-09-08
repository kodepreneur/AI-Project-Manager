<template>
  <div class="min-h-screen bg-dark-950 flex flex-col items-center justify-center p-4 relative overflow-hidden">
    <!-- Ambient glow behind card -->
    <div class="absolute w-96 h-96 bg-brand-600/10 rounded-full blur-3xl pointer-events-none -top-20 -left-20"></div>
    <div class="absolute w-96 h-96 bg-indigo-600/10 rounded-full blur-3xl pointer-events-none -bottom-20 -right-20"></div>

    <div class="w-full max-w-md relative z-10">
      <router-view />
    </div>
  </div>
</template>
