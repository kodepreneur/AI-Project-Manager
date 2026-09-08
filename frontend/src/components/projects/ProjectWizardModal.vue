<template>
  <div 
    v-if="isOpen" 
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
    @click.self="close"
  >
    <div class="w-full max-w-2xl bg-dark-900 border border-slate-700/80 rounded-2xl shadow-2xl overflow-hidden flex flex-col animate-scale-up">
      <!-- Header with Step Indicator -->
      <div class="px-6 py-4 border-b border-slate-800 bg-dark-950/40 flex items-center justify-between">
        <div>
          <h2 class="text-base font-bold text-white tracking-tight">Create New Project</h2>
          <p class="text-xs text-slate-400">Step {{ currentStep }} of 3: {{ stepTitles[currentStep - 1] }}</p>
        </div>
        <!-- Steps dots -->
        <div class="flex items-center gap-2">
          <div 
            v-for="step in 3" 
            :key="step"
            :class="[
              'w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-bold transition',
              currentStep === step ? 'bg-brand-600 text-white' : currentStep > step ? 'bg-emerald-600 text-white' : 'bg-dark-800 text-slate-500'
            ]"
          >
            {{ step }}
          </div>
        </div>
      </div>

      <!-- Step 1: Basic Information -->
      <div v-if="currentStep === 1" class="p-6 space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1.5">Project Name *</label>
          <input 
            v-model="form.name"
            type="text" 
            placeholder="e.g. Toko Buah ERP"
            class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
            @input="generateSlug"
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Project Slug</label>
            <input 
              v-model="form.slug"
              type="text" 
              placeholder="toko-buah-erp"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white font-mono focus:outline-none focus:border-brand-500"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Git Branch</label>
            <input 
              v-model="form.git_branch"
              type="text" 
              placeholder="main"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white font-mono focus:outline-none focus:border-brand-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1.5">Description</label>
          <textarea 
            v-model="form.description"
            rows="3" 
            placeholder="Brief summary of application purpose and key capabilities..."
            class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
          ></textarea>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1.5">Git Repository (optional)</label>
          <input 
            v-model="form.git_repo"
            type="text" 
            placeholder="https://github.com/organization/repo.git"
            class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500 font-mono"
          />
        </div>
      </div>

      <!-- Step 2: Technology Selection -->
      <div v-if="currentStep === 2" class="p-6 space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <!-- Backend -->
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Backend Framework</label>
            <select 
              v-model="form.framework"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
            >
              <option value="Laravel">Laravel (PHP)</option>
              <option value="Symfony">Symfony (PHP)</option>
              <option value="Native PHP">Native PHP</option>
              <option value="Node.js">Node.js (Express / Fastify)</option>
              <option value="Go">Go (Fiber / Gin)</option>
              <option value="Python">Python (FastAPI)</option>
            </select>
          </div>

          <!-- Frontend -->
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Frontend</label>
            <select 
              v-model="form.frontend_tech"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
            >
              <option value="Vue 3">Vue 3 + Vite</option>
              <option value="React">React + Vite</option>
              <option value="Blade">Blade / HTML</option>
              <option value="Next.js">Next.js</option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <!-- Database -->
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">Primary Database</label>
            <select 
              v-model="form.database_type"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
            >
              <option value="MySQL">MySQL</option>
              <option value="PostgreSQL">PostgreSQL</option>
              <option value="SQLite">SQLite</option>
              <option value="None">None</option>
            </select>
          </div>

          <!-- PHP Version -->
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1.5">PHP Version</label>
            <select 
              v-model="form.php_version"
              class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-sm text-white focus:outline-none focus:border-brand-500"
            >
              <option value="8.3">PHP 8.3 (Recommended)</option>
              <option value="8.4">PHP 8.4</option>
              <option value="8.2">PHP 8.2</option>
              <option value="8.1">PHP 8.1</option>
              <option value="auto">Auto-detect</option>
            </select>
          </div>
        </div>

        <!-- Project Template Presets -->
        <div class="pt-2">
          <label class="block text-xs font-semibold text-slate-400 mb-2">Quick Presets</label>
          <div class="grid grid-cols-3 gap-2">
            <button 
              type="button" 
              class="p-2 text-left rounded-lg border border-slate-800 bg-dark-850 hover:border-brand-500 transition text-xs"
              @click="applyPreset('Laravel', 'Vue 3', 'MySQL')"
            >
              <div class="font-semibold text-white">Laravel + Vue</div>
              <div class="text-[10px] text-slate-400">MySQL • PHP 8.3</div>
            </button>
            <button 
              type="button" 
              class="p-2 text-left rounded-lg border border-slate-800 bg-dark-850 hover:border-brand-500 transition text-xs"
              @click="applyPreset('Go', 'Vue 3', 'PostgreSQL')"
            >
              <div class="font-semibold text-white">Go + Vue</div>
              <div class="text-[10px] text-slate-400">PostgreSQL • Fiber</div>
            </button>
            <button 
              type="button" 
              class="p-2 text-left rounded-lg border border-slate-800 bg-dark-850 hover:border-brand-500 transition text-xs"
              @click="applyPreset('Node.js', 'React', 'SQLite')"
            >
              <div class="font-semibold text-white">Node + React</div>
              <div class="text-[10px] text-slate-400">SQLite • Express</div>
            </button>
          </div>
        </div>
      </div>

      <!-- Step 3: PRD Option -->
      <div v-if="currentStep === 3" class="p-6 space-y-4">
        <div class="text-xs text-slate-400">
          Before AI development begins, an approved PRD ensures clear milestones and acceptance criteria.
        </div>

        <div class="space-y-2.5">
          <label class="flex items-start gap-3 p-3.5 rounded-xl border border-brand-600/50 bg-brand-950/20 cursor-pointer">
            <input type="radio" v-model="form.prd_mode" value="ai" class="mt-0.5 text-brand-600" />
            <div>
              <div class="text-sm font-semibold text-white flex items-center gap-1.5">
                <Sparkles class="w-4 h-4 text-brand-400" />
                Generate PRD with AI (Recommended)
              </div>
              <p class="text-xs text-slate-400 mt-0.5">Creates a comprehensive 17-section PRD based on project objectives.</p>
            </div>
          </label>

          <label class="flex items-start gap-3 p-3.5 rounded-xl border border-slate-800 bg-dark-850 cursor-pointer">
            <input type="radio" v-model="form.prd_mode" value="manual" class="mt-0.5 text-brand-600" />
            <div>
              <div class="text-sm font-semibold text-white">Write PRD Manually</div>
              <p class="text-xs text-slate-400 mt-0.5">Start with a blank markdown document to draft your own specifications.</p>
            </div>
          </label>

          <label class="flex items-start gap-3 p-3.5 rounded-xl border border-slate-800 bg-dark-850 cursor-pointer">
            <input type="radio" v-model="form.prd_mode" value="skip" class="mt-0.5 text-brand-600" />
            <div>
              <div class="text-sm font-semibold text-white">Skip PRD for Now</div>
              <p class="text-xs text-slate-400 mt-0.5">Proceed directly to manual task creation and code development.</p>
            </div>
          </label>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="px-6 py-3 border-t border-slate-800 bg-dark-950/40 flex items-center justify-between">
        <button 
          v-if="currentStep > 1" 
          @click="currentStep--"
          class="px-4 py-2 rounded-lg text-xs font-semibold text-slate-300 hover:text-white hover:bg-slate-800 transition"
        >
          Back
        </button>
        <div v-else></div>

        <div class="flex items-center gap-2">
          <button 
            @click="close"
            class="px-4 py-2 rounded-lg text-xs font-semibold text-slate-400 hover:text-white hover:bg-slate-800 transition"
          >
            Cancel
          </button>
          
          <button 
            v-if="currentStep < 3"
            @click="nextStep"
            class="px-4 py-2 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition"
          >
            Continue
          </button>

          <button 
            v-else
            @click="submitCreate"
            :disabled="loading"
            class="px-5 py-2 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 text-white shadow transition flex items-center gap-1.5"
          >
            <span v-if="loading">Creating Project...</span>
            <span v-else>Initialize Project</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useProjectStore } from '@/stores/projects';
import { Sparkles } from 'lucide-vue-next';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits(['close']);

const router = useRouter();
const projectStore = useProjectStore();

const currentStep = ref(1);
const loading = ref(false);

const stepTitles = [
  'Basic Information',
  'Technology Stack',
  'Product Requirements Document (PRD)',
];

const form = reactive({
  name: '',
  slug: '',
  description: '',
  git_repo: '',
  git_branch: 'main',
  framework: 'Laravel',
  language: 'PHP',
  frontend_tech: 'Vue 3',
  database_type: 'MySQL',
  php_version: '8.3',
  prd_mode: 'ai',
});

function generateSlug() {
  form.slug = form.name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
}

function applyPreset(framework: string, frontend: string, db: string) {
  form.framework = framework;
  form.frontend_tech = frontend;
  form.database_type = db;
  if (framework.includes('Laravel') || framework.includes('PHP')) {
    form.language = 'PHP';
  } else if (framework === 'Go') {
    form.language = 'Go';
  } else {
    form.language = 'JavaScript';
  }
}

function nextStep() {
  if (currentStep.value === 1 && !form.name.trim()) {
    alert('Please enter a project name.');
    return;
  }
  currentStep.value++;
}

async function submitCreate() {
  loading.value = true;
  try {
    const project = await projectStore.createProject({
      name: form.name,
      slug: form.slug,
      description: form.description,
      git_repo: form.git_repo,
      git_branch: form.git_branch,
      framework: form.framework,
      language: form.language,
      frontend_tech: form.frontend_tech,
      database_type: form.database_type,
      php_version: form.php_version,
    });

    close();
    router.push(`/projects/${project.id}`);
  } catch (err: any) {
    alert(err.response?.data?.error?.message || 'Failed to create project.');
  } finally {
    loading.value = false;
  }
}

function close() {
  currentStep.value = 1;
  emit('close');
}
</script>
