<template>
  <div class="space-y-4">
    <!-- Top Action Bar -->
    <div class="flex flex-wrap items-center justify-between gap-3 p-4 bg-dark-900 border border-slate-800 rounded-xl">
      <div>
        <h3 class="text-sm font-bold text-white flex items-center gap-2">
          <Columns3 class="w-4 h-4 text-brand-400" />
          <span>Development Tasks & Kanban Board</span>
        </h3>
        <p class="text-xs text-slate-400">Track milestones, assign to Antigravity AI, and validate implementations.</p>
      </div>

      <div class="flex items-center gap-2">
        <button 
          @click="generateTasks"
          :disabled="generating"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-gradient-to-r from-purple-600 to-brand-600 hover:from-purple-500 hover:to-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Sparkles class="w-3.5 h-3.5" />
          <span>{{ generating ? 'Analyzing PRD...' : 'Generate Plan from PRD' }}</span>
        </button>

        <button 
          @click="openCreateTaskModal"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>New Task</span>
        </button>
      </div>
    </div>

    <!-- Kanban Columns Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-3.5 items-start">
      <div 
        v-for="col in columns" 
        :key="col.status"
        class="bg-dark-900 border border-slate-850 rounded-xl p-3 flex flex-col min-h-[480px]"
      >
        <!-- Column Header -->
        <div class="flex items-center justify-between pb-2.5 mb-2.5 border-b border-slate-800">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full" :class="col.color"></span>
            <span class="text-xs font-bold uppercase tracking-wider text-slate-300">{{ col.title }}</span>
          </div>
          <span class="text-[11px] font-mono px-1.5 py-0.5 rounded bg-dark-800 text-slate-400">
            {{ getColumnTasks(col.status).length }}
          </span>
        </div>

        <!-- Task Cards Container -->
        <div class="space-y-2.5 flex-1">
          <div 
            v-for="task in getColumnTasks(col.status)" 
            :key="task.id"
            class="glass-card p-3 rounded-lg border border-slate-800 hover:border-brand-500/50 cursor-pointer space-y-2 text-xs"
            @click="selectTask(task)"
          >
            <div class="flex items-start justify-between gap-1">
              <span class="font-semibold text-slate-100 line-clamp-2">{{ task.title }}</span>
              <span 
                :class="[
                  'text-[9px] uppercase font-bold px-1.5 py-0.5 rounded',
                  task.priority === 'urgent' ? 'bg-rose-950 text-rose-400 border border-rose-800' :
                  task.priority === 'high' ? 'bg-amber-950 text-amber-400 border border-amber-800' :
                  'bg-dark-800 text-slate-400'
                ]"
              >
                {{ task.priority }}
              </span>
            </div>

            <p v-if="task.phase_name" class="text-[10px] text-slate-400 font-mono line-clamp-1">
              {{ task.phase_name }}
            </p>

            <!-- Bottom Row: AI Status & Quick Transition -->
            <div class="flex items-center justify-between pt-1 border-t border-slate-800/60" @click.stop>
              <span 
                v-if="task.ai_status === 'running'"
                class="text-[10px] text-purple-400 flex items-center gap-1"
              >
                <Bot class="w-3 h-3 animate-spin" />
                Antigravity running...
              </span>
              <span v-else class="text-[10px] text-slate-500">
                Complexity: {{ task.complexity }}
              </span>

              <!-- Send to AI Trigger -->
              <button 
                v-if="task.status !== 'done'"
                @click="openSendToAi(task)"
                class="px-2 py-0.5 rounded bg-purple-950/60 hover:bg-purple-800 text-purple-300 border border-purple-800/60 text-[10px] font-medium flex items-center gap-1 transition"
                title="Send task to Antigravity CLI"
              >
                <Bot class="w-3 h-3" />
                Send AI
              </button>
            </div>
          </div>

          <div 
            v-if="getColumnTasks(col.status).length === 0" 
            class="h-24 flex items-center justify-center border border-dashed border-slate-850 rounded-lg text-slate-600 text-[11px]"
          >
            Empty
          </div>
        </div>
      </div>
    </div>

    <!-- Task Detail / Edit Modal -->
    <div 
      v-if="selectedTask" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
      @click.self="selectedTask = null"
    >
      <div class="w-full max-w-lg bg-dark-900 border border-slate-700 rounded-2xl shadow-2xl overflow-hidden animate-scale-up">
        <div class="px-6 py-4 border-b border-slate-800 bg-dark-950/40 flex items-center justify-between">
          <div>
            <span class="text-[10px] uppercase font-mono text-brand-400 font-bold">{{ selectedTask.phase_name || 'Task Detail' }}</span>
            <h3 class="text-sm font-bold text-white">{{ selectedTask.title }}</h3>
          </div>
          <button @click="selectedTask = null" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="p-6 space-y-4 text-xs">
          <!-- Status Selector -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-slate-400 mb-1 font-semibold">Column / Status</label>
              <select 
                v-model="selectedTask.status"
                class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-white"
                @change="updateTaskStatus(selectedTask.id, selectedTask.status)"
              >
                <option value="backlog">Backlog</option>
                <option value="todo">Todo</option>
                <option value="in_progress">In Progress</option>
                <option value="review">Review</option>
                <option value="done">Done</option>
              </select>
            </div>
            <div>
              <label class="block text-slate-400 mb-1 font-semibold">Priority</label>
              <select 
                v-model="selectedTask.priority"
                class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-white"
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="urgent">Urgent</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-slate-400 mb-1 font-semibold">Description</label>
            <textarea 
              v-model="selectedTask.description"
              rows="3"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-white"
            ></textarea>
          </div>

          <div>
            <label class="block text-slate-400 mb-1 font-semibold">Acceptance Criteria</label>
            <textarea 
              v-model="selectedTask.acceptance_criteria"
              rows="3"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-white font-mono"
            ></textarea>
          </div>

          <div>
            <label class="block text-slate-400 mb-1 font-semibold">AI Instructions for Antigravity</label>
            <textarea 
              v-model="selectedTask.ai_instructions"
              rows="2"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-white font-mono"
            ></textarea>
          </div>
        </div>

        <div class="px-6 py-3 border-t border-slate-800 bg-dark-950/40 flex items-center justify-between">
          <button 
            @click="deleteTask(selectedTask.id)"
            class="text-rose-400 hover:text-rose-300 text-xs font-semibold"
          >
            Delete Task
          </button>

          <div class="flex items-center gap-2">
            <button 
              @click="openSendToAi(selectedTask)"
              class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-purple-600 hover:bg-purple-500 text-white shadow flex items-center gap-1"
            >
              <Bot class="w-3.5 h-3.5" />
              Send to Antigravity
            </button>
            <button 
              @click="saveTaskChanges"
              class="px-4 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow"
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Send to Antigravity Confirmation Dialog -->
    <div 
      v-if="aiTargetTask" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
      @click.self="aiTargetTask = null"
    >
      <div class="w-full max-w-lg bg-dark-900 border border-purple-500/50 rounded-2xl shadow-2xl overflow-hidden animate-scale-up">
        <div class="px-6 py-4 border-b border-slate-800 bg-purple-950/30 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Bot class="w-5 h-5 text-purple-400" />
            <h3 class="text-sm font-bold text-white">Send Task to Antigravity CLI</h3>
          </div>
          <button @click="aiTargetTask = null" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="p-6 space-y-3.5 text-xs text-slate-300">
          <p>You are about to dispatch this task to the managed Antigravity child process:</p>
          <div class="p-3 bg-dark-950 border border-slate-800 rounded-lg space-y-1">
            <div class="font-bold text-white">{{ aiTargetTask.title }}</div>
            <div class="text-[11px] text-slate-400">{{ aiTargetTask.description }}</div>
          </div>

          <div class="space-y-1 text-slate-400">
            <div class="font-semibold text-white">Injected Context Includes:</div>
            <ul class="list-disc pl-4 space-y-0.5 text-[11px]">
              <li>Approved Project PRD & Specifications</li>
              <li>Technology stack (Framework, Runtime, Database)</li>
              <li>Task Acceptance Criteria & Instructions</li>
              <li>Project working directory path</li>
            </ul>
          </div>

          <div>
            <label class="block text-slate-300 mb-1 font-semibold">Additional prompt instructions (optional)</label>
            <input 
              v-model="customAiPrompt"
              type="text" 
              placeholder="e.g. Ensure strict type validation and add unit tests"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-white"
            />
          </div>
        </div>

        <div class="px-6 py-3 border-t border-slate-800 bg-dark-950/40 flex items-center justify-end gap-2">
          <button 
            @click="aiTargetTask = null"
            class="px-3 py-1.5 rounded-lg text-xs font-semibold text-slate-400 hover:text-white"
          >
            Cancel
          </button>
          <button 
            @click="confirmSendToAi"
            class="px-4 py-1.5 rounded-lg text-xs font-semibold bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white shadow flex items-center gap-1.5"
          >
            <Play class="w-3.5 h-3.5" />
            <span>Launch AI Execution</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTaskStore } from '@/stores/tasks';
import { api } from '@/services/api';
import type { Task } from '@/types';
import { Columns3, Sparkles, Plus, Bot, Play } from 'lucide-vue-next';

const props = defineProps<{
  projectId: number | string;
}>();

const emit = defineEmits(['start-ai']);

const taskStore = useTaskStore();
const generating = ref(false);
const selectedTask = ref<Task | null>(null);
const aiTargetTask = ref<Task | null>(null);
const customAiPrompt = ref('');

const columns: { title: string; status: Task['status']; color: string }[] = [
  { title: 'Backlog', status: 'backlog', color: 'bg-slate-500' },
  { title: 'Todo', status: 'todo', color: 'bg-blue-500' },
  { title: 'In Progress', status: 'in_progress', color: 'bg-amber-500' },
  { title: 'Review', status: 'review', color: 'bg-purple-500' },
  { title: 'Done', status: 'done', color: 'bg-emerald-500' },
];

onMounted(() => {
  taskStore.fetchTasks(props.projectId);
});

function getColumnTasks(status: Task['status']) {
  return taskStore.tasks.filter(t => t.status === status);
}

function selectTask(task: Task) {
  selectedTask.value = { ...task };
}

function openCreateTaskModal() {
  const title = prompt('Enter task title:');
  if (title) {
    taskStore.createTask(props.projectId, { title });
  }
}

async function generateTasks() {
  generating.value = true;
  try {
    const res = await taskStore.generateTasksFromPRD(props.projectId);
    alert(res.message);
  } catch (err: any) {
    alert(err.response?.data?.error?.message || 'Could not generate tasks. Ensure an approved PRD exists.');
  } finally {
    generating.value = false;
  }
}

async function updateTaskStatus(taskId: number, status: Task['status']) {
  await taskStore.updateTaskStatus(taskId, status);
}

async function saveTaskChanges() {
  if (!selectedTask.value) return;
  await taskStore.updateTask(selectedTask.value.id, selectedTask.value);
  selectedTask.value = null;
}

async function deleteTask(taskId: number) {
  if (confirm('Delete this task?')) {
    await taskStore.deleteTask(taskId);
    selectedTask.value = null;
  }
}

function openSendToAi(task: Task) {
  aiTargetTask.value = task;
  customAiPrompt.value = '';
}

async function confirmSendToAi() {
  if (!aiTargetTask.value) return;
  try {
    await api.post(`/projects/${props.projectId}/antigravity/start`, {
      task_id: aiTargetTask.value.id,
      prompt: customAiPrompt.value,
    });
    const target = aiTargetTask.value;
    aiTargetTask.value = null;
    selectedTask.value = null;
    await taskStore.fetchTasks(props.projectId);
    emit('start-ai', target);
  } catch (err: any) {
    alert(err.response?.data?.error?.message || 'Failed to dispatch AI task.');
  }
}
</script>
