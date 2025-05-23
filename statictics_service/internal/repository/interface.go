package repository

type StatsRepository interface {
	Save(source, action, timestamp string) error
	GetAll() (map[string]map[string]int, error)
}
