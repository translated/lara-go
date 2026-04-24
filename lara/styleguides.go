package lara

import "fmt"

type StyleguidesService struct {
	client *Client
}

func newStyleguidesService(client *Client) *StyleguidesService {
	return &StyleguidesService{
		client: client,
	}
}

func (s *StyleguidesService) List() ([]Styleguide, error) {
	var styleguides []Styleguide
	err := s.client.Get("/v2/styleguides", nil, nil, &styleguides)
	if err != nil {
		return nil, fmt.Errorf("failed to list styleguides: %w", err)
	}
	return styleguides, nil
}

func (s *StyleguidesService) Get(id string) (*Styleguide, error) {
	var styleguide Styleguide
	err := s.client.Get(fmt.Sprintf("/v2/styleguides/%s", id), nil, nil, &styleguide)
	if err != nil {
		if laraErr, ok := err.(*LaraError); ok && laraErr.Status == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get styleguide: %w", err)
	}
	return &styleguide, nil
}
