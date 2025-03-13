import request from "../../request/index.js";

const baseUrl = '/v1/admin/dashboard'

export const apiAdminDashboardUserObs = () => {
    return request.request({
        url: `${baseUrl}/userObs`,
        method: 'POST',
    })
}

export const apiAdminDashboardNodeObs = () => {
    return request.request({
        url: `${baseUrl}/nodeObs`,
        method: 'POST',
    })
}

export const apiAdminDashboardCount = () => {
    return request.request({
        url: `${baseUrl}/count`,
        method: 'POST',
    })
}
