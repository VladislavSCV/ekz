package cli

import "fmt"

func printHelp() {
	fmt.Print(`ekz — генератор проектов ДЭ 09.02.07 (Go + SQLite + Vite)

УСТАНОВКА
  go install github.com/VladislavSCV/ekz@v0.1.0
  go install github.com/VladislavSCV/ekz@main          (свежая)
  Не используйте @latest без тега — может быть v0.0.1 (старая).

КОМАНДЫ
  ekz                         мастер: сценарий, столбцы БД, страницы (↑↓ Enter)
  ekz help                    эта справка
  ekz presets                 шаблоны: conferences, food-delivery

  ekz schema init [файл]      пример project.yaml для правки (по умолч. project.yaml)
  ekz schema export <id> [файл]   выгрузить шаблон в YAML (без генерации кода)
  ekz schema help             справка по project.yaml

  ekz -name <папка> -config project.yaml     генерация из YAML
  ekz -name <папка> -preset conferences      быстро: Конференции.РФ
  ekz -name <папка> -preset food-delivery    быстро: доставка еды

  ekz themes                  (устар.) темы оформления
  ekz theme init [файл]       (устар.) theme.yaml

project.yaml — ОПИСАНИЕ БИЛЕТА (не код)
  portal          название сайта
  main.fields     столбцы БД и поля формы (string, int, date, enum)
  main.statuses   статусы для админки
  pages           login, register, create_form, cabinet, admin, slider
  reviews         отзывы true/false
  admin           логин/пароль (на экзамене часто Admin26 / Demo20)

КАК ПОЛУЧИТЬ project.yaml
  1) ekz schema export conferences     шаблон билета
  2) ekz schema init                   пустой пример
  3) ekz → «Только project.yaml»       мастер без кода
  4) после ekz -name X -preset …       файл в X/project.yaml

КАК СДЕЛАТЬ ЛЮБОЙ БИЛЕТ
  A) ekz schema export conferences → правка YAML → ekz -name proj -config project.yaml
  B) ekz  → «Свой билет с нуля» → столбцы и страницы стрелками
  C) ekz -name proj -preset conferences   без YAML

ЧТО ГЕНЕРИРУЕТСЯ (папка proj/)
  backend/     Go, Fiber, GORM, SQLite (users + заявки + reviews)
  frontend/    HTML, Vite, Bootstrap, 390px, слайдер 3 сек
  project.yaml копия схемы

ЗАПУСК ПРОЕКТА
  cd proj/backend && go mod tidy && go run .
  cd proj/frontend && npm install && npm run dev

ЭКЗАМЕН «Конференции.РФ» (вариант 2)
  ekz -name my-de -preset conferences
  ER-диаграмму и картинки из архива М2 — вручную; 3+ git commit по модулям.

Подробнее: README.md и EXAM.md в репозитории
`)
}
