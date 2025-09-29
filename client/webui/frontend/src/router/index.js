import {createRouter, createWebHistory} from "vue-router";
import {baseRouters} from "./routers/base.js";

const router = createRouter({
    history: createWebHistory('/extras/gostc/'),
    routes: baseRouters,
    scrollBehavior: () => ({left: 0, top: 0}),
})

router.beforeEach(async (to, from, next) => {
    next()
})

router.afterEach((to) => {
})

export default router

