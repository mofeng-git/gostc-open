import request from "../request/index.js";

const baseUrl = '/logger'

export const apiLoggerList = (data) => {
    return request.request({
        url: `${baseUrl}/list`, method: 'POST', data
    })
}
