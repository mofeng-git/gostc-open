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
