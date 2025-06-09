import request from "../../request/index.js";

const baseUrl = '/v1/normal/frp/client/cfg'


export const apiNormalFrpClientCfgPage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiNormalFrpClientCfgCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}

export const apiNormalFrpClientCfgMigrate = (data) => {
    return request.request({
        url: `${baseUrl}/migrate`,
        method: 'POST',
        data
    })
}

export const apiNormalFrpClientCfgEnable = (data) => {
    return request.request({
        url: `${baseUrl}/enable`,
        method: 'POST',
        data
    })
}

export const apiNormalFrpClientCfgUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}


export const apiNormalFrpClientCfgDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}
