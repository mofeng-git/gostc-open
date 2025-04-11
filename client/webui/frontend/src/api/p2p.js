import request from "../request/index.js";

const baseUrl = '/p2p'

export const apiP2PCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`, method: 'POST', data
    })
}

export const apiP2PDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`, method: 'POST', data
    })
}

export const apiP2PUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`, method: 'POST', data
    })
}

export const apiP2PList = (data) => {
    return request.request({
        url: `${baseUrl}/list`, method: 'POST', data
    })
}

export const apiP2PStatus = (data) => {
    return request.request({
        url: `${baseUrl}/status`, method: 'POST', data
    })
}
