// internal/repository/memory.go
package repository

import (
	"log"
	"sync"
	"time"
)

type InMemoryStatsRepo struct {
	mu        sync.RWMutex
	statistics map[string]map[string]int 
	timeline   []EventLogEntry         
}

type EventLogEntry struct {
	Source    string
	Action    string
	Timestamp time.Time
}

func NewInMemoryStatsRepo() *InMemoryStatsRepo {
	return &InMemoryStatsRepo{
		statistics: make(map[string]map[string]int),
		timeline:   []EventLogEntry{},
	}
}

func (r *InMemoryStatsRepo) Save(source, action string, ts string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		parsedTime = time.Now()
	}

	if _, ok := r.statistics[source]; !ok {
		r.statistics[source] = make(map[string]int)
	}
	r.statistics[source][action]++

	r.timeline = append(r.timeline, EventLogEntry{
		Source:    source,
		Action:    action,
		Timestamp: parsedTime,
	})

	log.Printf("Статистика обновлена: [%s %s] = %d\n",
		source, action, r.statistics[source][action])
}

func (r *InMemoryStatsRepo) GetAll() map[string]map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]map[string]int)
	for source, actions := range r.statistics {
		result[source] = make(map[string]int)
		for action, count := range actions {
			result[source][action] = count
		}
	}
	return result
}
