export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  created_at: string;
}

export interface Project {
  id: number;
  name: string;
  slug: string;
  description: string;
  path: string;
  git_repo: string;
  git_branch: string;
  framework: string;
  language: string;
  frontend_tech: string;
  database_type: string;
  php_version: string;
  status: 'Development' | 'Running' | 'Completed' | 'Paused' | 'Failed' | 'Archived';
  preview_port: number;
  preview_command: string;
  test_command: string;
  build_command: string;
  install_command: string;
  screenshot_url?: string;
  last_activity: string;
  created_at: string;
  updated_at: string;
  prd?: PRD;
  tasks?: Task[];
  processes?: ProjectProcess[];
}

export interface PRD {
  id: number;
  project_id: number;
  title: string;
  content: string;
  status: 'Draft' | 'Review' | 'Approved' | 'Archived';
  version: number;
  approved_at?: string;
  created_at: string;
  updated_at: string;
  versions?: PRDVersion[];
}

export interface PRDVersion {
  id: number;
  prd_id: number;
  version: number;
  content: string;
  summary: string;
  created_at: string;
}

export interface Task {
  id: number;
  project_id: number;
  phase_id?: number;
  phase_name: string;
  title: string;
  description: string;
  priority: 'low' | 'medium' | 'high' | 'urgent';
  complexity: 'low' | 'medium' | 'high';
  status: 'backlog' | 'todo' | 'in_progress' | 'review' | 'done' | 'cancelled';
  acceptance_criteria: string;
  ai_instructions: string;
  dependencies: string;
  files_affected: string;
  order: number;
  ai_status: 'idle' | 'running' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
}

export interface AntigravitySession {
  id: number;
  project_id: number;
  task_id?: number;
  pid: number;
  status: 'idle' | 'running' | 'paused' | 'completed' | 'failed';
  command: string;
  stdout: string;
  stderr: string;
  exit_code: number;
  start_time: string;
  end_time?: string;
  current_task: string;
}

export interface ProjectProcess {
  id: number;
  project_id: number;
  name: string;
  command: string;
  pid: number;
  port: number;
  status: 'running' | 'stopped' | 'failed';
  auto_restart: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProjectLog {
  id: number;
  project_id?: number;
  category: string;
  level: string;
  message: string;
  metadata?: string;
  created_at: string;
}

export interface SystemMetrics {
  cpu_percent: number;
  num_cpu: number;
  mem_total_mb: number;
  mem_used_mb: number;
  mem_percent: number;
  disk_total_gb: number;
  disk_used_gb: number;
  disk_percent: number;
  load_1: number;
  load_5: number;
  load_15: number;
  process_count: number;
  os: string;
  arch: string;
  server_time: string;
}
