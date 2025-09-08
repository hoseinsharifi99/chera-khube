package model

type DivarWidget struct {
	Widgets []interface{} `json:"widgets"`
}

type EventWidget struct {
	EventRow struct {
		Title      string `json:"title"`
		Subtitle   string `json:"subtitle"`
		Label      string `json:"label"`
		ImageId    string `json:"image_id"`
		HasDivider bool   `json:"has_divider"`
		IconName   string `json:"icon_name"`
	} `json:"event_row"`
}

type DescriptionWidget struct {
	DescriptionRow DescriptionRow `json:"description_row"`
}

type DescriptionRow struct {
	Text       string `json:"text"`
	HasDivider bool   `json:"has_divider"`
	Expandable bool   `json:"expandable"`
}
