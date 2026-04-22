package scheduler

import (
	"fmt"
	"log"
	"sync"

	"github.com/jaydeadlondon/project_na_go/internal/checker"
	"github.com/jaydeadlondon/project_na_go/internal/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Scheduler struct {
	db      *gorm.DB
	checker *checker.Checker
	cron    *cron.Cron
	jobs    map[uint]cron.EntryID
	mu      sync.Mutex
}

func NewScheduler(db *gorm.DB) *Scheduler {
	return &Scheduler{
		db:      db,
		checker: checker.NewChecker(db),
		cron:    cron.New(),
		jobs:    make(map[uint]cron.EntryID),
	}
}

func (s *Scheduler) Start() {
	log.Println("🕐 Starting scheduler...")

	var monitors []models.Monitor
	if err := s.db.Where("is_active = ?", true).Find(&monitors).Error; err != nil {
		log.Printf("Failed to load monitors: %v", err)
		return
	}

	for _, monitor := range monitors {
		s.AddMonitor(monitor)
	}

	s.cron.Start()
	log.Printf("✅ Scheduler started with %d monitors", len(monitors))
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) AddMonitor(monitor models.Monitor) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.jobs[monitor.ID]; exists {
		s.cron.Remove(entryID)
	}

	cronExpr := fmt.Sprintf("@every %dm", monitor.Interval)

	monitorCopy := monitor

	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.runCheck(&monitorCopy)
	})

	if err != nil {
		log.Printf("Failed to add monitor %d to scheduler: %v", monitor.ID, err)
		return
	}

	s.jobs[monitor.ID] = entryID
	log.Printf("📡 Monitor scheduled: %s (every %d min)", monitor.Name, monitor.Interval)
}

func (s *Scheduler) RemoveMonitor(monitorID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.jobs[monitorID]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, monitorID)
		log.Printf("Monitor %d removed from scheduler", monitorID)
	}
}

func (s *Scheduler) runCheck(monitor *models.Monitor) {
	var freshMonitor models.Monitor
	if err := s.db.First(&freshMonitor, monitor.ID).Error; err != nil {
		log.Printf("Monitor %d not found, skipping", monitor.ID)
		return
	}

	if !freshMonitor.IsActive {
		return
	}

	result := s.checker.CheckMonitor(freshMonitor)
	s.checker.SaveResult(&freshMonitor, result)
}

func (s *Scheduler) CheckNow(monitorID uint) {
	var monitor models.Monitor
	if err := s.db.First(&monitor, monitorID).Error; err != nil {
		return
	}
	go s.runCheck(&monitor)
}