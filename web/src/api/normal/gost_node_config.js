import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/node/config'

export const apiNormalGostNodeConfigList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}
