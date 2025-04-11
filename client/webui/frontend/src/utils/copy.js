export function copyToClipboard(value) {
    if (!navigator){
        return
    }
    return navigator.clipboard && navigator.clipboard.writeText(value)
}