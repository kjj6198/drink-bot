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
	DataSource  string   `json:"data_source"`
}

type Dialog struct {
	CallbackID  string    `json:"callback_id"`
	Title       string    `json:"title"`
	SubmitLabel string    `json:"submit_label"`
	Elements    []Element `json:"elements"`
}

type DialogOptions struct {
	Dialog    Dialog `json:"dialog"`
	Token     string `json:"token" binding:"required"`
	TriggerID string `json:"trigger_id"`
}

type SlackDialog struct {
	Type      string `json:"type" form:"type"`
	Token     string `json:"token" form:"token"`
	TimeStamp int64  `json:"action_ts"`
	User      string `json:"user" form:"user"`
}

type SlackUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackDialogParams struct {
	Type         string                 `json:"type" form:"type"`
	TimeStamp    string                 `json:"ts" form:"ts"`
	User         *SlackUser             `json:"user" form:"user"`
	SlackChannel *SlackChannel          `json:"channel" form:"channel"`
	Submission   map[string]interface{} `json:"submission" form:"submission"`
	CallbackID   string                 `json:"callback_id" form:"callback_id"`
}

type SlackDialogPayload struct {
	Payload *SlackDialogParams `json:"payload" form:"payload"`
}

func CreateMenuDialog(text string, triggerID string) DialogOptions {
	return DialogOptions{
		TriggerID: triggerID,
		Dialog: Dialog{
			CallbackID:  "submit-menu",
			Title:       "建立飲料訂單",
			SubmitLabel: "建立飲料訂單",
			Elements: []Element{
				Element{
					Label:    "訂單名稱",
					Name:     "name",
					Value:    text,
					Type:     "text",
					Hint:     "please select drink_shop",
					Optional: false,
				},
				Element{
					Label:      "店家名稱",
					Name:       "drink_shop",
					Type:       "select",
					Hint:       "please select drink_shop",
					DataSource: "external",
					Optional:   false,
				},
				Element{
					Label:       "時間（分）",
					Name:        "duration",
					Type:        "text",
					SubType:     "number",
					Placeholder: "e.g: 900",
					Hint:        "單位為(分)",
					Optional:    false,
				},
				Element{
					Label:      "notify to channel",
					Name:       "channel",
					Type:       "select",
					DataSource: "channels",
					Optional:   false,
				},
			},
		},
	}
}
