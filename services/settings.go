package services

import "strconv"

type SettingsService struct {
	db *DatabaseService
}

func NewSettingsService(db *DatabaseService) *SettingsService {
	return &SettingsService{db: db}
}

func (s *SettingsService) GetPollingInterval() int {
	v, err := s.db.GetSetting("polling_interval")
	if err != nil {
		return 300
	}
	n, _ := strconv.Atoi(v)
	if n <= 0 {
		return 300
	}
	return n
}

func (s *SettingsService) SetPollingInterval(sec int) error {
	return s.db.SetSetting("polling_interval", strconv.Itoa(sec))
}

func (s *SettingsService) GetTheme() string {
	v, err := s.db.GetSetting("theme")
	if err != nil {
		return "dark"
	}
	return v
}

func (s *SettingsService) SetTheme(t string) error {
	return s.db.SetSetting("theme", t)
}

func (s *SettingsService) GetWindowState() (int, int) {
	w, _ := strconv.Atoi(getS(s.db, "window_width", "1024"))
	h, _ := strconv.Atoi(getS(s.db, "window_height", "768"))
	return w, h
}

func (s *SettingsService) SaveWindowState(w, h int) {
	s.db.SetSetting("window_width", strconv.Itoa(w))
	s.db.SetSetting("window_height", strconv.Itoa(h))
}

func getS(db *DatabaseService, key, fallback string) string {
	v, err := db.GetSetting(key)
	if err != nil {
		return fallback
	}
	return v
}
