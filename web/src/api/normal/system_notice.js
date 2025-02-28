import request from "../../request/index.js";

const baseUrl = '/v1/normal/system/notice'

export const apiNormalSystemNoticeList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}
