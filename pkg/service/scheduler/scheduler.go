package scheduler

import (
	"avito2/pkg/repository"
	"log"
	"time"
)

type Scheduler struct {
	userSegmentRepository *repository.UserSegmentRepository
}

func NewScheduler(userSegmentRepository *repository.UserSegmentRepository) *Scheduler {
	return &Scheduler{
		userSegmentRepository: userSegmentRepository,
	}
}

func (s *Scheduler) ScheduleDeletion() {

	go func() {
		for {
			log.Default().Println("Delete cron job iteration started...")
			if err := (*s.userSegmentRepository).RemoveExpiredUserSegments(); err != nil {
				log.Default().Println("Something went wrong due cron-deleting")
			}

			log.Default().Println("Delete cron job iteration ended...")
			time.Sleep(40 * time.Second)
		}
	}()
}
