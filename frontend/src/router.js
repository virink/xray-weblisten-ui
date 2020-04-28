import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

// @ts-ignore
import Layout from '@/components/Layout'

export const constantRoutes = [{
    path: '/',
    component: Layout,
    redirect: '/projects', // Index
    children: [{
        path: '/projects',
        // @ts-ignore
        component: () => import('@/components/Projects'),
        name: 'Projects',
        meta: {
            title: 'Projects'
        }
    }, {
        path: '/vuls/:id',
        // @ts-ignore
        component: () => import('@/components/Vuls'),
        name: 'Vuls',
        meta: {
            title: 'Vuls'
        }
    }]
}]

const createRouter = () => new Router({
    mode: 'history', // require service support
    routes: constantRoutes
})

const router = createRouter()

export default router