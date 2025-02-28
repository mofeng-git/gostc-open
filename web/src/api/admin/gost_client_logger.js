import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client/logger'

export const apiAdminGostClientLoggerPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}
