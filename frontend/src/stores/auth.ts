import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '@/services/api';
import type { User } from '@/types';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('token'));
  const isSetup = ref<boolean>(true);
  const loading = ref<boolean>(false);

  async function checkSetupStatus() {
    try {
      const res = await api.get('/auth/setup-status');
      isSetup.value = res.data.data.is_setup;
      return isSetup.value;
    } catch {
      isSetup.value = true;
      return true;
    }
  }

  async function setupAdmin(payload: { username: string; email: string; password: string }) {
    loading.value = true;
    try {
      const res = await api.post('/auth/setup', payload);
      token.value = res.data.data.token;
      user.value = res.data.data.user;
      localStorage.setItem('token', token.value!);
      isSetup.value = true;
      return true;
    } finally {
      loading.value = false;
    }
  }

  async function login(usernameOrEmail: string, password: string) {
    loading.value = true;
    try {
      const res = await api.post('/auth/login', {
        username_or_email: usernameOrEmail,
        password,
      });
      token.value = res.data.data.token;
      user.value = res.data.data.user;
      localStorage.setItem('token', token.value!);
      return true;
    } finally {
      loading.value = false;
    }
  }

  async function fetchProfile() {
    if (!token.value) return null;
    try {
      const res = await api.get('/auth/me');
      user.value = res.data.data;
      return user.value;
    } catch {
      logout();
      return null;
    }
  }

  function logout() {
    token.value = null;
    user.value = null;
    localStorage.removeItem('token');
    window.location.href = '/login';
  }

  return {
    user,
    token,
    isSetup,
    loading,
    checkSetupStatus,
    setupAdmin,
    login,
    fetchProfile,
    logout,
  };
});
