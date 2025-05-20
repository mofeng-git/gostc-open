import request from "../../request/index.js"

const baseUrl = '/v1/auth'

// 登录
export const apiAuthLogin = (data) => {
    return request.request({
        url: `${baseUrl}/login`,
        method: 'POST',
        data
    })
}

// 登录
export const apiAuthLoginOtp = (data) => {
    return request.request({
        url: `${baseUrl}/loginOtp`,
        method: 'POST',
        data
    })
}

// 注册
export const apiAuthRegister = (data) => {
    return request.request({
        url: `${baseUrl}/register`,
        method: 'POST',
        data
    })
}

// 修改密码
export const apiAuthPassword = (data) => {
    return request.request({
        url: `${baseUrl}/password`,
        method: 'POST',
        data
    })
}

// 获取验证码
export const apiAuthCaptcha = () => {
    return request.request({
        url: `${baseUrl}/captcha?timestamp=${new Date().getTime()}`,
        method: 'POST'
    })
}

// 用户信息
export const apiAuthUserInfo = () => {
    return request.request({
        url: `${baseUrl}/userInfo`,
        method: 'POST'
    })
}

// 签到
export const apiAuthUserCheckin = () => {
    return request.request({
        url: `${baseUrl}/checkin`,
        method: 'POST'
    })
}

// 续期
export const apiAuthRenew = () => {
    return request.request({
        url: `${baseUrl}/renew`,
        method: 'POST'
    })
}

// 生成二步验证二维码
export const apiAuthGenOtp = () => {
    return request.request({
        url: `${baseUrl}/genOtp`,
        method: 'POST'
    })
}

// 开启二步验证
export const apiAuthOpenOtp = (data) => {
    return request.request({
        url: `${baseUrl}/openOtp`,
        method: 'POST',
        data
    })
}

// 关闭二步验证
export const apiAuthCloseOtp = (data) => {
    return request.request({
        url: `${baseUrl}/closeOtp`,
        method: 'POST',
        data
    })
}