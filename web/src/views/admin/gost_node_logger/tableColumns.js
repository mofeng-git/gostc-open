import {h, ref} from "vue";
import {NButton, NSpace} from "naive-ui";
import moment from "moment";

export function generateTableColumns() {
    return [
        {title: 'ID', key: 'id', ellipsis: {tooltip: true},width: 60},
        {
            title: '时间',
            key: 'createdAt',
            ellipsis: {tooltip: true},
            width:180,
            render(row) {
                return moment(row?.createdAt * 1000).format('yyyy-MM-DD HH:mm:ss')
            }
        },
        {title: '日志级别', key: 'level', ellipsis: {tooltip: true},width: 100},
        {title: '内容', key: 'content', ellipsis: {tooltip: true}},
    ]
}
