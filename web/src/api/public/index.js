import request from "../../request/index.js"

const baseUrl = '/v1/public/system/config'

export const apiSystemConfigQuery = () => {
    return request.request({
        url: `${baseUrl}/query`,
        method: 'POST',
        data:{}
    })
}