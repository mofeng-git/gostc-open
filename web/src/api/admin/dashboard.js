import request from "../../request/index.js";

const baseUrl = '/v1/admin/dashboard'

export const apiAdminDashboardUserObs = () => {
    return request.request({
        url: `${baseUrl}/userObs`, method: 'POST',
    })
}

export const apiAdminDashboardNodeObs = () => {
    return request.request({
        url: `${baseUrl}/nodeObs`, method: 'POST',
    })
}

export const apiAdminDashboardUserObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/userObsDate?date=` + date, method: 'POST',
    })
}

export const apiAdminDashboardNodeObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/nodeObsDate?date=` + date, method: 'POST',
    })
}

export const apiAdminDashboardCount = () => {
    return request.request({
        url: `${baseUrl}/count`, method: 'POST',
    })
}

export const apiAdminDashboardClientObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientObsDate?date=` + date, method: 'POST',
    })
}

export const apiAdminDashboardClientHostObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientHostObsDate?date=` + date, method: 'POST',
    })
}
export const apiAdminDashboardClientForwardObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientForwardObsDate?date=` + date, method: 'POST',
    })
}
export const apiAdminDashboardClientTunnelObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientTunnelObsDate?date=` + date, method: 'POST',
    })
}