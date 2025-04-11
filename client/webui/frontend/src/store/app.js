import {defineStore} from "pinia";

export const appStore = defineStore('app', {
    state: () => ({
        siteConfig: {
            title: 'GOSTC',
            favicon: '',
            isLoading: false,
        },
        width: window.innerWidth,
        height: window.innerHeight,
        auth: {
            token: '',
            expAt: 0,
        },
        userInfo: {
            account: '',
            amount: '',
            admin: 2,
            isLoading: false,
        }
    }),
    getters: {
        drawerWidthAdapter: (state) => {
            if (state?.width < 500) {
                return state?.width
            }
            return 400
        }
    },
    actions: {},
    persist: false
})