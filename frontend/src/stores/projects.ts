import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '@/services/api';
import type { Project } from '@/types';

export const useProjectStore = defineStore('projects', () => {
  const projects = ref<Project[]>([]);
  const currentProject = ref<Project | null>(null);
  const loading = ref<boolean>(false);
  const filterStatus = ref<string>('');
  const searchQuery = ref<string>('');

  async function fetchProjects() {
    loading.value = true;
    try {
      const params: Record<string, string> = {};
      if (filterStatus.value) params.status = filterStatus.value;
      if (searchQuery.value) params.search = searchQuery.value;

      const res = await api.get('/projects', { params });
      projects.value = res.data.data;
      return projects.value;
    } finally {
      loading.value = false;
    }
  }

  async function fetchProject(id: number | string) {
    loading.value = true;
    try {
      const res = await api.get(`/projects/${id}`);
      currentProject.value = res.data.data;
      return currentProject.value;
    } finally {
      loading.value = false;
    }
  }

  async function createProject(payload: Partial<Project> & { template?: string }) {
    loading.value = true;
    try {
      const res = await api.post('/projects', payload);
      await fetchProjects();
      return res.data.data;
    } finally {
      loading.value = false;
    }
  }

  async function updateProject(id: number, payload: Partial<Project>) {
    const res = await api.put(`/projects/${id}`, payload);
    if (currentProject.value?.id === id) {
      currentProject.value = { ...currentProject.value, ...res.data.data };
    }
    await fetchProjects();
    return res.data.data;
  }

  async function deleteProject(id: number) {
    await api.delete(`/projects/${id}`);
    await fetchProjects();
  }

  async function archiveProject(id: number) {
    const res = await api.post(`/projects/${id}/archive`);
    await fetchProjects();
    return res.data.data;
  }

  return {
    projects,
    currentProject,
    loading,
    filterStatus,
    searchQuery,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    archiveProject,
  };
});
