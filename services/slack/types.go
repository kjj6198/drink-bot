package slack

type SlackMessageInput struct {
	UserID      string `form:"user_id"`
	UserName    string `form:"user_name"`
	Command     string `form:"command"`
	Text        string `form:"text"`
	TriggerID   string `form:"trigger_id"`
	ChannelID   string `form:"channel_id"`
	Token       string `form:"token"`
	ResponseURL string `form:"response_url"`
}

type Option struct {
	Label string `json:"label" form:"label"`
	Value string `json:"value" form:"value"`
}

type Element struct {
	Label       string   `json:"label" form:"label"`
	Name        string   `json:"name" form:"name"`
	Type        string   `json:"type" form:"type"`       // text, textarea, select
	SubType     string   `json:"subtype" form:"subtype"` // number,
	Placeholder string   `json:"placeholder" form:"placeholder"`
	Value       string   `json:"value" form:"value"`
	Hint        string   `json:"hint" form:"hint"`
	Optional    bool     `json:"optional,omitempty" form:"optional,omitempty"`
	Options     []Option `json:"options,omitempty"`
}

type Dialog struct {
	CallbackID  string `json:"callback_id"`
	Title       string `json:"title"`
	SubmitLabel string `json:"submit_label"`
	Elements    []Element
}

type DialogOptions struct {
	Dialog    Dialog `json:"dialog"`
	Token     string `json:"token" binding:"required"`
	TriggerID string `json:"trigger_id"`
}
