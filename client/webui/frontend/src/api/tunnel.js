import request from "../request/index.js";

const baseUrl = '/tunnel'

export const apiTunnelCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`, method: 'POST', data
    })
}

export const apiTunnelDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`, method: 'POST', data
    })
}

export const apiTunnelUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`, method: 'POST', data
    })
}

export const apiTunnelList = (data) => {
    return request.request({
        url: `${baseUrl}/list`, method: 'POST', data
    })
}

export const apiTunnelStatus = (data) => {
    return request.request({
        url: `${baseUrl}/status`, method: 'POST', data
    })
}
