import Vue from 'vue'
import VueRouter from 'vue-router'

import Device from '@/views/Devices.vue'
import Votes from '@/views/Votes.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: Device
  },
  {
    path: '/votes',
    name: 'votes',
    component: Votes
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
