import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client/tunnel'

export const apiAdminGostClientTunnelList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiAdminGostClientTunnelPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientTunnelCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientTunnelConfig = (data) => {
    return request.request({
        url: `${baseUrl}/config`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientTunnelUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiAdminGostClientTunnelDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}