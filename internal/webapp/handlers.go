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
	webhookSecret := os.Getenv("TELEGRAM_WEBHOOK_SECRET")

	// Serve static mini app assets.
	r.StaticFS("/web", http.Dir("web"))
	r.StaticFile("/robots.txt", "web/robots.txt")

	// Redirect root to the mini app to avoid noisy 404s.
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/web/")
	})

	// Telegram webhook endpoint (use instead of long polling).
	r.POST("/api/telegram/webhook", func(c *gin.Context) {
		if webhookSecret != "" {
			if c.GetHeader("X-Telegram-Bot-Api-Secret-Token") != webhookSecret {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		var update bot.Update
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid update"})
			return
		}

		botClient.HandleUpdate(&update)
		c.Status(http.StatusOK)
	})

	// Game result endpoint.
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
