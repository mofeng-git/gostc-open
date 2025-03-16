import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client/forward'

export const apiAdminGostClientForwardList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiAdminGostClientForwardPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientForwardCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientForwardConfig = (data) => {
    return request.request({
        url: `${baseUrl}/config`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientForwardDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}