import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client/p2p'

export const apiAdminGostClientP2PList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiAdminGostClientP2PPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientP2PCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientP2PConfig = (data) => {
    return request.request({
        url: `${baseUrl}/config`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientP2PUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiAdminGostClientP2PDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}