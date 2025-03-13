import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/node'

export const apiNormalGostNodeList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}

export const apiNormalGostNodeCleanPort = (data) => {
    return request.request({
        url: `${baseUrl}/cleanPort`,
        method: 'POST',
        data
    })
}


export const apiNormalGostNodePage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostNodeCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostNodeUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiNormalGostNodeDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}
