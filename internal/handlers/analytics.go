package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jaydeadlondon/project_na_go/internal/models"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	db *gorm.DB
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

func (h *AnalyticsHandler) Overview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var monitors []models.Monitor
	if err := h.db.Where("user_id = ?", userID).Find(&monitors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get monitors",
		})
	}

	totalMonitors := len(monitors)
	upCount := 0
	downCount := 0
	pendingCount := 0

	for _, m := range monitors {
		switch m.CurrentStatus {
		case models.StatusUp:
			upCount++
		case models.StatusDown:
			downCount++
		default:
			pendingCount++
		}
	}

	var totalChecks int64
	var upChecks int64

	h.db.Model(&models.MonitorCheck{}).
		Joins("JOIN monitors ON monitors.id = monitor_checks.monitor_id").
		Where("monitors.user_id = ? AND monitor_checks.checked_at > ?",
			userID, time.Now().Add(-24*time.Hour)).
		Count(&totalChecks)

	h.db.Model(&models.MonitorCheck{}).
		Joins("JOIN monitors ON monitors.id = monitor_checks.monitor_id").
		Where("monitors.user_id = ? AND monitor_checks.status = ? AND monitor_checks.checked_at > ?",
			userID, models.CheckStatusUp, time.Now().Add(-24*time.Hour)).
		Count(&upChecks)

	overallUptime := 100.0
	if totalChecks > 0 {
		overallUptime = float64(upChecks) / float64(totalChecks) * 100
	}

	var activeIncidents int64
	h.db.Model(&models.Incident{}).
		Joins("JOIN monitors ON monitors.id = incidents.monitor_id").
		Where("monitors.user_id = ? AND incidents.resolved_at IS NULL", userID).
		Count(&activeIncidents)

	return c.JSON(fiber.Map{
		"total_monitors":   totalMonitors,
		"up":               upCount,
		"down":             downCount,
		"pending":          pendingCount,
		"overall_uptime":   overallUptime,
		"active_incidents": activeIncidents,
	})
}

func (h *AnalyticsHandler) MonitorChart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")
	period := c.Query("period", "24h")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	var since time.Time
	var groupBy string

	switch period {
	case "7d":
		since = time.Now().Add(-7 * 24 * time.Hour)
		groupBy = "hour"
	case "30d":
		since = time.Now().Add(-30 * 24 * time.Hour)
		groupBy = "day"
	default: // 24h
		since = time.Now().Add(-24 * time.Hour)
		groupBy = "hour"
	}

	type ChartPoint struct {
		Time           time.Time `json:"time"`
		AvgResponseMs  float64   `json:"avg_response_ms"`
		UptimePercent  float64   `json:"uptime_percent"`
		TotalChecks    int64     `json:"total_checks"`
	}

	var points []ChartPoint

	rows, err := h.db.Raw(`
		SELECT
			date_trunc(?, checked_at) as time,
			AVG(response_time_ms) as avg_response_ms,
			COUNT(*) as total_checks,
			SUM(CASE WHEN status = 'up' THEN 1 ELSE 0 END) * 100.0 / COUNT(*) as uptime_percent
		FROM monitor_checks
		WHERE monitor_id = ? AND checked_at > ?
		GROUP BY date_trunc(?, checked_at)
		ORDER BY time ASC
	`, groupBy, monitor.ID, since, groupBy).Rows()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get chart data",
		})
	}
	defer rows.Close()

	for rows.Next() {
		var point ChartPoint
		if err := rows.Scan(
			&point.Time,
			&point.AvgResponseMs,
			&point.TotalChecks,
			&point.UptimePercent,
		); err != nil {
			continue
		}
		points = append(points, point)
	}

	if points == nil {
		points = []ChartPoint{}
	}

	return c.JSON(fiber.Map{
		"monitor": monitor,
		"period":  period,
		"data":    points,
	})
}

func (h *AnalyticsHandler) MonitorIncidents(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	var incidents []models.Incident
	if err := h.db.Where("monitor_id = ?", monitorID).
		Order("started_at DESC").
		Limit(50).
		Find(&incidents).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get incidents",
		})
	}

	if incidents == nil {
		incidents = []models.Incident{}
	}

	return c.JSON(fiber.Map{
		"data":  incidents,
		"total": len(incidents),
	})
}