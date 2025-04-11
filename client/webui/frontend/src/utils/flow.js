export function flowFormat(value) {
    if (value < 1024) {
        return value + 'B'
    }
    if (value < 1024 * 1024) {
        return (value / 1024).toFixed(1) + 'KB'
    }
    if (value < 1024 * 1024 * 1024) {
        return (value / 1024 / 1024).toFixed(1) + 'MB'
    }
    if (value < 1024 * 1024 * 1024 * 1024) {
        return (value / 1024 / 1024 / 1024).toFixed(1) + 'GB'
    }
    if (value < 1024 * 1024 * 1024 * 1024 * 1024) {
        return (value / 1024 / 1024 / 1024 / 1024).toFixed(1) + 'TB'
    }
    return '--'
}