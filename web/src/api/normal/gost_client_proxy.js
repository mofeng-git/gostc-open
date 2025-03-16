import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client/proxy'

export const apiNormalGostClientProxyList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiNormalGostClientProxyPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientProxyCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientProxyRenew = (data) => {
    return request.request({
        url: `${baseUrl}/renew`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientProxyEnable = (data) => {
    return request.request({
        url: `${baseUrl}/enable`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientProxyUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientProxyDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}