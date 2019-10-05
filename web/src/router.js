import Vue from 'vue'
import Router from 'vue-router'
import Book from './views/Book.vue'
import Settings from './views/Settings.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Book',
      component: Book,
      meta: {
        requiresCredentials: true,
        priority: 1,
        icon: 'mdi-plus',
      },
    },
    {
      path: '/settings',
      name: 'Settings',
      component: Settings,
      meta: {
        requiresCredentials: false,
        priority: 10,
        icon: 'mdi-settings',
      },
    },
  ],
})
