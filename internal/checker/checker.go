package checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jaydeadlondon/project_na_go/internal/models"
	"github.com/jaydeadlondon/project_na_go/internal/telegram"
	"gorm.io/gorm"
)

type Checker struct {
	db         *gorm.DB
	httpClient *http.Client
	tgBot      *telegram.Bot
}

func NewChecker(db *gorm.DB, tgBot *telegram.Bot) *Checker {
	return &Checker{
		db:    db,
		tgBot: tgBot,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

type CheckResult struct {
	Status         models.CheckStatus
	StatusCode     *int
	ResponseTimeMs *int64
	ErrorMessage   *string
}

func (c *Checker) CheckMonitor(monitor models.Monitor) CheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		monitor.URL,
		nil,
	)
	if err != nil {
		errMsg := err.Error()
		return CheckResult{
			Status:       models.CheckStatusDown,
			ErrorMessage: &errMsg,
		}
	}

	req.Header.Set("User-Agent", "UptimeMonitor/1.0")

	resp, err := c.httpClient.Do(req)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		errMsg := err.Error()
		return CheckResult{
			Status:         models.CheckStatusDown,
			ResponseTimeMs: &responseTime,
			ErrorMessage:   &errMsg,
		}
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	status := models.CheckStatusUp

	if statusCode >= 400 {
		status = models.CheckStatusDown
	}

	return CheckResult{
		Status:         status,
		StatusCode:     &statusCode,
		ResponseTimeMs: &responseTime,
	}
}

func (c *Checker) SaveResult(monitor *models.Monitor, result CheckResult) {
	now := time.Now()

	check := models.MonitorCheck{
		MonitorID:      monitor.ID,
		Status:         result.Status,
		StatusCode:     result.StatusCode,
		ResponseTimeMs: result.ResponseTimeMs,
		ErrorMessage:   result.ErrorMessage,
		CheckedAt:      now,
	}

	if err := c.db.Create(&check).Error; err != nil {
		log.Printf("Failed to save check for monitor %d: %v", monitor.ID, err)
		return
	}

	previousStatus := monitor.CurrentStatus

	monitor.CurrentStatus = models.MonitorStatus(result.Status)
	monitor.LastCheckedAt = &now

	if err := c.db.Save(monitor).Error; err != nil {
		log.Printf("Failed to update monitor %d status: %v", monitor.ID, err)
		return
	}

	c.handleIncident(monitor, previousStatus, result)
}

func (c *Checker) handleIncident(
	monitor *models.Monitor,
	previousStatus models.MonitorStatus,
	result CheckResult,
) {
	now := time.Now()

	if previousStatus != models.StatusDown && result.Status == models.CheckStatusDown {
		incident := models.Incident{
			MonitorID: monitor.ID,
			StartedAt: now,
		}
		if err := c.db.Create(&incident).Error; err != nil {
			log.Printf("Failed to create incident: %v", err)
		}

		log.Printf("🔴 Monitor DOWN: %s (%s)", monitor.Name, monitor.URL)

		c.sendDownNotification(monitor, result)
		return
	}

	if previousStatus == models.StatusDown && result.Status == models.CheckStatusUp {
		var incident models.Incident
		err := c.db.Where("monitor_id = ? AND resolved_at IS NULL", monitor.ID).
			First(&incident).Error

		if err == nil {
			duration := int64(now.Sub(incident.StartedAt).Seconds())
			incident.ResolvedAt = &now
			incident.DurationSeconds = &duration

			if err := c.db.Save(&incident).Error; err != nil {
				log.Printf("Failed to resolve incident: %v", err)
			}

			log.Printf("🟢 Monitor UP: %s — was down for %ds", monitor.Name, duration)

			c.sendUpNotification(monitor, duration)
		}
	}
}

func (c *Checker) sendDownNotification(monitor *models.Monitor, result CheckResult) {
	if c.tgBot == nil {
		return
	}

	var user models.User
	if err := c.db.First(&user, monitor.UserID).Error; err != nil {
		return
	}

	if user.TelegramChatID == nil {
		return
	}

	errMsg := "Неизвестная ошибка"
	if result.ErrorMessage != nil {
		errMsg = *result.ErrorMessage
	} else if result.StatusCode != nil {
		errMsg = fmt.Sprintf("HTTP %d", *result.StatusCode)
	}

	c.tgBot.SendDownAlert(*user.TelegramChatID, *monitor, errMsg)
}

func (c *Checker) sendUpNotification(monitor *models.Monitor, duration int64) {
	if c.tgBot == nil {
		return
	}

	var user models.User
	if err := c.db.First(&user, monitor.UserID).Error; err != nil {
		return
	}

	if user.TelegramChatID == nil {
		return
	}

	c.tgBot.SendUpAlert(*user.TelegramChatID, *monitor, duration)
}