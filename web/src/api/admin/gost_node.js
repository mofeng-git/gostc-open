import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/node'

export const apiAdminGostNodeList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodeCleanPort = (data) => {
    return request.request({
        url: `${baseUrl}/cleanPort`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodePage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodeCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostNodeUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiAdminGostNodeDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}


export const apiAdminGostNodeQuery = (data) => {
    return request.request({
        url: `${baseUrl}/query`,
        method: 'POST',
        data
    })
}
