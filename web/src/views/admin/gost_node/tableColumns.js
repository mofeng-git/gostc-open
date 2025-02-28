import {h, ref} from "vue";
import {NButton, NSpace} from "naive-ui";
import router from "../../../router/index.js";
import Online from "../../../icon/online.vue";

export function generateTableColumns(
    update = (row) => {
    },
    remove = (row) => {
    },
    bind = (row)=>{
    },
) {
    return [
        {
            title: '名称',
            key: 'name',
            ellipsis: {tooltip: true},
            render(row) {
                return h('div', {
                    style: {
                        display: 'flex',
                        alignItems: 'center',
                        cursor: 'pointer',
                    },
                    onClick:()=>{
                        bind(row)
                    }
                }, {
                    default: () => [
                        h(Online, {online: row.online === 1, size: 10}, {}),
                        h('span', {
                            style: {
                                marginLeft: '8px',
                            }
                        }, {default: () => row.name}),
                    ]
                });
            }
        },
        {title: '连接密钥', key: 'key', ellipsis: {tooltip: true}, width: 330},
        {title: '地址', key: 'address', ellipsis: {tooltip: true}},
        {
            title: '功能',
            key: 'address',
            ellipsis: {tooltip: true},
            render(row) {
                let data = []
                if (row.tunnel === 1 && row.web === 1) {
                    data.push('域名解析')
                }
                if (row.forward === 1) {
                    data.push('端口转发')
                }
                if (row.tunnel === 1) {
                    data.push('私有隧道')
                }
                return data?.join(',');
            }
        },
        {
            title: '规则',
            key: 'rules',
            ellipsis: {tooltip: true},
            render(row) {
                return row.ruleNames?.join(',')
            }
        },
        {
            title: '标签',
            key: 'tags',
            ellipsis: {tooltip: true},
            render(row) {
                return row.tags?.join(',')
            }
        },
        {title: '介绍', key: 'remark', ellipsis: {tooltip: true}},
        {
            title: '操作',
            key: 'operator',
            width: 4*67,
            // fixed: 'right',
            render(row) {
                return generateButtonGroup(
                    generateObsButton(row),
                    generateLoggerButton(row),
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


function generateLoggerButton(row) {
    return h(NButton, {
        size: 'small',
        type: 'info',
        onClick() {
            router.push({name: 'AdminGostNodeLogger', query: {nodeCode: row.code}})
        },
    }, {
        default: () => '日志',
    })
}

function generateObsButton(row) {
    return h(NButton, {
        size: 'small',
        type: 'info',
        onClick() {
            router.push({name: 'NormalGostNodeObs', query: {nodeCode: row.code}})
        },
    }, {
        default: () => '流量',
    })
}

