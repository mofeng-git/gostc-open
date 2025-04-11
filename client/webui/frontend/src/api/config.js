import request from "../request/index.js";

const baseUrl = '/config'

export const apiConfigQuery = (data) => {
    return request.request({
        url: `${baseUrl}/query`, method: 'POST', data
    })
}
