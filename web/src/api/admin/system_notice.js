import request from "../../request/index.js";

const baseUrl = '/v1/admin/system/notice'

export const apiAdminSystemNoticePage = (data) => {
    return request.request({
        url: `${baseUrl}/page`,
        method: 'POST',
        data
    })
}

export const apiAdminSystemNoticeList = (data) => {
    return request.request({
        url: `${baseUrl}/list`,
        method: 'POST',
        data
    })
}


export const apiAdminSystemNoticeCreate = (data) => {
    return request.request({
        url: `${baseUrl}/create`,
        method: 'POST',
        data
    })
}


export const apiAdminSystemNoticeUpdate = (data) => {
    return request.request({
        url: `${baseUrl}/update`,
        method: 'POST',
        data
    })
}

export const apiAdminSystemNoticeDelete = (data) => {
    return request.request({
        url: `${baseUrl}/delete`,
        method: 'POST',
        data
    })
}