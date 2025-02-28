import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/node/logger'

export const apiAdminGostNodeLoggerPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}
