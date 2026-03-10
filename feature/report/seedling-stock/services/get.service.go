package seedlingstock

func (s *Service) GetReport(startDate, endDate string) ([]string, error) {
	seedling_stocks, err := s.repo.GetReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, r := range seedling_stocks {
		result = append(result, r.Name)
	}

	return result, nil
}
