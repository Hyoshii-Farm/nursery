package predator

func (s *Service) GetReport(startDate, endDate string) ([]string, error) {
	predators, err := s.repo.GetReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, r := range predators {
		result = append(result, r.Name)
	}

	return result, nil
}
