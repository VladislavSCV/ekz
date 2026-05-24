import { api, getUser, requireAuth } from '../api.js'

if (!requireAuth()) throw new Error('auth')

const form = document.getElementById('booking-form')
const alertBox = document.getElementById('booking-alert')

function hint(name, msg) {
  const el = document.getElementById(`${name}-hint`)
  if (el) el.textContent = msg || ''
}

form.addEventListener('submit', async (e) => {
  e.preventDefault()
  alertBox.classList.add('d-none')
  ;['room_type', 'start_date', 'payment_method'].forEach((f) => hint(f, ''))
  const payload = {
    room_type: form.room_type.value,
    start_date: form.start_date.value.trim(),
    payment_method: form.payment_method.value,
  }
  let ok = true
  if (!payload.room_type) {
    hint('room_type', 'Выберите тип помещения')
    ok = false
  }
  if (!/^\d{2}\.\d{2}\.\d{4}$/.test(payload.start_date)) {
    hint('start_date', 'Формат: ДД.ММ.ГГГГ')
    ok = false
  }
  if (!payload.payment_method) {
    hint('payment_method', 'Выберите способ оплаты')
    ok = false
  }
  if (!ok) return
  try {
    await api('/bookings', { method: 'POST', body: JSON.stringify(payload) })
    alertBox.textContent = 'Заявка успешно создана'
    alertBox.className = 'alert alert-success alert-inline'
    form.reset()
    setTimeout(() => (window.location.href = '/cabinet.html'), 800)
  } catch (err) {
    alertBox.textContent = err.message
    alertBox.className = 'alert alert-warning alert-inline'
    alertBox.classList.remove('d-none')
  }
})

const user = getUser()
if (user?.is_admin) window.location.href = '/admin.html'
