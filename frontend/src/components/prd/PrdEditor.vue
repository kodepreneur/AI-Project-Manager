<template>
  <div class="space-y-4">
    <!-- PRD Top Bar -->
    <div class="flex flex-wrap items-center justify-between gap-3 p-4 bg-dark-900 border border-slate-800 rounded-xl">
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-blue-500/10 text-blue-400 border border-blue-500/20">
          <FileText class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-sm text-white">{{ prd?.title || 'Product Requirements Document' }}</h3>
            <span 
              :class="[
                'text-[10px] font-semibold px-2 py-0.5 rounded-full border',
                prd?.status === 'Approved' 
                  ? 'bg-emerald-950/70 text-emerald-400 border-emerald-800'
                  : 'bg-amber-950/70 text-amber-400 border-amber-800'
              ]"
            >
              {{ prd?.status || 'Draft' }}
            </span>
            <span v-if="prd?.version" class="text-[10px] font-mono text-slate-400">v{{ prd.version }}</span>
          </div>
          <p class="text-xs text-slate-400">Defines architecture, user flows, and development milestones.</p>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2">
        <button 
          @click="showAiModal = true"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Sparkles class="w-3.5 h-3.5" />
          <span>Generate with AI</span>
        </button>

        <button 
          v-if="prd?.status !== 'Approved'"
          @click="approvePRD"
          :disabled="!prd?.content"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white shadow transition flex items-center gap-1.5"
        >
          <CheckCircle2 class="w-3.5 h-3.5" />
          <span>Approve PRD</span>
        </button>

        <button 
          @click="savePRD"
          :disabled="saving"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
        >
          <Save class="w-3.5 h-3.5" />
          <span>{{ saving ? 'Saving...' : 'Save Draft' }}</span>
        </button>

        <button 
          @click="exportMarkdown"
          class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800 border border-slate-800 transition"
          title="Export Markdown"
        >
          <Download class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Mode Toggle: Edit / Preview -->
    <div class="flex items-center justify-between border-b border-slate-800 pb-2">
      <div class="flex items-center gap-2 text-xs">
        <button 
          @click="activeTab = 'edit'"
          :class="['px-3 py-1 rounded-md font-medium transition', activeTab === 'edit' ? 'bg-dark-800 text-white' : 'text-slate-400 hover:text-slate-200']"
        >
          Editor
        </button>
        <button 
          @click="activeTab = 'preview'"
          :class="['px-3 py-1 rounded-md font-medium transition', activeTab === 'preview' ? 'bg-dark-800 text-white' : 'text-slate-400 hover:text-slate-200']"
        >
          Rendered Preview
        </button>
      </div>

      <span class="text-[11px] text-slate-500 font-mono">17 structured PRD sections supported</span>
    </div>

    <!-- Editor Mode -->
    <div v-if="activeTab === 'edit'" class="relative">
      <textarea 
        v-model="editorContent"
        rows="26"
        placeholder="# Product Requirements Document&#10;&#10;Click 'Generate with AI' to build a structured 17-section PRD automatically..."
        class="w-full p-4 bg-dark-900 border border-slate-800 rounded-xl font-mono text-xs leading-relaxed text-slate-200 focus:outline-none focus:border-brand-500/50"
      ></textarea>
    </div>

    <!-- Rendered Preview Mode -->
    <div 
      v-else 
      class="p-6 bg-dark-900 border border-slate-800 rounded-xl min-h-[500px] prose prose-invert max-w-none text-xs leading-relaxed"
    >
      <pre class="whitespace-pre-wrap font-sans text-slate-200 bg-transparent p-0 border-0">{{ editorContent || 'No PRD content yet.' }}</pre>
    </div>

    <!-- AI PRD Generator Modal -->
    <div 
      v-if="showAiModal" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
      @click.self="showAiModal = false"
    >
      <div class="w-full max-w-xl bg-dark-900 border border-slate-700 rounded-2xl shadow-2xl overflow-hidden animate-scale-up">
        <div class="px-6 py-4 border-b border-slate-800 bg-dark-950/40 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Sparkles class="w-5 h-5 text-purple-400" />
            <h3 class="text-sm font-bold text-white">AI PRD Generator</h3>
          </div>
          <button @click="showAiModal = false" class="text-slate-400 hover:text-white">✕</button>
        </div>

        <div class="p-6 space-y-3.5">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Project Idea *</label>
            <textarea 
              v-model="aiForm.project_idea"
              rows="2"
              placeholder="e.g. Modern POS and inventory management system with offline caching"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Target Users</label>
              <input 
                v-model="aiForm.target_users"
                type="text" 
                placeholder="Cashiers, Store Managers"
                class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Business Goals</label>
              <input 
                v-model="aiForm.business_goals"
                type="text" 
                placeholder="Speed up checkout by 50%"
                class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Main Features (one per line)</label>
            <textarea 
              v-model="aiForm.main_features"
              rows="3"
              placeholder="Product catalog and barcode scanner&#10;Order cart & multi-payment gateway&#10;Receipt generation & email invoice&#10;Inventory deduction on checkout"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
            ></textarea>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Technical Requirements (optional)</label>
            <input 
              v-model="aiForm.technical_requirements"
              type="text" 
              placeholder="REST API, SQLite/MySQL, WebSocket real-time updates"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500"
            />
          </div>
        </div>

        <div class="px-6 py-3 border-t border-slate-800 bg-dark-950/40 flex items-center justify-end gap-2">
          <button 
            @click="showAiModal = false"
            class="px-3 py-1.5 rounded-lg text-xs font-semibold text-slate-400 hover:text-white"
          >
            Cancel
          </button>
          <button 
            @click="generatePRD"
            :disabled="generating"
            class="px-4 py-1.5 rounded-lg text-xs font-semibold bg-purple-600 hover:bg-purple-500 text-white shadow transition flex items-center gap-1.5"
          >
            <Sparkles class="w-3.5 h-3.5" />
            <span>{{ generating ? 'Generating 17 Sections...' : 'Generate PRD' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue';
import { api } from '@/services/api';
import type { PRD } from '@/types';
import { FileText, Sparkles, CheckCircle2, Save, Download } from 'lucide-vue-next';

const props = defineProps<{
  projectId: number | string;
}>();

const emit = defineEmits(['approved']);

const prd = ref<PRD | null>(null);
const editorContent = ref('');
const activeTab = ref<'edit' | 'preview'>('edit');
const saving = ref(false);
const showAiModal = ref(false);
const generating = ref(false);

const aiForm = reactive({
  project_idea: '',
  target_users: '',
  business_goals: '',
  main_features: '',
  technical_requirements: '',
  additional_notes: '',
});

onMounted(() => {
  loadPRD();
});

watch(() => props.projectId, () => {
  loadPRD();
});

async function loadPRD() {
  try {
    const res = await api.get(`/projects/${props.projectId}/prd`);
    prd.value = res.data.data;
    editorContent.value = prd.value?.content || '';
  } catch {
    prd.value = null;
    editorContent.value = '';
  }
}

async function savePRD() {
  saving.value = true;
  try {
    const res = await api.post(`/projects/${props.projectId}/prd`, {
      title: prd.value?.title || 'Product Requirements Document',
      content: editorContent.value,
      status: prd.value?.status || 'Draft',
    });
    prd.value = res.data.data;
    alert('PRD saved successfully.');
  } finally {
    saving.value = false;
  }
}

async function approvePRD() {
  if (!confirm('Approve this PRD? Once approved, you can generate development tasks for Antigravity.')) return;
  try {
    const res = await api.post(`/projects/${props.projectId}/prd/approve`);
    prd.value = res.data.data;
    emit('approved');
    alert('PRD approved! You can now generate development tasks in the TODO board.');
  } catch (err: any) {
    alert(err.response?.data?.error?.message || 'Failed to approve PRD.');
  }
}

async function generatePRD() {
  if (!aiForm.project_idea.trim()) {
    alert('Please enter a project idea.');
    return;
  }
  generating.value = true;
  try {
    const res = await api.post(`/projects/${props.projectId}/prd/generate`, aiForm);
    editorContent.value = res.data.data.content;
    showAiModal.value = false;
    // Auto-save as draft
    await savePRD();
  } catch (err: any) {
    alert(err.response?.data?.error?.message || 'Failed to generate PRD.');
  } finally {
    generating.value = false;
  }
}

function exportMarkdown() {
  window.open(`/api/projects/${props.projectId}/prd/export`, '_blank');
}
</script>
