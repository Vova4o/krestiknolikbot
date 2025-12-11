# Telegram Mini App «Крестики-нолики»

Мини-приложение Telegram, которое запускается через официальную Menu Button и предлагает сыграть в крестики-нолики против компьютера. При победе фронтенд отправляет результат на backend, который верифицирует `initData`, извлекает `user.id` и отправляет пользователю сообщение с промокодом.

## Что входит

- Backend на Go (Gin + long polling Telegram Bot API)
- Frontend (HTML/CSS/JS) с авторизацией через `initData`
- Дизайн под целевую аудиторию женщин 25–40 лет

## Быстрый старт

1. Скопируйте `.env.example` в `.env` и заполните значения.
2. Загрузите зависимости:
   ```bash
   go mod tidy
   ```
3. Запустите сервер:
   ```bash
   go run ./cmd/server/main.go
   ```
4. Откройте mini app локально: `http://localhost:8080/web/`.

## Подключение Menu Button в BotFather

- `/setmenubutton` → выбрать бота → тип `web_app`
- URL: `https://<ваш-домен>/web/`

## Docker

### Сборка и запуск Dockerfile
```bash
docker build -t tg-tictactoe-miniapp .
docker run --rm -p 8080:8080 \
  -e TELEGRAM_BOT_TOKEN=... \
  -e PORT=8080 \
  tg-tictactoe-miniapp
```

### docker-compose (локально / Dockploy)
```bash
docker-compose up --build
```
- Укажите `TELEGRAM_BOT_TOKEN` и (при необходимости) `PORT` в переменных окружения (Dockploy хранит их в интерфейсе; локально можно прокинуть через shell или `.env`).
- Порты не проброшены наружу, чтобы избежать конфликтов — Dockploy сам назначит. Для локальной отладки добавьте в `docker-compose.yml` секцию `ports: - "8080:8080"` или используйте `docker run -p 8080:8080 ...`.
- Публичный URL в Menu Button: `https://<домен>/web/` (контейнер отдаёт статику по этому пути).

## Проверка результатов

- Победа: пользователю отправляется «Победа! Промокод выдан: <код>»
- Поражение: «Проигрыш»
- Ничья: «Ничья»

## Продакшн

- Храните `TELEGRAM_BOT_TOKEN` в переменных окружения/secret manager.
- Защитите домен HTTPS (для mini app внутри Telegram обязателен HTTPS).
- При деплое убедитесь, что `PORT` открыт и статика из `web/` доступна по `/web/`.

export TELEGRAM_BOT_TOKEN='8511938965:AAGa_N1IJurukyfui3MttOv_6FkbARao_YI'
export TELEGRAM_WEBHOOK_SECRET='8cbe1719ad6bd1589f775fe09f8976c6582865d5713b2833d0119ab747fbb1ab'
curl -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook" \
  -d "url=https://tictactoe.vova4o.com/api/telegram/webhook" \
  -d "secret_token=$TELEGRAM_WEBHOOK_SECRET" \
  -d "drop_pending_updates=true"
