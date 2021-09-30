import { nextTick } from 'vue'
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { LoadingBarApi } from 'naive-ui'
import ForbiddenPage from '../pages/403.vue'
import NotFoundPage from '../pages/404.vue'
import ErrorPage from '../pages/500.vue'
import LoginPage from '../pages/Login.vue'
import InitPage from '../pages/Init.vue'
import { store } from "../store";

var loadingBar: LoadingBarApi;

export function initLoadingBar(bar: LoadingBarApi) {
  loadingBar = bar
}

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: () => import('../pages/Home.vue'),
    meta: {
      title: '首页',
    }
  },
  {
    name: 'login',
    path: '/login',
    component: LoginPage,
    meta: {
      layout: "empty",
      anonymous: true,
      title: '登录',
    }
  },
  {
    name: 'init',
    path: '/init',
    component: InitPage,
    meta: {
      layout: "empty",
      anonymous: true,
      title: '初始化',
    }
  },
  {
    path: "/profile",
    component: () => import('../pages/Profile.vue'),
    meta: {
      title: '个人资料',
    }
  },
  {
    path: "/config/notice",
    component: () => import('../pages/config/Notice.vue'),
    meta: {
      title: '通知设置',
    }
  },
  {
    path: "/config/advance",
    component: () => import('../pages/config/Advance.vue'),
    meta: {
      title: '高级设置',
    }
  },
  {
    path: "/tasks",
    component: () => import('../pages/task/List.vue'),
    meta: {
      title: '任务列表',
    }
  },
  {
    path: "/tasks/new",
    component: () => import('../pages/task/Edit.vue'),
    meta: {
      title: '新建任务',
    }
  },
  {
    path: "/tasks/:name",
    component: () => import('../pages/task/View.vue'),
    meta: {
      title: '任务详情',
    }
  },
  {
    path: "/tasks/:name/edit",
    component: () => import('../pages/task/Edit.vue'),
    meta: {
      title: '任务编辑',
    }
  },
  {
    path: "/jobs",
    component: () => import('../pages/job/List.vue'),
    meta: {
      title: '作业管理',
    }
  },
  {
    path: "/jobs/:id",
    component: () => import('../pages/job/View.vue'),
    meta: {
      title: '作业详情',
    }
  },
  {
    path: "/account/users",
    component: () => import('../pages/account/user/List.vue'),
    meta: {
      title: '用户管理',
    }
  },
  {
    path: "/account/users/new",
    component: () => import('../pages/account/user/Edit.vue'),
    meta: {
      title: '新建用户',
    }
  },
  {
    path: "/account/users/:id",
    component: () => import('../pages/account/user/View.vue'),
    meta: {
      title: '用户详情',
    }
  },
  {
    path: "/account/users/:id/edit",
    component: () => import('../pages/account/user/Edit.vue'),
    meta: {
      title: '用户编辑',
    }
  },
  {
    name: "role.list",
    path: "/account/roles",
    component: () => import('../pages/account/role/List.vue'),
    meta: {
      title: '角色管理',
    }
  },
  {
    name: "role.new",
    path: "/account/roles/new",
    component: () => import('../pages/account/role/Edit.vue'),
    meta: {
      title: '新建角色',
    }
  },
  {
    name: "role.view",
    path: "/account/roles/:id",
    component: () => import('../pages/account/role/View.vue'),
    meta: {
      title: '角色详情',
    }
  },
  {
    name: "role.edit",
    path: "/account/roles/:id/edit",
    component: () => import('../pages/account/role/Edit.vue'),
    meta: {
      title: '角色编辑',
    }
  },
  {
    name: '403',
    path: '/403',
    component: ForbiddenPage,
    meta: {
      layout: "simple",
      anonymous: true,
      title: '禁止访问',
    }
  },
  {
    name: '404',
    path: '/404',
    component: NotFoundPage,
    meta: {
      layout: "simple",
      anonymous: true,
      title: '页面不存在',
    }
  },
  {
    name: '500',
    path: '/500',
    component: ErrorPage,
    meta: {
      layout: "simple",
      anonymous: true,
      title: '服务器错误',
    }
  },
  {
    name: 'not-found',
    path: '/:pathMatch(.*)*',
    redirect: { name: '404' }
  },
]

function createSiteRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes,
  })

  router.beforeEach(function (to, from, next) {
    if (!from || to.path !== from.path) {
      loadingBar?.start()
      window.document.title = to.meta.title ? `${to.meta.title} - Skynet` : 'Skynet'
    }

    if (to.matched.some(record => !record.meta.anonymous)) {
      // this route requires auth, if not logged in, redirect to login page.      
      if (store.getters.anonymous) {
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
        return
      }
    }

    next()
  })

  router.afterEach(function (to, from) {
    if (!from || to.path !== from.path) {
        loadingBar?.finish()
      if (to.hash && to.hash !== from.hash) {
        nextTick(() => {
          const el = document.querySelector(to.hash)
          if (el) el.scrollIntoView()
        })
      }
    }
  })

  return router
}

export const router = createSiteRouter()