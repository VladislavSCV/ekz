import { api, setSession } from './api.js'

const loginRe = /^[a-zA-Z0-9]+$/

function setHint(id, msg) {
  const el = document.getElementById(id)
  if (el) el.textContent = msg || ''
}

export function initLoginForm(formId, alertId, onSuccess) {
  const form = document.getElementById(formId)
  const alertBox = document.getElementById(alertId)
  form?.addEventListener('submit', async (e) => {
    e.preventDefault()
    alertBox?.classList.add('d-none')
    setHint('login-hint', '')
    setHint('password-hint', '')
    const login = form.login.value.trim()
    const password = form.password.value
    let ok = true
    if (login.length < 6 || !loginRe.test(login)) {
      setHint('login-hint', 'Логин: минимум 6 символов, латиница и цифры')
      ok = false
    }
    if (password.length < 8) {
      setHint('password-hint', 'Пароль: минимум 8 символов')
      ok = false
    }
    if (!ok) return
    try {
      const data = await api('/login', { method: 'POST', body: JSON.stringify({ login, password }) })
      setSession(data.token, data.user)
      onSuccess(data.user)
    } catch (err) {
      alertBox.textContent = err.message
      alertBox.classList.remove('d-none', 'alert-success')
      alertBox.classList.add('alert-warning')
    }
  })
}

export function initRegisterForm(formId, onSuccess) {
  const form = document.getElementById(formId)
  form?.addEventListener('submit', async (e) => {
    e.preventDefault()
    const fields = ['full_name', 'phone', 'email', 'login', 'password']
    let ok = true
    fields.forEach((f) => setHint(`${f}-hint`, ''))
    const payload = {}
    fields.forEach((f) => {
      payload[f] = form[f].value.trim()
    })
    if (!payload.full_name) {
      setHint('full_name-hint', 'Укажите ФИО')
      ok = false
    }
    if (!payload.phone) {
      setHint('phone-hint', 'Укажите телефон')
      ok = false
    }
    if (!payload.email?.includes('@')) {
      setHint('email-hint', 'Укажите e-mail')
      ok = false
    }
    if (payload.login.length < 6 || !loginRe.test(payload.login)) {
      setHint('login-hint', 'Логин: минимум 6 символов, латиница и цифры')
      ok = false
    }
    if (payload.password.length < 8) {
      setHint('password-hint', 'Пароль: минимум 8 символов')
      ok = false
    }
    if (!ok) return
    try {
      const data = await api('/register', { method: 'POST', body: JSON.stringify(payload) })
      setSession(data.token, data.user)
      onSuccess(data.user)
    } catch (err) {
      setHint('login-hint', err.message)
    }
  })
}
