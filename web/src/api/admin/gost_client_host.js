import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client/host'

export const apiAdminGostClientHostList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiAdminGostClientHostPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientHostCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientHostConfig = (data) => {
    return request.request({
        url: `${baseUrl}/config`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientHostUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiAdminGostClientHostDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}