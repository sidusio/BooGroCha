import Vue from 'vue'
import Router from 'vue-router'
import Book from './views/Book.vue'
import About from './views/About.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'book',
      component: Book
    },
    {
      path: '/about',
      name: 'about',
      component: About
    },
  ]
})
