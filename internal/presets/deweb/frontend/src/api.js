const API = '/api'

export function getToken() {
  return localStorage.getItem('token')
}

export function getUser() {
  const raw = localStorage.getItem('user')
  return raw ? JSON.parse(raw) : null
}

export function setSession(token, user) {
  localStorage.setItem('token', token)
  localStorage.setItem('user', JSON.stringify(user))
}

export function clearSession() {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

export async function api(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`
  const res = await fetch(API + path, { ...options, headers })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Ошибка запроса')
  return data
}

export function requireAuth(redirect = '/index.html') {
  if (!getToken()) {
    window.location.href = redirect
    return false
  }
  return true
}

export function redirectByRole() {
  const user = getUser()
  if (!user) return (window.location.href = '/index.html')
  if (user.is_admin) window.location.href = '/admin.html'
  else window.location.href = '/cabinet.html'
}
