package service

import "github.com/BogdanBratsky/proverb/internal/model"

type ProverbService struct {
	proverbs []model.Proverb
}

func NewProverbService() *ProverbService {
	return &ProverbService{proverbs: make([]model.Proverb, 0)}
}

func (s *ProverbService) GetProverb() {}
