package telegram

import (
	"fmt"
	"log"
	"strings"

	"github.com/jaydeadlondon/project_na_go/internal/models"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type Bot struct {
	bot *tele.Bot
	db  *gorm.DB
}

func NewBot(token string, db *gorm.DB) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is empty")
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	b := &Bot{bot: bot, db: db}

	b.registerHandlers()

	return b, nil
}

func (b *Bot) Start() {
	log.Println("✅ Telegram bot started")
	go b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
	log.Println("Telegram bot stopped")
}

func (b *Bot) registerHandlers() {
	b.bot.Handle("/start", b.handleStart)

	b.bot.Handle("/status", b.handleStatus)

	b.bot.Handle("/help", b.handleHelp)
}

func (b *Bot) handleStart(c tele.Context) error {
	args := c.Args()

	if len(args) == 0 {
		return c.Send(
			"👋 Привет! Я бот для мониторинга сайтов.\n\n" +
				"Чтобы получать уведомления, привяжи свой аккаунт:\n" +
				"1. Войди в приложение\n" +
				"2. Перейди в Настройки → Telegram\n" +
				"3. Нажми «Привязать Telegram»\n\n" +
				"Или используй команду /help",
		)
	}

	linkToken := args[0]
	var user models.User

	if err := b.db.Where("link_token = ?", linkToken).First(&user).Error; err != nil {
		return c.Send("❌ Неверный токен. Попробуй снова.")
	}

	chatID := c.Chat().ID
	if err := b.db.Model(&user).Updates(map[string]interface{}{
		"telegram_chat_id": chatID,
		"link_token":       nil,
	}).Error; err != nil {
		return c.Send("❌ Ошибка привязки. Попробуй снова.")
	}

	return c.Send(
		"✅ Аккаунт успешно привязан!\n\n" +
			"Теперь ты будешь получать уведомления когда:\n" +
			"🔴 Сайт упадёт\n" +
			"🟢 Сайт восстановится\n\n" +
			"Используй /status чтобы посмотреть статус мониторов",
	)
}

func (b *Bot) handleStatus(c tele.Context) error {
	chatID := c.Chat().ID

	var user models.User
	if err := b.db.Where("telegram_chat_id = ?", chatID).First(&user).Error; err != nil {
		return c.Send(
			"❌ Аккаунт не привязан.\n" +
				"Используй /start чтобы привязать аккаунт.",
		)
	}

	var monitors []models.Monitor
	if err := b.db.Where("user_id = ?", user.ID).Find(&monitors).Error; err != nil {
		return c.Send("❌ Ошибка получения мониторов.")
	}

	if len(monitors) == 0 {
		return c.Send("📭 У тебя пока нет мониторов.")
	}

	var sb strings.Builder
	sb.WriteString("📊 *Статус мониторов:*\n\n")

	for _, m := range monitors {
		emoji := "⏳"
		switch m.CurrentStatus {
		case models.StatusUp:
			emoji = "🟢"
		case models.StatusDown:
			emoji = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%s *%s*\n`%s`\n\n", emoji, m.Name, m.URL))
	}

	return c.Send(sb.String(), tele.ModeMarkdown)
}

func (b *Bot) handleHelp(c tele.Context) error {
	return c.Send(
		"🤖 *Uptime Monitor Bot*\n\n" +
			"*Команды:*\n" +
			"/start — привязать аккаунт\n" +
			"/status — статус мониторов\n" +
			"/help — помощь\n\n" +
			"*Как привязать аккаунт:*\n" +
			"Перейди в настройки приложения и нажми «Привязать Telegram»",
		tele.ModeMarkdown,
	)
}

func (b *Bot) SendDownAlert(chatID int64, monitor models.Monitor, errMsg string) {
	text := fmt.Sprintf(
		"🔴 *Сайт недоступен!*\n\n"+
			"*Монитор:* %s\n"+
			"*URL:* `%s`\n"+
			"*Причина:* %s",
		monitor.Name,
		monitor.URL,
		errMsg,
	)

	chat := &tele.Chat{ID: chatID}
	if _, err := b.bot.Send(chat, text, tele.ModeMarkdown); err != nil {
		log.Printf("Failed to send down alert: %v", err)
	}
}

func (b *Bot) SendUpAlert(chatID int64, monitor models.Monitor, durationSeconds int64) {
	duration := formatDuration(durationSeconds)

	text := fmt.Sprintf(
		"🟢 *Сайт восстановлен!*\n\n"+
			"*Монитор:* %s\n"+
			"*URL:* `%s`\n"+
			"*Время простоя:* %s",
		monitor.Name,
		monitor.URL,
		duration,
	)

	chat := &tele.Chat{ID: chatID}
	if _, err := b.bot.Send(chat, text, tele.ModeMarkdown); err != nil {
		log.Printf("Failed to send up alert: %v", err)
	}
}

func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%d сек", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%d мин", seconds/60)
	}
	return fmt.Sprintf("%d ч %d мин", seconds/3600, (seconds%3600)/60)
}