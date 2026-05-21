checkServer()

let resume = null

const audio = document.getElementById('bg-music')
const btnMute = document.getElementById('btn-mute')

btnMute.addEventListener('click', () => {
    audio.muted = !audio.muted
    btnMute.textContent = audio.muted ? '🔇' : '🔊'
})

const btn = document.getElementById('btn-generate')
btn.addEventListener('click', generate)


const btnPreview = document.getElementById('btn-preview')
btnPreview.addEventListener('click', preview)

function preview() {
    const url = URL.createObjectURL(resume)
    chrome.tabs.create({ url })
}

const btnDownload = document.getElementById('btn-download')
btnDownload.addEventListener('click', download)

function download() {
    const filename = document.getElementById('filename-input').value || 'resume.pdf'
    const url = URL.createObjectURL(resume)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)

}

function log(msg) {
    const box = document.getElementById('log-box')
    const line = document.createElement('p')
    line.textContent = '> ' + msg
    box.prepend(line)
}

const btnUpload = document.getElementById('btn-upload')
btnUpload.addEventListener('click', upload)
function upload() {
    const filename = document.getElementById('filename-input').value || 'resume.pdf'

    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
        resume.arrayBuffer().then(buffer => {
            chrome.tabs.sendMessage(tabs[0].id, {
                type: 'UPLOAD_RESUME',
                data: Array.from(new Uint8Array(buffer)),
                filename: filename
            })
        })
    })
}

document.getElementById('btn-retry').addEventListener('click', () => {
    document.getElementById('state-error').style.display = 'none'
    document.getElementById('state-idle').style.display = 'block'
    checkServer()
})


const btnBack = document.getElementById('btn-back')
btnBack.addEventListener('click', reset)

function reset() {
    resume = null

    document.getElementById('state-ready').style.display = 'none'
    document.getElementById('state-idle').style.display = 'block'
}


function checkServer() {
    const dot = document.getElementById('server-dot')
    const status = document.getElementById('server-status')

    dot.classList.add('red')    // red

    fetch('http://localhost:8080/health')
    .then(() => {
        dot.classList.remove('red') // back to cyan
        status.textContent = 'connection established'
    })

    .catch(() => {
        dot.classList.add('red')
        status.textContent = 'SERVER UNREACHABLE'
    })
}

function generate() {
    const loadState = document.getElementById('state-loading')
    const idleState = document.getElementById('state-idle')
    const readyState = document.getElementById('state-ready')
    document.getElementById('log-box').innerHTML = ''   // clear last run
    loadState.style.display = 'block'
    idleState.style.display = 'none'
    
    audio.play().catch(err => {
        log('music blocked: ' + err.message)
    })

    log('scraping job description...')
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
        chrome.scripting.executeScript({
            target: { tabId: tabs[0].id },
            func: () => document.body.innerText
        }, (results) => {
            if (chrome.runtime.lastError || !results || !results[0]) {
                audio.pause()
                audio.currentTime = 0
                log('ERROR: cannot read this page')
                loadState.style.display = 'none'
                document.getElementById('error-msg').textContent = 'cannot script this page type'
                document.getElementById('state-error').style.display = 'block'
                return
            }
            const jd = results[0].result
            log('jd captured. feeding to trimmer...')
            fetch('http://localhost:8080/api/resume/rewrite', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jd: jd })
            })
            .then(res => {
                log('models cooking... do not open the oven')
                return res.blob()
            })
            .then(blob => {
                audio.pause()
                audio.currentTime = 0
                log('pdf hot off the press 🔥')
                resume = blob
                loadState.style.display = 'none'
                readyState.style.display = 'block'
            })
            .catch(err => {
                audio.pause()
                audio.currentTime = 0
                log('ERROR: ' + err.message)
                loadState.style.display = 'none'
                document.getElementById('error-msg').textContent = err.message
                document.getElementById('state-error').style.display = 'block'
            })
        })
    })
}
