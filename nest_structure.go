package gonest

import "fmt"

type Structure struct {
	Name           string           `json:"name"`
	CountryCode    string           `json:"country_code"`
	TZ             string           `json:"time_zone"`
	Away           string           `json:"away"`
	Thermostats    []string         `json:"thermostats"`
	StructureID    string           `json:"structure_id"`
	RHREnrollement bool             `json:"rhr_enrollment"`
	Wheres         map[string]Where `json:"wheres"`
}

type Where struct {
	WhereID string `json:"where_id"`
	Name    string `json:"name"`
}

func (h *Handler) ReadStructures() (map[string]Structure, error) {
	n, err := h.Get()
	if err != nil {
		return nil, err
	}

	structures := make(map[string]Structure)
	for sid, t := range n.Structures {
		structures[sid] = t
	}
	return structures, nil
}

func (h *Handler) ReadStructure(id string) (Structure, error) {
	n, err := h.Get()
	if err != nil {
		return Structure{}, err
	}

	if s, ok := n.Structures[id]; ok {
		return s, nil
	}
	return Structure{}, fmt.Errorf("unknown structure id: %s", id)
}

func (h *Handler) SetAway(status string) error {
	structures, err := h.ReadStructures()
	if err != nil {
		return err
	}
	for sid := range structures {
		if err := h.Set(fmt.Sprintf("structures/%s", sid), fmt.Sprintf(`{"away": "%s"}`, status)); err != nil {
			return err
		}
	}
	return nil
}
