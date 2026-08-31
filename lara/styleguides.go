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

func (s *StyleguidesService) Create(name, content string) (*Styleguide, error) {
	body := map[string]interface{}{
		"name":    name,
		"content": content,
	}

	var styleguide Styleguide
	err := s.client.Post("/v2/styleguides", body, nil, nil, &styleguide)
	if err != nil {
		return nil, fmt.Errorf("failed to create styleguide: %w", err)
	}
	return &styleguide, nil
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

func (s *StyleguidesService) Delete(id string) (*Styleguide, error) {
	var styleguide Styleguide
	err := s.client.Delete(fmt.Sprintf("/v2/styleguides/%s", id), nil, nil, &styleguide)
	if err != nil {
		return nil, fmt.Errorf("failed to delete styleguide: %w", err)
	}
	return &styleguide, nil
}

func (s *StyleguidesService) Update(id string, name, content *string) (*Styleguide, error) {
	body := map[string]interface{}{}
	if name != nil {
		body["name"] = *name
	}
	if content != nil {
		body["content"] = *content
	}

	var styleguide Styleguide
	err := s.client.Put(fmt.Sprintf("/v2/styleguides/%s", id), body, nil, nil, &styleguide)
	if err != nil {
		return nil, fmt.Errorf("failed to update styleguide: %w", err)
	}
	return &styleguide, nil
}

func (s *StyleguidesService) GetShares(id string) (*StyleguideShares, error) {
	var shares StyleguideShares
	err := s.client.Get(fmt.Sprintf("/v2/styleguides/%s/shares", id), nil, nil, &shares)
	if err != nil {
		return nil, fmt.Errorf("failed to get styleguide shares: %w", err)
	}
	return &shares, nil
}

func (s *StyleguidesService) AddAccountShare(id string) (*Styleguide, error) {
	return s.share("POST", fmt.Sprintf("/v2/styleguides/%s/shares", id), nil)
}

func (s *StyleguidesService) AddAccountShareWithName(id, name string) (*Styleguide, error) {
	return s.share("POST", fmt.Sprintf("/v2/styleguides/%s/shares", id), &name)
}

func (s *StyleguidesService) RenameAccountShare(id, name string) (*Styleguide, error) {
	return s.share("PUT", fmt.Sprintf("/v2/styleguides/%s/shares", id), &name)
}

func (s *StyleguidesService) RevokeAccountShare(id string) (*Styleguide, error) {
	return s.share("DELETE", fmt.Sprintf("/v2/styleguides/%s/shares", id), nil)
}

func (s *StyleguidesService) AddGroupShare(id, groupID string) (*Styleguide, error) {
	return s.share("POST", fmt.Sprintf("/v2/styleguides/%s/shares/groups/%s", id, groupID), nil)
}

func (s *StyleguidesService) AddGroupShareWithName(id, groupID, name string) (*Styleguide, error) {
	return s.share("POST", fmt.Sprintf("/v2/styleguides/%s/shares/groups/%s", id, groupID), &name)
}

func (s *StyleguidesService) RenameGroupShare(id, groupID, name string) (*Styleguide, error) {
	return s.share("PUT", fmt.Sprintf("/v2/styleguides/%s/shares/groups/%s", id, groupID), &name)
}

func (s *StyleguidesService) RevokeGroupShare(id, groupID string) (*Styleguide, error) {
	return s.share("DELETE", fmt.Sprintf("/v2/styleguides/%s/shares/groups/%s", id, groupID), nil)
}

func (s *StyleguidesService) share(method, path string, name *string) (*Styleguide, error) {
	body := map[string]interface{}{}
	if name != nil {
		body["name"] = *name
	}
	var styleguide Styleguide
	var err error
	switch method {
	case "POST":
		err = s.client.Post(path, body, nil, nil, &styleguide)
	case "PUT":
		err = s.client.Put(path, body, nil, nil, &styleguide)
	case "DELETE":
		err = s.client.Delete(path, nil, nil, &styleguide)
	default:
		return nil, fmt.Errorf("unsupported method %q for styleguide share", method)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update styleguide share: %w", err)
	}
	return &styleguide, nil
}
