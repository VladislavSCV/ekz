import { initLoginForm } from '../auth.js'
import { getToken, getUser } from '../api.js'

if (getToken()) {
  const u = getUser()
  if (u?.is_admin) window.location.href = '/admin.html'
  else window.location.href = '/cabinet.html'
} else {
  initLoginForm('login-form', 'login-alert')
}
