<template>
  <div class="space-y-4">
    <!-- Top Git Control Bar -->
    <div class="flex flex-wrap items-center justify-between gap-3 p-4 bg-dark-900 border border-slate-800 rounded-xl">
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-orange-500/10 text-orange-400 border border-orange-500/20">
          <GitBranch class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-sm text-white">Git Version Control</h3>
            <span class="text-[10px] font-mono px-2 py-0.5 rounded bg-dark-800 border border-slate-700 text-brand-400">
              {{ gitStatus?.current_branch || 'main' }}
            </span>
          </div>
          <p class="text-xs text-slate-400">
            {{ gitStatus?.total_changes || 0 }} modified files detected in working tree.
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button 
          @click="fetchGitData"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-dark-800 hover:bg-dark-750 text-slate-300 border border-slate-700 transition flex items-center gap-1.5"
        >
          <RotateCcw class="w-3.5 h-3.5" />
          <span>Refresh</span>
        </button>
      </div>
    </div>

    <!-- 2-Column Git Layout: Changes & Commit vs Diff/History -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Left Column: Changed Files & Commit Box -->
      <div class="space-y-4">
        <!-- Changed Files List -->
        <div class="bg-dark-900 border border-slate-800 rounded-xl p-4">
          <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3 flex items-center justify-between">
            <span>Changed Files</span>
            <span class="font-mono text-slate-400">{{ gitStatus?.files?.length || 0 }}</span>
          </h4>

          <div class="space-y-1 max-h-56 overflow-y-auto pr-1">
            <div 
              v-for="(f, idx) in gitStatus?.files" 
              :key="idx"
              class="px-2.5 py-1.5 rounded-lg bg-dark-850 border border-slate-800/80 flex items-center justify-between text-xs font-mono"
            >
              <span class="truncate text-slate-200" :title="f.file">{{ f.file }}</span>
              <span 
                :class="[
                  'text-[9px] uppercase font-bold px-1.5 py-0.2 rounded',
                  f.action === 'modified' ? 'text-amber-400 bg-amber-950/60' :
                  f.action === 'added' ? 'text-emerald-400 bg-emerald-950/60' :
                  f.action === 'deleted' ? 'text-rose-400 bg-rose-950/60' : 'text-blue-400 bg-blue-950/60'
                ]"
              >
                {{ f.action }}
              </span>
            </div>

            <div v-if="!gitStatus?.files?.length" class="text-center py-6 text-xs text-slate-500">
              Working tree is clean.
            </div>
          </div>
        </div>

        <!-- Commit Box -->
        <div class="bg-dark-900 border border-slate-800 rounded-xl p-4 space-y-3">
          <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400">Commit Changes</h4>
          <textarea 
            v-model="commitMessage"
            rows="2"
            placeholder="feat: implement new api endpoints and tests"
            class="w-full px-3 py-2 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white focus:outline-none focus:border-brand-500 font-mono"
          ></textarea>
          <button 
            @click="submitCommit"
            :disabled="committing || !commitMessage.trim()"
            class="w-full py-2 rounded-lg text-xs font-semibold bg-brand-600 hover:bg-brand-500 disabled:opacity-50 text-white shadow transition flex items-center justify-center gap-1.5"
          >
            <GitCommit class="w-3.5 h-3.5" />
            <span>{{ committing ? 'Committing...' : 'Stage & Commit Changes' }}</span>
          </button>
        </div>
      </div>

      <!-- Right Column: Tabs for Diff and Recent Commits -->
      <div class="lg:col-span-2 bg-dark-900 border border-slate-800 rounded-xl p-4 flex flex-col">
        <div class="flex items-center gap-2 border-b border-slate-800 pb-2 mb-3">
          <button 
            @click="tab = 'diff'"
            :class="['px-3 py-1 rounded-md text-xs font-semibold transition', tab === 'diff' ? 'bg-dark-800 text-white' : 'text-slate-400 hover:text-slate-200']"
          >
            Working Tree Diff
          </button>
          <button 
            @click="tab = 'history'"
            :class="['px-3 py-1 rounded-md text-xs font-semibold transition', tab === 'history' ? 'bg-dark-800 text-white' : 'text-slate-400 hover:text-slate-200']"
          >
            Recent Commits
          </button>
        </div>

        <!-- Diff Viewer -->
        <div v-if="tab === 'diff'" class="flex-1 bg-dark-950 border border-slate-850 rounded-lg p-3 overflow-auto max-h-[420px] font-mono text-[11px] leading-relaxed text-slate-300">
          <pre v-if="diffContent" class="whitespace-pre-wrap">{{ diffContent }}</pre>
          <div v-else class="text-center py-12 text-slate-500">
            No active diffs to display.
          </div>
        </div>

        <!-- Commit History List -->
        <div v-else class="flex-1 space-y-2 overflow-y-auto max-h-[420px]">
          <div 
            v-for="(c, idx) in commits" 
            :key="idx"
            class="p-3 bg-dark-850 border border-slate-800 rounded-lg flex items-start justify-between gap-3 text-xs"
          >
            <div>
              <div class="font-semibold text-slate-200">{{ c.message }}</div>
              <div class="text-[10px] text-slate-400 mt-1">
                {{ c.author }} • {{ c.time }}
              </div>
            </div>
            <span class="font-mono text-[10px] bg-dark-800 px-2 py-0.5 rounded text-brand-400 border border-slate-700">
              {{ c.hash }}
            </span>
          </div>

          <div v-if="commits.length === 0" class="text-center py-12 text-xs text-slate-500">
            No commits recorded yet.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { api } from '@/services/api';
import { GitBranch, RotateCcw, GitCommit } from 'lucide-vue-next';

const props = defineProps<{
  projectId: number | string;
}>();

const tab = ref<'diff' | 'history'>('diff');
const gitStatus = ref<any>(null);
const diffContent = ref('');
const commits = ref<any[]>([]);
const commitMessage = ref('');
const committing = ref(false);

onMounted(() => {
  fetchGitData();
});

async function fetchGitData() {
  try {
    const [statusRes, diffRes, commitsRes] = await Promise.all([
      api.get(`/projects/${props.projectId}/git/status`),
      api.get(`/projects/${props.projectId}/git/diff`),
      api.get(`/projects/${props.projectId}/git/commits`),
    ]);
    gitStatus.value = statusRes.data.data;
    diffContent.value = diffRes.data.data.diff;
    commits.value = commitsRes.data.data;
  } catch (e) {
    console.error(e);
  }
}

async function submitCommit() {
  if (!commitMessage.value.trim()) return;
  committing.value = true;
  try {
    await api.post(`/projects/${props.projectId}/git/commit`, {
      message: commitMessage.value,
    });
    commitMessage.value = '';
    await fetchGitData();
    alert('Changes committed successfully.');
  } catch (err: any) {
    alert(err.response?.data?.error || 'Commit failed.');
  } finally {
    committing.value = false;
  }
}
</script>
