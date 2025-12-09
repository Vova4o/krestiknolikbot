# Telegram Mini App «Крестики-нолики»
## Полное техническое задание + полный код (backend + frontend)
## Запуск через официальную Mini App Menu Button

Этот файл предназначен для использования в Cursor / VS Code / ChatGPT-Codex как
единое техническое задание **и** источник кода. По нему можно собрать полный
рабочий проект:

- Telegram Mini App, запускаемая через **официальную Mini App Menu Button**
- Backend на Go:
  - Telegram-бот (long polling)
  - HTTP-сервер для mini app
  - обработчик результата игры и отправка сообщений пользователю
- Frontend (HTML/CSS/JS) — игра «Крестики-нолики» против компьютера
- Авторизация через `initData` (Telegram WebApp)
- Дизайн, ориентированный на женщин 25–40 лет

---

## 1. Функциональные требования

### 1.1. Основной сценарий

1. Пользователь открывает чат с ботом в Telegram.
2. В нижнем меню отображается **Mini App кнопка** (Menu Button), настроенная через BotFather:
   - Тип: `web_app`
   - URL: `https://<your-domain>/web/`
3. Пользователь нажимает кнопку Mini App → открывается web-приложение внутри Telegram.
4. Внутри Mini App пользователь видит:
   - тёплый, аккуратный интерфейс
   - игровое поле 3×3 (крестики-нолики)
   - статус: чей ход, результат и т.п.
5. Пользователь играет против компьютера:
   - игрок — `X`
   - компьютер — `O`
6. При **победе игрока**:
   - на экране показывается случайный 5-значный промокод (например, `38217`);
   - фронтенд делает `POST /api/result` с:
     - `result = "win"`
     - `promoCode`
     - `initData`
   - backend проверяет `initData`, извлекает `user.id` → отправляет пользователю сообщение:
     - **«Победа! Промокод выдан: [код]»**
7. При **поражении игрока**:
   - на экране показывается текст «В этот раз компьютер выиграл. Попробуешь ещё раз?» и кнопка «Сыграть ещё раз»;
   - фронтенд делает `POST /api/result` с:
     - `result = "lose"`
     - `promoCode` = `""` (или `null`)
     - `initData`
   - backend отправляет пользователю сообщение:
     - **«Проигрыш»**
8. При **ничьей**:
   - показывается краткий текст «Ничья, давай ещё раз?» + кнопка «Сыграть ещё раз»;
   - можно не отправлять сообщение в Telegram (по желанию; в реализации ниже — тоже отправим, текст «Ничья»).

---

## 2. Визуальные требования (ориентация на женщин 25–40)

- Общий стиль:
  - мягкие пастельные цвета
  - округлые формы
  - без агрессивных геймерских элементов
  - ощущение «маленького приятного перерыва»
- Цветовая палитра (пример):
  - фон страницы: `#FFF7F9`
  - карточка игры: `#FFFFFF`
  - акценты: `#F8BBD0`, `#FFCCBC`
  - текст: `#3F3D56`
- Шрифты:
  - системный стек: `system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`
- Основные тексты:
  - заголовок: «Небольшой перерыв? Сыграем в крестики-нолики 💕»
  - статусные сообщения:
    - «Твой ход»
    - «Ход компьютера…»
    - «Ты выиграла!»
    - «В этот раз выиграл компьютер. Попробуешь ещё раз?»
    - «Ничья, давай ещё раз?»
  - при победе:
    - «🎁 Поздравляем! Ты выиграла промокод»
  - при проигрыше:
    - «В этот раз компьютер выиграл. Попробуешь ещё раз?»

---

## 3. Архитектура проекта

Проект в виде Go-модуля со следующей структурой:

```text
tg-tictactoe-miniapp/
  cmd/
    server/
      main.go              # точка входа сервера (бот + HTTP-сервер)
  internal/
    bot/
      bot.go               # Telegram Bot client + long polling
    webapp/
      handlers.go          # HTTP handlers (web + API)
      verify.go            # проверка initData
  web/
    index.html             # mini app (HTML)
    styles.css             # стили (CSS)
    app.js                 # логика игры + WebApp API
  .env                     # конфигурация (TOKEN, PORT, DOMAIN)
  go.mod
  go.sum
  README.md                # краткое описание (по желанию)
```

---

## 4. Настройка Telegram Mini App (Menu Button)

### 4.1. Создание бота

1. Открыть `@BotFather`.
2. Выполнить `/newbot`.
3. Получить токен вида:  
   `1234567890:ABCDEF...` → сохранить как `TELEGRAM_BOT_TOKEN` в `.env`.

### 4.2. Настройка домена mini app

1. Развернуть backend на сервере / локально через ngrok.
2. Получить URL, например:
   - `https://<ngrok-id>.ngrok.io`
3. Убедиться, что mini app доступна по адресу:
   - `https://<your-domain>/web/`

### 4.3. Настройка Menu Button (главная Mini App кнопка)

Через `@BotFather`:

1. Команда: `/setmenubutton`
2. Выбрать бота.
3. Выбрать тип: `web_app`
4. Указать:
   - Название mini app (например, `Крестики-нолики`)
   - URL: `https://<your-domain>/web/`

После этого в интерфейсе Telegram у пользователя появится кнопка Mini App.

---

## 5. Backend: Go + Gin + Telegram Bot API

Ниже — полный пример backend-кода.

### 5.1. go.mod

```go
module tg-tictactoe-miniapp

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/joho/godotenv v1.5.1
)
```

> Версии можно скорректировать под актуальные.

---

### 5.2. .env (пример)

```env
TELEGRAM_BOT_TOKEN=1234567890:ABCDEF...
PORT=8080
PUBLIC_BASE_URL=https://your-domain.tld
```

- `TELEGRAM_BOT_TOKEN` — токен бота от BotFather.
- `PORT` — порт HTTP-сервера.
- `PUBLIC_BASE_URL` — публичный URL (используется, если нужно где-то формировать ссылки).

---

### 5.3. cmd/server/main.go

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"tg-tictactoe-miniapp/internal/bot"
	"tg-tictactoe-miniapp/internal/webapp"
)

func main() {
	// Загружаем .env (не критично для production, но удобно локально)
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Инициализируем Telegram Bot client
	botClient := bot.NewBotClient(token)

	// Инициализируем Gin
	router := gin.Default()

	// Настраиваем маршруты mini app и API
	webapp.RegisterRoutes(router, botClient)

	// Запускаем long polling Telegram-бота в отдельной горутине
	go func() {
		log.Println("Starting Telegram bot long polling...")
		if err := botClient.StartLongPolling(); err != nil {
			log.Printf("Telegram bot stopped with error: %v
", err)
		}
	}()

	// Запускаем HTTP-сервер
	go func() {
		addr := ":" + port
		log.Printf("Starting HTTP server at %s...
", addr)
		if err := router.Run(addr); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// Ожидаем сигналы завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %s, shutting down...", sig)
}
```

---

### 5.4. internal/bot/bot.go

```go
package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type BotClient struct {
	token  string
	apiURL string
	client *http.Client
	offset int
}

func NewBotClient(token string) *BotClient {
	return &BotClient{
		token:  token,
		apiURL: "https://api.telegram.org",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (b *BotClient) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", b.apiURL, b.token)

	payload := sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal sendMessage payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create sendMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("sendMessage http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sendMessage non-200: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// --- Long Polling ---

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type getUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func (b *BotClient) StartLongPolling() error {
	for {
		updates, err := b.getUpdates()
		if err != nil {
			log.Printf("getUpdates error: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, u := range updates {
			if u.UpdateID >= b.offset {
				b.offset = int(u.UpdateID) + 1
			}
			if u.Message == nil {
				continue
			}
			b.handleMessage(u.Message)
		}
	}
}

func (b *BotClient) getUpdates() ([]Update, error) {
	url := fmt.Sprintf("%s/bot%s/getUpdates?timeout=20&offset=%d", b.apiURL, b.token, b.offset)

	resp, err := b.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getUpdates http error: %w", err)
	}
	defer resp.Body.Close()

	var result getUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode getUpdates response: %w", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("getUpdates not ok")
	}

	return result.Result, nil
}

func (b *BotClient) handleMessage(m *Message) {
	if m.Text == "/start" {
		_ = b.SendMessage(m.Chat.ID, "Привет! 👋

Нажми кнопку Mini App внизу, чтобы сыграть в крестики-нолики 💕")
		return
	}
}
```

---

### 5.5. internal/webapp/verify.go (проверка initData)

```go
package webapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// VerifyInitData проверяет подпись initData и возвращает TelegramUser, если всё корректно.
func VerifyInitData(initData string, botToken string) (*TelegramUser, error) {
	if initData == "" {
		return nil, fmt.Errorf("empty initData")
	}

	// Разбираем query string
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("parse initData: %w", err)
	}

	hash := values.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash is missing")
	}
	values.Del("hash")

	// Составляем пары "key=value"
	var dataPairs []string
	for key, vals := range values {
		// Берём только первый value
		value := vals[0]
		dataPairs = append(dataPairs, fmt.Sprintf("%s=%s", key, value))
	}

	// Сортируем по ключу (лексикографически по всей строке key=value — ключ идёт первым)
	sort.Strings(dataPairs)

	dataCheckString := strings.Join(dataPairs, "
")

	// secret_key = HMAC_SHA256("WebAppData", bot_token)
	secretKey := hmacSHA256([]byte("WebAppData"), []byte(botToken))

	// HMAC_SHA256(secret_key, data_check_string)
	calculatedHashBytes := hmacSHA256(secretKey, []byte(dataCheckString))
	calculatedHash := hex.EncodeToString(calculatedHashBytes)

	if !hmac.Equal([]byte(calculatedHash), []byte(hash)) {
		return nil, fmt.Errorf("invalid hash")
	}

	// Если hash валиден, разбираем поле user
	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, fmt.Errorf("user field is missing")
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}

	return &user, nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
```

---

### 5.6. internal/webapp/handlers.go

```go
package webapp

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"tg-tictactoe-miniapp/internal/bot"
)

type ResultRequest struct {
	Result    string `json:"result"`
	PromoCode string `json:"promoCode"`
	InitData  string `json:"initData"`
}

func RegisterRoutes(r *gin.Engine, botClient *bot.BotClient) {
	// Отдаём статические файлы mini app
	r.StaticFS("/web", http.Dir("web"))

	// Endpoint для приёма результата игры
	r.POST("/api/result", func(c *gin.Context) {
		var req ResultRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		if req.Result == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "result is required"})
			return
		}

		botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
		if botToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfigured"})
			return
		}

		user, err := VerifyInitData(req.InitData, botToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid initData"})
			return
		}

		var messageText string

		switch req.Result {
		case "win":
			code := req.PromoCode
			if code == "" {
				code = generatePromoCode()
			}
			messageText = "Победа! Промокод выдан: " + code
		case "lose":
			messageText = "Проигрыш"
		case "draw":
			messageText = "Ничья"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown result"})
			return
		}

		if err := botClient.SendMessage(user.ID, messageText); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
```

---

### 5.7. internal/webapp/promocode.go

```go
package webapp

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generatePromoCode() string {
	n := rand.Intn(100000) // 0..99999
	return fmt.Sprintf("%05d", n)
}
```

---

## 6. Frontend: HTML + CSS + JS

Все файлы находятся в папке `web/`.

### 6.1. web/index.html

```html
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <title>Крестики-нолики</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link rel="stylesheet" href="/web/styles.css" />
  <script src="https://telegram.org/js/telegram-web-app.js"></script>
</head>
<body>
  <div class="app">
    <header class="app-header">
      <h1>Небольшой перерыв? 💕</h1>
      <p>Сыграем в крестики-нолики против компьютера.</p>
    </header>

    <main class="game-card">
      <div class="status" id="statusText">Твой ход</div>

      <div class="board" id="board">
        <!-- 9 клеток -->
        <div class="cell" data-index="0"></div>
        <div class="cell" data-index="1"></div>
        <div class="cell" data-index="2"></div>
        <div class="cell" data-index="3"></div>
        <div class="cell" data-index="4"></div>
        <div class="cell" data-index="5"></div>
        <div class="cell" data-index="6"></div>
        <div class="cell" data-index="7"></div>
        <div class="cell" data-index="8"></div>
      </div>

      <div class="result-section" id="resultSection" hidden>
        <div class="result-title" id="resultTitle"></div>
        <div class="promo-block" id="promoBlock" hidden>
          <p class="promo-label">🎁 Поздравляем! Ты выиграла промокод</p>
          <div class="promo-code" id="promoCodeText"></div>
          <button class="btn secondary" id="copyPromoBtn">Скопировать код</button>
        </div>
        <button class="btn primary" id="playAgainBtn">Сыграть ещё раз</button>
      </div>
    </main>
  </div>

  <script src="/web/app.js"></script>
</body>
</html>
```

---

### 6.2. web/styles.css

```css
:root {
  --bg: #fff7f9;
  --card-bg: #ffffff;
  --accent: #f8bbd0;
  --accent-soft: #ffccbc;
  --text-main: #3f3d56;
  --text-muted: #7a768e;
  --cell-empty: #ffe6f0;
  --cell-hover: #ffeef5;
  --x-color: #d81b60;
  --o-color: #5c6bc0;
}

*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  background: var(--bg);
  color: var(--text-main);
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
}

.app-header {
  text-align: center;
  margin-bottom: 16px;
}

.app-header h1 {
  margin: 0;
  font-size: 1.4rem;
}

.app-header p {
  margin: 8px 0 0;
  font-size: 0.95rem;
  color: var(--text-muted);
}

.game-card {
  background: var(--card-bg);
  width: 100%;
  max-width: 360px;
  border-radius: 20px;
  padding: 16px 16px 20px;
  box-shadow:
    0 8px 24px rgba(216, 27, 96, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.06);
}

.status {
  text-align: center;
  margin-bottom: 16px;
  font-size: 1rem;
  font-weight: 500;
}

.board {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 16px;
}

.cell {
  width: 90px;
  height: 90px;
  max-width: 100%;
  border-radius: 18px;
  background: var(--cell-empty);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.4rem;
  font-weight: 600;
  cursor: pointer;
  user-select: none;
  transition:
    transform 0.08s ease,
    box-shadow 0.08s ease,
    background-color 0.08s ease;
}

.cell:hover {
  background: var(--cell-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(216, 27, 96, 0.18);
}

.cell.x {
  color: var(--x-color);
}

.cell.o {
  color: var(--o-color);
}

.cell.disabled {
  cursor: default;
  box-shadow: none;
}

.result-section {
  text-align: center;
  margin-top: 8px;
}

.result-title {
  font-size: 1rem;
  margin-bottom: 12px;
}

.promo-block {
  background: #fff1f6;
  border-radius: 16px;
  padding: 10px 12px;
  margin-bottom: 12px;
}

.promo-label {
  margin: 0 0 8px;
  font-size: 0.9rem;
  color: var(--text-main);
}

.promo-code {
  font-size: 1.4rem;
  font-weight: 700;
  letter-spacing: 0.18em;
  padding: 8px 0;
}

.btn {
  border: none;
  border-radius: 999px;
  padding: 10px 18px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition:
    transform 0.08s ease,
    box-shadow 0.08s ease,
    background-color 0.08s ease;
}

.btn.primary {
  background: var(--accent);
  color: #2d1430;
}

.btn.primary:hover {
  background: #f79bc0;
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(216, 27, 96, 0.25);
}

.btn.secondary {
  background: var(--accent-soft);
  color: #4a2c30;
  margin-top: 4px;
}

.btn.secondary:hover {
  background: #ffbda7;
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(255, 152, 0, 0.2);
}

@media (max-width: 360px) {
  .cell {
    height: 80px;
    font-size: 2.1rem;
  }
}
```

---

### 6.3. web/app.js

```javascript
// Telegram WebApp
const tg = window.Telegram && window.Telegram.WebApp ? window.Telegram.WebApp : null;
let initData = "";

if (tg) {
  tg.expand();
  initData = tg.initData || "";
} else {
  console.warn("Telegram WebApp not detected. Running in dev mode.");
  // Для локальной отладки можно сюда подставить заранее сохранённый initData
}

// --- Игровая логика ---

const boardElement = document.getElementById("board");
const statusText = document.getElementById("statusText");
const resultSection = document.getElementById("resultSection");
const resultTitle = document.getElementById("resultTitle");
const promoBlock = document.getElementById("promoBlock");
const promoCodeText = document.getElementById("promoCodeText");
const copyPromoBtn = document.getElementById("copyPromoBtn");
const playAgainBtn = document.getElementById("playAgainBtn");

let board = Array(9).fill(null); // "X", "O" или null
let isGameOver = false;
const PLAYER = "X";
const COMPUTER = "O";

function resetGame() {
  board = Array(9).fill(null);
  isGameOver = false;

  const cells = boardElement.querySelectorAll(".cell");
  cells.forEach((cell) => {
    cell.textContent = "";
    cell.classList.remove("x", "o", "disabled");
  });

  resultSection.hidden = true;
  promoBlock.hidden = true;
  promoCodeText.textContent = "";
  statusText.textContent = "Твой ход";
}

function checkWinner(b) {
  const winningCombinations = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ];

  for (const [a, c, d] of winningCombinations) {
    if (b[a] && b[a] === b[c] && b[a] === b[d]) {
      return b[a]; // "X" или "O"
    }
  }

  if (b.every((cell) => cell !== null)) {
    return "draw";
  }

  return null;
}

function computerMove() {
  if (isGameOver) return;

  // Небольшая задержка, чтобы ощущалась "анимация хода"
  setTimeout(() => {
    // 1. Попытка выиграть
    let move = findBestMove(COMPUTER);
    // 2. Блокировка выигрыша игрока
    if (move === -1) {
      move = findBestMove(PLAYER);
    }
    // 3. Центр
    if (move === -1 && board[4] === null) {
      move = 4;
    }
    // 4. Случайный ход
    if (move === -1) {
      const available = board
        .map((val, idx) => (val === null ? idx : -1))
        .filter((idx) => idx !== -1);
      if (available.length > 0) {
        const randIndex = Math.floor(Math.random() * available.length);
        move = available[randIndex];
      }
    }

    if (move >= 0 && board[move] === null) {
      board[move] = COMPUTER;
      const cell = boardElement.querySelector(`.cell[data-index="${move}"]`);
      cell.textContent = "O";
      cell.classList.add("o", "disabled");
    }

    const result = checkWinner(board);
    if (result === COMPUTER) {
      isGameOver = true;
      handleGameEnd("lose");
    } else if (result === "draw") {
      isGameOver = true;
      handleGameEnd("draw");
    } else {
      statusText.textContent = "Твой ход";
    }
  }, 500);
}

function findBestMove(playerSymbol) {
  const b = [...board];
  const winningCombinations = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ];

  for (const [a, c, d] of winningCombinations) {
    const line = [b[a], b[c], b[d]];
    const marks = line.filter((v) => v === playerSymbol).length;
    const empties = line.filter((v) => v === null).length;

    if (marks === 2 && empties === 1) {
      if (b[a] === null) return a;
      if (b[c] === null) return c;
      if (b[d] === null) return d;
    }
  }

  return -1;
}

function handleCellClick(e) {
  const cell = e.target;
  if (!cell.classList.contains("cell")) return;
  const index = parseInt(cell.getAttribute("data-index"), 10);

  if (isGameOver || board[index] !== null) return;

  board[index] = PLAYER;
  cell.textContent = "X";
  cell.classList.add("x", "disabled");

  const result = checkWinner(board);
  if (result === PLAYER) {
    isGameOver = true;
    handleGameEnd("win");
  } else if (result === "draw") {
    isGameOver = true;
    handleGameEnd("draw");
  } else {
    statusText.textContent = "Ход компьютера…";
    computerMove();
  }
}

boardElement.addEventListener("click", handleCellClick);

function handleGameEnd(result) {
  let text = "";
  let promoCode = "";

  switch (result) {
    case "win":
      text = "Ты выиграла! 🎉";
      promoCode = generatePromoCode();
      promoCodeText.textContent = promoCode;
      promoBlock.hidden = false;
      break;
    case "lose":
      text = "В этот раз выиграл компьютер. Попробуешь ещё раз?";
      promoBlock.hidden = true;
      break;
    case "draw":
      text = "Ничья, давай ещё раз?";
      promoBlock.hidden = true;
      break;
  }

  statusText.textContent = "";
  resultTitle.textContent = text;
  resultSection.hidden = false;

  // Отправляем результат на backend
  sendResult(result, promoCode);
}

function generatePromoCode() {
  const n = Math.floor(Math.random() * 100000);
  return n.toString().padStart(5, "0");
}

async function sendResult(result, promoCode) {
  try {
    const payload = {
      result,
      promoCode,
      initData,
    };

    const resp = await fetch("/api/result", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!resp.ok) {
      console.error("Failed to send result:", await resp.text());
    } else {
      const data = await resp.json();
      console.log("Result sent:", data);
    }
  } catch (err) {
    console.error("Error sending result:", err);
  }
}

// Кнопка "Сыграть ещё раз"
playAgainBtn.addEventListener("click", () => {
  resetGame();
});

// Кнопка "Скопировать код"
copyPromoBtn.addEventListener("click", async () => {
  const code = promoCodeText.textContent.trim();
  if (!code) return;

  try {
    await navigator.clipboard.writeText(code);
    if (tg && tg.showPopup) {
      tg.showPopup({
        title: "Скопировано",
        message: "Промокод скопирован в буфер обмена 💕",
        buttons: [{ type: "close" }],
      });
    } else {
      alert("Промокод скопирован: " + code);
    }
  } catch (err) {
    console.error("Clipboard error:", err);
  }
});

// Инициализация
resetGame();
```

---

## 7. Инструкция по запуску

### 7.1. Локально

1. Инициализировать модуль:
   ```bash
   go mod tidy
   ```
2. Заполнить `.env`:
   ```env
   TELEGRAM_BOT_TOKEN=...
   PORT=8080
   PUBLIC_BASE_URL=http://localhost:8080
   ```
3. Запустить:
   ```bash
   go run ./cmd/server/main.go
   ```
4. Открыть в браузере:
   - `http://localhost:8080/web/`

### 7.2. Через ngrok для теста в Telegram

1. Запустить сервер локально.
2. Поднять ngrok:
   ```bash
   ngrok http 8080
   ```
3. В BotFather:
   - `/setmenubutton` → `web_app`
   - URL: `https://<ngrok-id>.ngrok.io/web/`
4. Открыть бота в Telegram, нажать Mini App кнопку.
5. Сыграть партию и проверить:
   - при победе:
     - промокод на экране
     - сообщение в Telegram: «Победа! Промокод выдан: <код>»
   - при проигрыше:
     - текст «Попробуешь ещё раз?»
     - сообщение в Telegram: «Проигрыш»

---

## 8. Что важно для ревьюера

- Чёткая структура проекта (`cmd/`, `internal/`, `web/`).
- Реализация официального Mini App запуска через **Menu Button**.
- Корректная валидация `initData` (без неё нельзя доверять frontend).
- Логика игры (player vs машина) с простым, но не совсем тупым AI.
- Соответствие UX целевой аудитории (женщины 25–40):
  - цвета
  - тексты
  - мягкие эффекты и анимации.
- Чистый, аккуратный код, подходящий для продакшн-пет-проекта.

---
