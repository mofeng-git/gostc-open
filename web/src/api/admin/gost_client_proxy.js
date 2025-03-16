import request from "../../request/index.js";

const baseUrl = '/v1/admin/gost/client/proxy'

export const apiAdminGostClientProxyList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiAdminGostClientProxyPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientProxyConfig = (data) => {
    return request.request({
        url: `${baseUrl}/config`,
        method: 'POST',
        data
    })
}

export const apiAdminGostClientProxyDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}