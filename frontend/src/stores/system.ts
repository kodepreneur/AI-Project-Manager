import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '@/services/api';
import type { SystemMetrics } from '@/types';

export const useSystemStore = defineStore('system', () => {
  const metrics = ref<SystemMetrics | null>(null);
  const overview = ref<any>(null);
  const loading = ref<boolean>(false);

  async function fetchMetrics() {
    try {
      const res = await api.get('/system/metrics');
      metrics.value = res.data.data;
      return metrics.value;
    } catch {
      return null;
    }
  }

  async function fetchOverview() {
    try {
      const res = await api.get('/system/overview');
      overview.value = res.data.data;
      return overview.value;
    } catch {
      return null;
    }
  }

  return {
    metrics,
    overview,
    loading,
    fetchMetrics,
    fetchOverview,
  };
});
