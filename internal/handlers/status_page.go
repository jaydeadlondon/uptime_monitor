package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jaydeadlondon/project_na_go/internal/models"
	"gorm.io/gorm"
	"time"
)

type StatusPageHandler struct {
	db *gorm.DB
}

func NewStatusPageHandler(db *gorm.DB) *StatusPageHandler {
	return &StatusPageHandler{db: db}
}

type PublicMonitor struct {
	Name          string              `json:"name"`
	URL           string              `json:"url"`
	CurrentStatus models.MonitorStatus `json:"current_status"`
	Uptime24h     float64             `json:"uptime_24h"`
	Uptime7d      float64             `json:"uptime_7d"`
	AvgResponseMs float64             `json:"avg_response_ms"`
	LastCheckedAt *time.Time          `json:"last_checked_at"`
}

func (h *StatusPageHandler) GetStatusPage(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Status page not found",
		})
	}

	var monitors []models.Monitor
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).
		Find(&monitors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get monitors",
		})
	}

	publicMonitors := make([]PublicMonitor, 0, len(monitors))

	for _, m := range monitors {
		publicMonitors = append(publicMonitors, PublicMonitor{
			Name:          m.Name,
			URL:           m.URL,
			CurrentStatus: m.CurrentStatus,
			Uptime24h:     calculateUptime(h.db, m.ID, 24),
			Uptime7d:      calculateUptime(h.db, m.ID, 24*7),
			AvgResponseMs: calculateAvgResponseTime(h.db, m.ID),
			LastCheckedAt: m.LastCheckedAt,
		})
	}

	systemStatus := "operational"
	for _, m := range monitors {
		if m.CurrentStatus == models.StatusDown {
			systemStatus = "degraded"
			break
		}
	}

	return c.JSON(fiber.Map{
		"system_status": systemStatus,
		"monitors":      publicMonitors,
		"total":         len(publicMonitors),
	})
}