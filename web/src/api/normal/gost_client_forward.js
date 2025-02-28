import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/client/forward'

export const apiNormalGostClientForwardList = () => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
    })
}

export const apiNormalGostClientForwardPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientForwardCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientForwardRenew = (data) => {
    return request.request({
        url: `${baseUrl}/renew`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientForwardEnable = (data) => {
    return request.request({
        url: `${baseUrl}/enable`,
        method: 'POST',
        data
    })
}

export const apiNormalGostClientForwardUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiNormalGostClientForwardMatcher = (data) => {
    return request.request({
        url: `${baseUrl}/matcher`,
        method: 'POST',
        data
    })
}


export const apiNormalGostClientForwardAdmission = (data) => {
    return request.request({
        url: `${baseUrl}/admission`,
        method: 'POST',
        data
    })
}


export const apiNormalGostClientForwardDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}