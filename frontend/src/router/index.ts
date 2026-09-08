import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '/setup',
          name: 'Setup',
          component: () => import('@/pages/SetupPage.vue'),
          meta: { title: 'First-time Setup' },
        },
        {
          path: '/login',
          name: 'Login',
          component: () => import('@/pages/LoginPage.vue'),
          meta: { title: 'Login' },
        },
      ],
    },
    {
      path: '/',
      component: () => import('@/layouts/AppLayout.vue'),
      children: [
        {
          path: '',
          redirect: '/dashboard',
        },
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/pages/DashboardPage.vue'),
          meta: { title: 'Dashboard' },
        },
        {
          path: 'projects',
          name: 'Projects',
          component: () => import('@/pages/ProjectsPage.vue'),
          meta: { title: 'Projects' },
        },
        {
          path: 'projects/:id',
          name: 'ProjectDetail',
          component: () => import('@/pages/ProjectDetailPage.vue'),
          meta: { title: 'Project Control' },
        },
        {
          path: 'kanban',
          name: 'Kanban',
          component: () => import('@/pages/KanbanPage.vue'),
          meta: { title: 'TODO Board' },
        },
        {
          path: 'prds',
          name: 'PRDs',
          component: () => import('@/pages/PrdListPage.vue'),
          meta: { title: 'PRDs & Specifications' },
        },
        {
          path: 'antigravity',
          name: 'Antigravity',
          component: () => import('@/pages/AntigravityPage.vue'),
          meta: { title: 'Antigravity Sessions' },
        },
        {
          path: 'infrastructure/:type?',
          name: 'Infrastructure',
          component: () => import('@/pages/InfrastructurePage.vue'),
          meta: { title: 'Infrastructure' },
        },
        {
          path: 'logs',
          name: 'Logs',
          component: () => import('@/pages/LogsPage.vue'),
          meta: { title: 'System Logs' },
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('@/pages/SettingsPage.vue'),
          meta: { title: 'Settings' },
        },
      ],
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  const isSetup = await authStore.checkSetupStatus();

  if (!isSetup && to.path !== '/setup') {
    return next('/setup');
  }

  if (isSetup && to.path === '/setup') {
    return next('/login');
  }

  const token = localStorage.getItem('token');
  if (!token && to.path !== '/login' && to.path !== '/setup') {
    return next('/login');
  }

  if (token && (to.path === '/login' || to.path === '/setup')) {
    return next('/dashboard');
  }

  next();
});

export default router;
