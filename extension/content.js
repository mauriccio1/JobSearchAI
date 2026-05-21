chrome.runtime.onMessage.addListener((message) => {
    if (message.type == 'UPLOAD_RESUME') {
        const bytes = new Uint8Array(message.data)
        const blob = new Blob([bytes], { type: 'application/pdf' })
        const file = new File([blob], message.filename, { type: 'application/pdf' })

        const input = document.querySelector('input[type="file"]')
        if (!input) return

        const dt = new DataTransfer()
        dt.items.add(file)
        input.files = dt.files

        input.dispatchEvent(new Event('change', { bubbles: true }))
        input.dispatchEvent(new Event('input', { bubbles: true }))
    }
})
