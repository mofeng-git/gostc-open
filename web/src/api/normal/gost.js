import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost'

export const apiNormalGostInfo = (data) => {
    return request.request({
        url: `${baseUrl}/info`,
        method: 'POST',
        data
    })
}
