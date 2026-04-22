package handlers

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
	"github.com/jaydeadlondon/project_na_go/internal/models"
	"gorm.io/gorm"
)

type TelegramHandler struct {
	db      *gorm.DB
	botName string
}

func NewTelegramHandler(db *gorm.DB, botName string) *TelegramHandler {
	return &TelegramHandler{db: db, botName: botName}
}

func (h *TelegramHandler) GenerateLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	token, err := generateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	if err := h.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("link_token", token).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save token",
		})
	}

	botLink := "https://t.me/" + h.botName + "?start=" + token

	return c.JSON(fiber.Map{
		"link":  botLink,
		"token": token,
	})
}

func (h *TelegramHandler) Unlink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	if err := h.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"telegram_chat_id": nil,
			"link_token":       nil,
		}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unlink telegram",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Telegram unlinked successfully",
	})
}

func (h *TelegramHandler) Status(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"linked":           user.TelegramChatID != nil,
		"telegram_chat_id": user.TelegramChatID,
	})
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}