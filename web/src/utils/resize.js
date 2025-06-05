export const resizeDirective = {
    mounted(el, binding) {
        const callback = binding.value
        const observer = new ResizeObserver(entries => {
            entries.forEach(entry => {
                callback(entry.contentRect)
            })
        })
        observer.observe(el)
        el._resizeObserver = observer
    },
    unmounted(el) {
        if (el._resizeObserver) {
            el._resizeObserver.disconnect()
        }
    }
}