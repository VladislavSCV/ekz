import { getToken, getUser } from '../api.js'
import { initLoginForm } from '../auth.js'

function redirect(user) {
  if (user.is_admin) window.location.href = '/admin.html'
  else if (document.querySelector('[data-after-login]'))
    window.location.href = document.querySelector('[data-after-login]').dataset.afterLogin
  else window.location.href = '/cabinet.html'
}

if (getToken()) redirect(getUser())
else initLoginForm('login-form', 'login-alert', redirect)
