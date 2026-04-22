package checker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jaydeadlondon/project_na_go/internal/models"
	"gorm.io/gorm"
)

type Checker struct {
	db         *gorm.DB
	httpClient *http.Client
}

func NewChecker(db *gorm.DB) *Checker {
	return &Checker{
		db: db,
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

	c.handleIncident(monitor, previousStatus, result.Status)
}

func (c *Checker) handleIncident(
	monitor *models.Monitor,
	previousStatus models.MonitorStatus,
	currentStatus models.CheckStatus,
) {
	now := time.Now()

	if previousStatus != models.StatusDown && currentStatus == models.CheckStatusDown {
		incident := models.Incident{
			MonitorID: monitor.ID,
			StartedAt: now,
		}
		if err := c.db.Create(&incident).Error; err != nil {
			log.Printf("Failed to create incident for monitor %d: %v", monitor.ID, err)
		}
		log.Printf("🔴 Monitor DOWN: %s (%s)", monitor.Name, monitor.URL)
		return
	}

	if previousStatus == models.StatusDown && currentStatus == models.CheckStatusUp {
		var incident models.Incident
		err := c.db.Where("monitor_id = ? AND resolved_at IS NULL", monitor.ID).
			First(&incident).Error

		if err == nil {
			duration := int64(now.Sub(incident.StartedAt).Seconds())
			incident.ResolvedAt = &now
			incident.DurationSeconds = &duration

			if err := c.db.Save(&incident).Error; err != nil {
				log.Printf("Failed to resolve incident for monitor %d: %v", monitor.ID, err)
			}
			log.Printf("🟢 Monitor UP: %s (%s) — was down for %ds", monitor.Name, monitor.URL, duration)
		}
	}
}