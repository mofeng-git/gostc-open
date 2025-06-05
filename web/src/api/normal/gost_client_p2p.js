import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client/p2p'

export const apiNormalGostClientP2PList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiNormalGostClientP2PPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientP2PCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientP2PRenew = (data) => {
    return request.request({
        url: `${baseUrl}/renew`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientP2PMigrate = (data) => {
    return request.request({
        url: `${baseUrl}/migrate`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientP2PEnable = (data) => {
    return request.request({
        url: `${baseUrl}/enable`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientP2PUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiNormalGostClientP2PDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}