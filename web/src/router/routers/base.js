export const baseRouters = [
    {
        path: '/',
        // name: 'Home',
        // component: () => import('../../views/system/home/index.vue'),
        // meta: {
        //     title: '首页',
        //     layout: 'empty',
        //     hidden: 1,
        //     icon: ''
        // },
        redirect: '/login'
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../../views/login/index.vue'),
        meta: {
            title: '登录',
            layout: 'empty',
            hidden: 1,
            icon: '',
        }
    },
    {
        path: '/403',
        name: '403',
        component: () => import('../../views/public/403.vue'),
        meta: {
            title: '403',
            layout: 'empty',
            hidden: 1,
            icon: '',
        }
    },
    {
        path: '/:pathMatch(.*)',
        name: '404',
        component: () => import('../../views/public/404.vue'),
        meta: {
            title: '404',
            layout: 'empty'
        }
    }
]