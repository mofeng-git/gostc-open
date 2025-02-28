import {createRouter, createWebHistory} from "vue-router";
import {baseRouters} from "./routers/base.js";
import {normalRouters} from "./routers/normal.js";
import {apiAuthUserInfo} from "../api/auth/index.js";
import {appStore} from "../store/app.js";
import {localStore} from "../store/local.js";
import {apiSystemConfigQuery} from "../api/public/index.js";
import {adminRouters} from "./routers/admin.js";

export const allRouters = baseRouters.concat(normalRouters, adminRouters)

const router = createRouter({
    history: createWebHistory(),
    routes: allRouters,
    scrollBehavior: () => ({left: 0, top: 0}),
})

function setFavicon(url) {
    let favicon = document.querySelector('[rel="icon"]')
    favicon.setAttribute('href', url)
}

const initSiteConfig = async () => {
    if (appStore().siteConfig.isLoading) {
        return
    }
    let res = await apiSystemConfigQuery()
    if (res.data.favicon) {
        setFavicon(res.data.favicon)
    }
    appStore().siteConfig = res.data
    appStore().siteConfig.isLoading = true
}

const initUserInfo = async () => {
    if (appStore().userInfo.isLoading) {
        return
    }
    let res = await apiAuthUserInfo()
    appStore().userInfo = res.data
    appStore().userInfo.isLoading = true
}

// 公开访问的路由
const publicRouterName = baseRouters.map(item => item.name).filter(item => {
    return item !== 'Login'
})
// 登录用户可访问的路由
const normalRouterName = getRouterAllNames(normalRouters)
// 管理员可访问的路由
const adminRouterName = getRouterAllNames(adminRouters)

function getRouterAllNames(r) {
    let result = []
    if (!r) {
        return result
    }
    r.forEach(item => {
        result.push(item.name)
        if (item?.children) {
            result.push(...getRouterAllNames(item.children))
        }
    })
    return result
}

router.beforeEach(async (to, from, next) => {
    await initSiteConfig()
    document.title = appStore().siteConfig.title + ' - ' + to.meta?.title
    if (publicRouterName.indexOf(to.name) >= 0) {
        next()
        return
    }
    if (!localStore().auth.token) {
        // 未登录
        if (to.name === 'Login') {
            next()
        } else {
            next('/login')
        }
    } else {
        await initUserInfo()
        // 已登录
        if (to.name === 'Login') {
            next({name: normalRouterName[0]})
        } else {
            if (normalRouterName.indexOf(to.name) >= 0) {
                next()
            } else if (adminRouterName.indexOf(to.name) >= 0 && appStore().userInfo?.admin === 1) {
                next()
            } else {
                next({name: '403'})
            }
            // // if (appStore().userInfo?.admin===1){
            // //
            // // }
            // next()
        }
    }
    localStore().menuKey = to.name
})

router.afterEach((to) => {
})

export default router

