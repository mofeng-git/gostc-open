import {h, ref} from "vue";
import {NButton, NSpace, NTag} from "naive-ui";
import {flowFormat} from "../../../utils/flow.js";

export function generateTableColumns(
    update = (row) => {
    },
    remove = (row) => {
    },
) {
    return [
        {title: '账号', key: 'account', ellipsis: {tooltip: true}},
        {
            title: '角色',
            key: 'admin',
            ellipsis: {tooltip: true},
            render(row) {
                return generateAdminTag(row)
            }
        },
        {title: '积分', key: 'amount', ellipsis: {tooltip: true}, width: 190},
        {
            title: '流量(IN|OUT)',
            ellipsis: {tooltip: true},
            width: 190, render(row) {
                return flowFormat(row.inputBytes) + ' | ' + flowFormat(row.outputBytes);
            }
        },
        {title: '注册时间', key: 'createdAt', ellipsis: {tooltip: true}, width: 190},
        {
            title: '操作',
            key: 'operator',
            width: 134,
            render(row) {
                return generateButtonGroup(
                    generateUpdateButton(row, update),
                    generateDeleteButton(row, remove),
                )
            }
        },
    ]
}

function generateAdminTag(row) {
    return h(NTag, {
        size: 'small',
        type: row.admin === 1 ? 'success' : '',
    }, {
        default: () => row.admin === 1 ? '管理员' : '普通用户'
    })
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