import {h, ref} from "vue";
import {NButton, NSpace} from "naive-ui";
import moment from "moment";
import {cLimiterText, limiterText, rLimiterText} from "../gost_node_config/index.js";

export function generateTableColumns(
    update = (row) => {
    },
    remove = (row) => {
    },
) {
    return [
        {title: '用户账号', key: 'userAccount', ellipsis: {tooltip: true}},
        {title: '名称', key: 'name', ellipsis: {tooltip: true}},
        {
            title: '续费方式',
            key: 'chargingType',
            ellipsis: {tooltip: true},
            render(row) {
                switch (row.chargingType) {
                    case 1:
                        return `一次性计费,${row.amount}积分`
                    case 2:
                        return `循环计费,${row.amount}积分/${row.cycle}天`
                    case 3:
                        return '免费'
                }
            }
        },
        {
            title: '速率',
            key: 'limiter',
            ellipsis: {tooltip: true},
            render(row) {
                return limiterText(row.limiter);
            }
        },
        {
            title: '并发数',
            key: 'rLimiter',
            ellipsis: {tooltip: true},
            render(row) {
                return rLimiterText(row.rLimiter);
            }
        },
        {
            title: '连接数',
            key: 'cLimiter',
            ellipsis: {tooltip: true},
            render(row) {
                return cLimiterText(row.cLimiter);
            }
        },
        {
            title: '到期时间',
            key: 'expAt',
            ellipsis: {tooltip: true},
            width: 200,
            render: (row) => {
                if (row.chargingType === 2) {
                    return moment(new Date(row.expAt * 1000)).format('YYYY-MM-DD HH:mm:ss')
                }
                return '-'
            }
        },
        {
            title: '使用情况',
            key: 'tunnelType',
            ellipsis: {tooltip: true},
            render(row) {
                switch (row.tunnelType) {
                    case 1:
                        return `域名解析`
                    case 2:
                        return '端口转发'
                    case 3:
                        return '私有隧道'
                }
                return '未使用'
            }
        },
        {
            title: '操作',
            key: 'operator',
            width: 134,
            fixed: 'right',
            render(row) {
                return generateButtonGroup(
                    generateUpdateButton(row, update),
                    generateDeleteButton(row, remove),
                )
            }
        },
    ]
}

function generateButtonGroup(...btn) {
    return h(NSpace, {}, {
        default: () => [...btn]
    })
}

function generateUpdateButton(row, event) {
    return h(NButton, {
        size: 'small',
        type: 'success',
        onClick() {
            event(row)
        },
    }, {
        default: () => '修改',
    })
}


function generateDeleteButton(row, event) {
    return h(NButton, {
        size: 'small',
        type: 'error',
        onClick() {
            const loading = ref(false)
            $dialog.warning({
                maskClosable: true,
                closable: false,
                showIcon: true,
                title: '删除',
                content: '确认删除吗?',
                negativeText: '取消',
                positiveText: '确认',
                negativeButtonProps: {
                    focusable: false,
                },
                loading: loading,
                positiveButtonProps: {
                    type: 'error',
                    focusable: false,
                },
                async onNegativeClick() {

                },
                async onPositiveClick() {
                    try {
                        loading.value = true
                        await event(row)
                    } finally {
                        loading.value = false
                    }
                }
            })
        },
    }, {
        default: () => '删除',
    })
}

