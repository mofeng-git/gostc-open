import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client'

export const apiAdminGostClientPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}