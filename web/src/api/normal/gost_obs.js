import request from "../../request/index.js";

const baseUrl = '/v1/normal/gost/obs'

export const apiNormalGostObsTunnelMonth = (data) => {
    return request.request({
        url: `${baseUrl}/tunnel/month`,
        method: 'POST',
        data
    })
}

export const apiNormalGostObsClientMonth = (data) => {
    return request.request({
        url: `${baseUrl}/client/month`,
        method: 'POST',
        data
    })
}

export const apiNormalGostObsNodeMonth = (data) => {
    return request.request({
        url: `${baseUrl}/node/month`,
        method: 'POST',
        data
    })
}

export const apiNormalGostObsUserMonth = (data) => {
    return request.request({
        url: `${baseUrl}/user/month`,
        method: 'POST',
        data
    })
}