import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client'

export const apiNormalGostClientPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}