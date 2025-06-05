import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client/tunnel'

export const apiNormalGostClientTunnelList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiNormalGostClientTunnelPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientTunnelCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientTunnelRenew = (data) => {
    return request.request({
        url: `${baseUrl}/renew`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientTunnelMigrate = (data) => {
    return request.request({
        url: `${baseUrl}/migrate`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientTunnelEnable = (data) => {
    return request.request({
        url: `${baseUrl}/enable`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientTunnelUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiNormalGostClientTunnelDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}