package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jaydeadlondon/project_na_go/internal/models"
	"github.com/jaydeadlondon/project_na_go/internal/scheduler"
	"gorm.io/gorm"
)

type MonitorHandler struct {
	db        *gorm.DB
	scheduler *scheduler.Scheduler
}

func NewMonitorHandler(db *gorm.DB, scheduler *scheduler.Scheduler) *MonitorHandler {
	return &MonitorHandler{db: db, scheduler: scheduler}
}

type CreateMonitorRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

type UpdateMonitorRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
	IsActive *bool  `json:"is_active"`
}

func (h *MonitorHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	req := new(CreateMonitorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	if req.Interval <= 0 {
		req.Interval = 5
	}

	if req.Interval > 60 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Interval must be between 1 and 60 minutes",
		})
	}

	monitor := models.Monitor{
		UserID:        userID,
		Name:          req.Name,
		URL:           req.URL,
		Interval:      req.Interval,
		IsActive:      true,
		CurrentStatus: models.StatusPending,
	}

	if err := h.db.Create(&monitor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create monitor",
		})
	}

	h.scheduler.AddMonitor(monitor)

	h.scheduler.CheckNow(monitor.ID)

	return c.Status(fiber.StatusCreated).JSON(monitor)
}

func (h *MonitorHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var monitors []models.Monitor
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&monitors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get monitors",
		})
	}

	return c.JSON(fiber.Map{
		"data":  monitors,
		"total": len(monitors),
	})
}

func (h *MonitorHandler) GetOne(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	return c.JSON(monitor)
}

func (h *MonitorHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	req := new(UpdateMonitorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name != "" {
		monitor.Name = req.Name
	}
	if req.URL != "" {
		monitor.URL = req.URL
	}
	if req.Interval > 0 && req.Interval <= 60 {
		monitor.Interval = req.Interval
	}
	if req.IsActive != nil {
		monitor.IsActive = *req.IsActive
	}

	if err := h.db.Save(&monitor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update monitor",
		})
	}

	if monitor.IsActive {
		h.scheduler.AddMonitor(monitor)
	} else {
		h.scheduler.RemoveMonitor(monitor.ID)
	}

	return c.JSON(monitor)
}

func (h *MonitorHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	h.scheduler.RemoveMonitor(monitor.ID)

	h.db.Where("monitor_id = ?", monitor.ID).Delete(&models.MonitorCheck{})
	h.db.Where("monitor_id = ?", monitor.ID).Delete(&models.Incident{})

	if err := h.db.Delete(&monitor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete monitor",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Monitor deleted successfully",
	})
}

func (h *MonitorHandler) GetChecks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	var checks []models.MonitorCheck
	if err := h.db.Where("monitor_id = ?", monitorID).
		Order("checked_at DESC").
		Limit(100).
		Find(&checks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get checks",
		})
	}

	return c.JSON(fiber.Map{
		"data":  checks,
		"total": len(checks),
	})
}

func (h *MonitorHandler) GetStats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	monitorID := c.Params("id")

	var monitor models.Monitor
	if err := h.db.Where("id = ? AND user_id = ?", monitorID, userID).
		First(&monitor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Monitor not found",
		})
	}

	stats := fiber.Map{
		"uptime_24h":        calculateUptime(h.db, monitor.ID, 24),
		"uptime_7d":         calculateUptime(h.db, monitor.ID, 24*7),
		"uptime_30d":        calculateUptime(h.db, monitor.ID, 24*30),
		"avg_response_time": calculateAvgResponseTime(h.db, monitor.ID),
	}

	return c.JSON(stats)
}

func calculateUptime(db *gorm.DB, monitorID uint, hours int) float64 {
	var total int64
	var up int64

	db.Model(&models.MonitorCheck{}).
		Where("monitor_id = ? AND checked_at > NOW() - INTERVAL '? hours'", monitorID, hours).
		Count(&total)

	if total == 0 {
		return 100.0
	}

	db.Model(&models.MonitorCheck{}).
		Where("monitor_id = ? AND status = ? AND checked_at > NOW() - INTERVAL '? hours'",
			monitorID, models.CheckStatusUp, hours).
		Count(&up)

	return float64(up) / float64(total) * 100
}

func calculateAvgResponseTime(db *gorm.DB, monitorID uint) float64 {
	var result struct {
		Avg float64
	}

	db.Model(&models.MonitorCheck{}).
		Select("AVG(response_time_ms) as avg").
		Where("monitor_id = ? AND checked_at > NOW() - INTERVAL '24 hours'", monitorID).
		Scan(&result)

	return result.Avg
}