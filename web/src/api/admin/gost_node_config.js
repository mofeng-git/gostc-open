import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/node/config'

export const apiAdminGostNodeConfigPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodeConfigList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodeConfigCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodeConfigUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiAdminGostNodeConfigDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}