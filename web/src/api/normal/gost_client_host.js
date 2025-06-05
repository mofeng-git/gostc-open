import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client/host'

export const apiNormalGostClientHostList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiNormalGostClientHostPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostDomain = (data) => {
    return request.request({
        url: `${baseUrl}/domain`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostRenew = (data) => {
    return request.request({
        url: `${baseUrl}/renew`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostMigrate = (data) => {
    return request.request({
        url: `${baseUrl}/migrate`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostEnable = (data) => {
    return request.request({
        url: `${baseUrl}/enable`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiNormalGostClientHostDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientHostAdmission = (data) => {
    return request.request({
        url: `${baseUrl}/admission`,
        method: 'POST',
        data
    })
}