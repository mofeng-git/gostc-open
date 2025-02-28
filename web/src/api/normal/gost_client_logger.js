import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client/logger'

export const apiNormalGostClientLoggerPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}
