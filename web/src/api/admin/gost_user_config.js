import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/user/config'

export const apiAdminGostUserConfigPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostUserConfigList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}

export const apiAdminGostUserConfigNodeList = (data) => {
    return request.request({
        url: `${baseUrl}/node/list`,
        method: 'POST',
        data
    })
}

export const apiAdminGostUserConfigCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostUserConfigUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiAdminGostUserConfigDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}