import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from '@/services/api';
import type { Task } from '@/types';

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([]);
  const loading = ref<boolean>(false);

  async function fetchTasks(projectId: number | string) {
    loading.value = true;
    try {
      const res = await api.get(`/projects/${projectId}/tasks`);
      tasks.value = res.data.data;
      return tasks.value;
    } finally {
      loading.value = false;
    }
  }

  async function createTask(projectId: number | string, payload: Partial<Task>) {
    const res = await api.post(`/projects/${projectId}/tasks`, payload);
    tasks.value.push(res.data.data);
    return res.data.data;
  }

  async function updateTask(taskId: number, payload: Partial<Task>) {
    const res = await api.put(`/tasks/${taskId}`, payload);
    const index = tasks.value.findIndex(t => t.id === taskId);
    if (index !== -1) {
      tasks.value[index] = res.data.data;
    }
    return res.data.data;
  }

  async function updateTaskStatus(taskId: number, status: Task['status']) {
    // Optimistic update
    const index = tasks.value.findIndex(t => t.id === taskId);
    if (index !== -1) {
      tasks.value[index].status = status;
    }
    try {
      await api.put(`/tasks/${taskId}/status`, { status });
    } catch {
      // Revert if error
    }
  }

  async function deleteTask(taskId: number) {
    await api.delete(`/tasks/${taskId}`);
    tasks.value = tasks.value.filter(t => t.id !== taskId);
  }

  async function generateTasksFromPRD(projectId: number | string) {
    loading.value = true;
    try {
      const res = await api.post(`/projects/${projectId}/tasks/generate-from-prd`);
      await fetchTasks(projectId);
      return res.data;
    } finally {
      loading.value = false;
    }
  }

  return {
    tasks,
    loading,
    fetchTasks,
    createTask,
    updateTask,
    updateTaskStatus,
    deleteTask,
    generateTasksFromPRD,
  };
});
