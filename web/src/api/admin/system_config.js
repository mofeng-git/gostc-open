import request from "../../request/index.js";

const baseUrl = '/v1/admin/system/config'

export const apiAdminSystemConfigBase = (data) => {
    return request.request({
        url: `${baseUrl}/base`,
        method: 'POST',
        data
    })
}

export const apiAdminSystemConfigGost = (data) => {
    return request.request({
        url: `${baseUrl}/gost`,
        method: 'POST',
        data
    })
}

export const apiAdminSystemConfigQuery = (data) => {
    return request.request({
        url: `${baseUrl}/query`,
        method: 'POST',
        data
    })
}
