<template>
  <div class="space-y-4">
    <!-- Top Database Header -->
    <div class="flex flex-wrap items-center justify-between gap-3 p-4 bg-dark-900 border border-slate-800 rounded-xl">
      <div class="flex items-center gap-3">
        <div class="p-2 rounded-lg bg-amber-500/10 text-amber-400 border border-amber-500/20">
          <Database class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-sm text-white">Database Control</h3>
            <span class="text-[10px] font-mono px-2 py-0.5 rounded bg-dark-800 border border-slate-700 text-amber-400">
              {{ engine }}
            </span>
          </div>
          <p class="text-xs text-slate-400">Manage database connections, check statuses, and execute queries.</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button 
          @click="testConnection"
          :disabled="testing"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-dark-800 hover:bg-dark-750 text-slate-300 border border-slate-700 transition flex items-center gap-1.5"
        >
          <Activity class="w-3.5 h-3.5 text-emerald-400" />
          <span>{{ testing ? 'Testing...' : 'Test Connection' }}</span>
        </button>
      </div>
    </div>

    <!-- Connection Config and SQL Query Console -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Connection Form -->
      <div class="bg-dark-900 border border-slate-800 rounded-xl p-4 space-y-3">
        <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400">Connection Settings</h4>

        <div>
          <label class="block text-[11px] font-semibold text-slate-400 mb-1">Host</label>
          <input 
            v-model="conn.host"
            type="text" 
            class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
          />
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="block text-[11px] font-semibold text-slate-400 mb-1">Port</label>
            <input 
              v-model="conn.port"
              type="number" 
              class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
            />
          </div>
          <div>
            <label class="block text-[11px] font-semibold text-slate-400 mb-1">Database</label>
            <input 
              v-model="conn.database"
              type="text" 
              class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
            />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="block text-[11px] font-semibold text-slate-400 mb-1">User</label>
            <input 
              v-model="conn.user"
              type="text" 
              class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
            />
          </div>
          <div>
            <label class="block text-[11px] font-semibold text-slate-400 mb-1">Password</label>
            <input 
              v-model="conn.password"
              type="password" 
              class="w-full px-2.5 py-1.5 bg-dark-850 border border-slate-700 rounded-lg text-xs text-white font-mono"
            />
          </div>
        </div>

        <div v-if="testResult" :class="['p-2.5 rounded-lg text-xs border', testResult.success ? 'bg-emerald-950/40 text-emerald-300 border-emerald-800' : 'bg-rose-950/40 text-rose-300 border-rose-800']">
          {{ testResult.message }}
        </div>
      </div>

      <!-- SQL Console -->
      <div class="lg:col-span-2 bg-dark-900 border border-slate-800 rounded-xl p-4 flex flex-col space-y-3">
        <div class="flex items-center justify-between">
          <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center gap-1.5">
            <Terminal class="w-3.5 h-3.5 text-brand-400" />
            <span>SQL Query Console</span>
          </h4>

          <button 
            @click="executeQuery"
            :disabled="runningQuery || !sqlQuery.trim()"
            class="px-3 py-1 rounded-md text-xs font-semibold bg-brand-600 hover:bg-brand-500 disabled:opacity-50 text-white shadow transition flex items-center gap-1"
          >
            <Play class="w-3 h-3" />
            <span>{{ runningQuery ? 'Running...' : 'Execute SQL' }}</span>
          </button>
        </div>

        <textarea 
          v-model="sqlQuery"
          rows="4"
          placeholder="SELECT * FROM users LIMIT 10;"
          class="w-full px-3 py-2 bg-dark-950 border border-slate-800 rounded-lg text-xs text-brand-300 font-mono focus:outline-none focus:border-brand-500"
        ></textarea>

        <!-- Query Output / Results Table -->
        <div class="flex-1 min-h-[160px] bg-dark-950 border border-slate-850 rounded-lg p-3 overflow-auto font-mono text-xs">
          <div v-if="queryResult">
            <div class="text-[10px] text-slate-500 mb-2">
              Execution time: {{ queryResult.execution_time }} • Affected: {{ queryResult.affected_rows }}
            </div>
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="border-b border-slate-800 text-slate-400">
                  <th v-for="col in queryResult.columns" :key="col" class="py-1 px-2 font-semibold">
                    {{ col }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, rIdx) in queryResult.rows" :key="rIdx" class="border-b border-slate-900 text-slate-300">
                  <td v-for="(val, cIdx) in row" :key="cIdx" class="py-1 px-2">
                    {{ val }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div v-else class="text-center py-10 text-slate-500 text-xs">
            Ready to execute SQL queries.
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { api } from '@/services/api';
import { Database, Activity, Terminal, Play } from 'lucide-vue-next';

const props = withDefaults(defineProps<{
  engine?: string;
  projectId?: number | string;
}>(), {
  engine: 'MySQL',
});

const conn = reactive({
  host: '127.0.0.1',
  port: props.engine === 'PostgreSQL' ? 5432 : 3306,
  database: 'app_db',
  user: 'root',
  password: '',
});

const testing = ref(false);
const testResult = ref<any>(null);
const sqlQuery = ref('SHOW TABLES;');
const runningQuery = ref(false);
const queryResult = ref<any>(null);

async function testConnection() {
  testing.value = true;
  testResult.value = null;
  const endpoint = props.engine === 'PostgreSQL' ? '/postgres/test' : '/mysql/test';
  try {
    const res = await api.post(endpoint, conn);
    testResult.value = res.data;
  } catch (err: any) {
    testResult.value = {
      success: false,
      message: err.response?.data?.message || 'Connection test failed.',
    };
  } finally {
    testing.value = false;
  }
}

async function executeQuery() {
  if (!sqlQuery.value.trim()) return;
  runningQuery.value = true;
  const endpoint = props.engine === 'PostgreSQL' ? '/postgres/query' : '/mysql/query';
  try {
    const res = await api.post(endpoint, {
      ...conn,
      query: sqlQuery.value,
    });
    queryResult.value = res.data.data;
  } catch (err: any) {
    alert(err.response?.data?.error?.message || 'Query execution failed.');
  } finally {
    runningQuery.value = false;
  }
}
</script>
