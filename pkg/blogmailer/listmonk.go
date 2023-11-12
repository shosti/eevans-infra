package blogmailer

import "encoding/json"

// Campaign represents a Listmonk campaign
type Campaign struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	Lists       []List   `json:"lists,omitempty"`
	Body        string   `json:"body,omitempty"`
	Altbody     string   `json:"altbody,omitempty"`
	SendAt      string   `json:"send_at,omitempty"`
}

// CampaignStatusUpdate represents a listmonk campaign status update
type CampaignStatusUpdate struct {
	Status string `json:"status"`
}

const CampaignStatusScheduled = "scheduled"

// CreateCampaignResponse represents the response from a CreateCampaign call
type CreateCampaignResponse struct {
	Data *Campaign `json:"data"`
}

// List represents a Listmonk list
type List struct {
	ID int
}

// Bizarrely, the Listmonk API accepts integers but returns objects for list IDs
// so we have to do some hackery
func (l List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.ID)
}

type listObject struct {
	ID int `json:"id"`
}

func (l *List) UnmarshalJSON(data []byte) error {
	var lo listObject
	err := json.Unmarshal(data, &lo)
	if err == nil {
		l.ID = lo.ID
		return nil
	}
	var id int
	err = json.Unmarshal(data, &id)
	if err == nil {
		l.ID = id
		return nil
	}
	return err
}
