<template>
  <div class="glass-panel p-8 rounded-2xl shadow-2xl border border-slate-800">
    <div class="text-center mb-6">
      <div class="w-12 h-12 rounded-xl bg-gradient-to-tr from-brand-600 to-indigo-500 flex items-center justify-center mx-auto mb-3 shadow-lg shadow-brand-500/20">
        <Sparkles class="w-6 h-6 text-white" />
      </div>
      <h1 class="text-xl font-bold text-white tracking-tight">AI Project Manager</h1>
      <p class="text-xs text-slate-400 mt-1">Sign in to your development control panel</p>
    </div>

    <form @submit.prevent="handleLogin" class="space-y-4">
      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1.5">Username or Email</label>
        <input 
          v-model="usernameOrEmail"
          type="text" 
          required
          placeholder="admin"
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
        />
      </div>

      <div>
        <label class="block text-xs font-semibold text-slate-300 mb-1.5">Password</label>
        <input 
          v-model="password"
          type="password" 
          required
          placeholder="••••••••••••"
          class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500 font-mono"
        />
      </div>

      <div v-if="error" class="p-3 rounded-lg bg-rose-950/50 border border-rose-800 text-rose-300 text-xs">
        {{ error }}
      </div>

      <button 
        type="submit"
        :disabled="loading"
        class="w-full py-2.5 rounded-lg text-sm font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow-lg shadow-brand-600/20 transition flex items-center justify-center gap-2"
      >
        <span v-if="loading">Signing in...</span>
        <span v-else>Access Control Panel</span>
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { Sparkles } from 'lucide-vue-next';

const router = useRouter();
const authStore = useAuthStore();

const usernameOrEmail = ref('');
const password = ref('');
const error = ref('');
const loading = ref(false);

async function handleLogin() {
  error.value = '';
  loading.value = true;
  try {
    await authStore.login(usernameOrEmail.value, password.value);
    router.push('/dashboard');
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Invalid username or password.';
  } finally {
    loading.value = false;
  }
}
</script>
