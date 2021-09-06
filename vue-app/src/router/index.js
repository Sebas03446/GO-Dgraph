import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import vi from '../views/vi.vue'
import About from '../views/About.vue'
import SyncClients from '../views/SyncClients.vue'
//import { component } from 'vue/types/umd'
Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: About
  },
  {
    path: '/v1',
    name: 'vi',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: vi
  },
  {
    path:'/v1/syncData',
    name: 'SyncClients',
    component: SyncClients
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
  linkActiveClass:"active",
  linkExactActiveClass:"exact-active"
})

export default router
