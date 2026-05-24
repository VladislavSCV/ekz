export function initSlider(root) {
  const slides = [...root.querySelectorAll('.slide')]
  const prev = root.querySelector('[data-slider-prev]')
  const next = root.querySelector('[data-slider-next]')
  let index = 0
  let timer

  function show(i) {
    index = (i + slides.length) % slides.length
    slides.forEach((s, n) => s.classList.toggle('active', n === index))
  }

  function nextSlide() {
    show(index + 1)
  }

  function start() {
    clearInterval(timer)
    timer = setInterval(nextSlide, 3000)
  }

  prev?.addEventListener('click', () => {
    show(index - 1)
    start()
  })
  next?.addEventListener('click', () => {
    show(index + 1)
    start()
  })

  show(0)
  start()
}
