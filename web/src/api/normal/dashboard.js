import request from "../../request/index.js";

const baseUrl = '/v1/normal/dashboard'

export const apiNormalDashboardCount = () => {
    return request.request({
        url: `${baseUrl}/count`, method: 'POST',
    })
}

export const apiNormalDashboardClientObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientObsDate?date=` + date, method: 'POST',
    })
}

export const apiNormalDashboardClientHostObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientHostObsDate?date=` + date, method: 'POST',
    })
}
export const apiNormalDashboardClientForwardObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientForwardObsDate?date=` + date, method: 'POST',
    })
}
export const apiNormalDashboardClientTunnelObsDate = (date) => {
    return request.request({
        url: `${baseUrl}/clientTunnelObsDate?date=` + date, method: 'POST',
    })
}