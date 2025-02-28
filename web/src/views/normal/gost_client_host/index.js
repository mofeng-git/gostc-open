export const limiterText = (limiter) => {
    if (limiter > 0) {
        return limiter + 'mbps'
    }
    return '无限制'
}

export const cLimiterText = (limiter) => {
    if (limiter > 0) {
        return limiter
    }
    return '无限制'
}

export const rLimiterText = (limiter) => {
    if (limiter > 0) {
        return limiter
    }
    return '无限制'
}

export const configText = (config) => {
    switch (config.chargingType) {
        case 1:
            return `一次性计费,${config.amount}积分`
        case 2:
            return `循环计费,${config.amount}积分/${config.cycle}天`
        case 3:
            return `免费`
    }
}

export const configExpText = (config) => {
    switch (config.chargingType) {
        case 1:
            return `免续费`
        case 2:
            return config.expAt
        case 3:
            return `免续费`
    }
}