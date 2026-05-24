import { initRegisterForm } from '../auth.js'

initRegisterForm('register-form', (user) => {
  window.location.href = user.is_admin ? '/admin.html' : '/create.html'
})
