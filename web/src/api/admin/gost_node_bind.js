import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/node/bind'

export const apiNormalGostNodeBindUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}

