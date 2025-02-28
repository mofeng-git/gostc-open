import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/node/rule'

export const apiAdminGostNodeRuleList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}
