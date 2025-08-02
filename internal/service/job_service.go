package service

import (
	"github.com/google/uuid"
	"github.com/satyam-svg/resume-parser/internal/model"
	"gorm.io/gorm"
)

type JobService struct {
	DB *gorm.DB
}

func (s *JobService) CreateJob(job *model.Job) error {
	return s.DB.Create(job).Error
}

func (s *JobService) GetJobs() ([]model.Job, error) {
	var jobs []model.Job
	err := s.DB.Order("created_at desc").Find(&jobs).Error
	return jobs, err
}

func (js *JobService) GetJobsByRecruiterID(recruiterID string) ([]model.Job, error) {
	var jobs []model.Job
	err := js.DB.Where("recruiter_id = ?", recruiterID).Find(&jobs).Error
	return jobs, err
}

func (js *JobService) GetJobsByRecruiter(recruiterID string) ([]model.Job, error) {
	var jobs []model.Job
	if err := js.DB.Where("recruiter_id = ?", recruiterID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (js *JobService) GetJobByID(id uint) (*model.Job, error) {
	var job model.Job
	if err := js.DB.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *JobService) GetRelevantStudents() ([]model.User, error) {
	var users []model.User
	if err := s.DB.Where("role = ?", "applicant").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (js *JobService) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	if err := js.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (js *JobService) GetAllJobs() ([]model.Job, error) {
	var jobs []model.Job
	if err := js.DB.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (js *JobService) UpdateUserCredits(userID uuid.UUID, credits int) error {
	return js.DB.Model(&model.User{}).Where("id = ?", userID).Update("credits", credits).Error
}
