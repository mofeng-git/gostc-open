import request from "../request/index.js";

const baseUrl = '/client'

export const apiClientCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`, method: 'POST', data
    })
}

export const apiClientDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`, method: 'POST', data
    })
}

export const apiClientUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`, method: 'POST', data
    })
}

export const apiClientList = (data) => {
    return request.request({
        url: `${baseUrl}/list`, method: 'POST', data
    })
}

export const apiClientStatus = (data) => {
    return request.request({
        url: `${baseUrl}/status`, method: 'POST', data
    })
}
