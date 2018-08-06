import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'landing-page',
      component: require('@/pages/login/login').default
    },
    {
      path:"/index",
      name:"Index",
      component:require('@/pages/index/index').default
    },
    {
      path: '*',
      redirect: '/'
    }
  ]
})
