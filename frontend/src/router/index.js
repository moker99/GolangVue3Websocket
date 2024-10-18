import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/Register.vue')
    },

    {
      path: '/chat',
      name: 'Chat',
      component: () => import('../views/Chat.vue')
    },
    {
      path: '/index',
      name: 'Index',
      component: () => import('../views/Index.vue'),
      children:[
        {
          path: '/friends',
          name: 'Friends',
          component: () => import('../views/Friends.vue')
        },
        {
          path: '/member',
          name: 'Member',
          component: () => import('../views/Member.vue')
        },
      ]
    },
    
  ]
})

export default router
