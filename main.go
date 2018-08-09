package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// AttachmentAction is a button or menu to be included in the attachment. Required when
// using message buttons or menus and otherwise not useful. A maximum of 5 actions may be
// provided per attachment.
type AttachmentAction struct {
	Name            string                        `json:"name"`                       // Required.
	Text            string                        `json:"text"`                       // Required.
	Style           string                        `json:"style,omitempty"`            // Optional. Allowed values: "default", "primary", "danger".
	Type            string                        `json:"type"`                       // Required. Must be set to "button" or "select".
	Value           string                        `json:"value,omitempty"`            // Optional.
	DataSource      string                        `json:"data_source,omitempty"`      // Optional.
	MinQueryLength  int                           `json:"min_query_length,omitempty"` // Optional. Default value is 1.
	Options         []AttachmentActionOption      `json:"options,omitempty"`          // Optional. Maximum of 100 options can be provided in each menu.
	SelectedOptions []AttachmentActionOption      `json:"selected_options,omitempty"` // Optional. The first element of this array will be set as the pre-selected option for this menu.
	OptionGroups    []AttachmentActionOptionGroup `json:"option_groups,omitempty"`    // Optional.
	Confirm         []ConfirmationField           `json:"confirm,omitempty"`          // Optional.
	URL             string                        `json:"url,omitempty"`              // Optional.
}

// AttachmentActionOption the individual option to appear in action menu.
type AttachmentActionOption struct {
	Text        string `json:"text"`                  // Required.
	Value       string `json:"value"`                 // Required.
	Description string `json:"description,omitempty"` // Optional. Up to 30 characters.
}

// AttachmentActionOptionGroup is a semi-hierarchal way to list available options to appear in action menu.
type AttachmentActionOptionGroup struct {
	Text    string                   `json:"text"`    // Required.
	Options []AttachmentActionOption `json:"options"` // Required.
}

// AttachmentActionCallback is sent from Slack when a user clicks a button in an interactive message (aka AttachmentAction)
type AttachmentActionCallback struct {
	Actions    []AttachmentAction `json:"actions"`
	CallbackID string             `json:"callback_id"`

	Name  string `json:"name"`
	Value string `json:"value"`

	ActionTs     string `json:"action_ts"`
	MessageTs    string `json:"message_ts"`
	AttachmentID string `json:"attachment_id"`
	Token        string `json:"token"`
	ResponseURL  string `json:"response_url"`
	TriggerID    string `json:"trigger_id"`
}

// ConfirmationField are used to ask users to confirm actions
type ConfirmationField struct {
	Title       string `json:"title,omitempty"`        // Optional.
	Text        string `json:"text"`                   // Required.
	OkText      string `json:"ok_text,omitempty"`      // Optional. Defaults to "Okay"
	DismissText string `json:"dismiss_text,omitempty"` // Optional. Defaults to "Cancel"
}

// Attachment contains all the information for an attachment
type Attachment struct {
	Color    string `json:"color,omitempty"`
	Fallback string `json:"fallback"`

	CallbackID string `json:"callback_id,omitempty"`
	ID         int    `json:"id,omitempty"`

	AuthorID      string `json:"author_id,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	AuthorLink    string `json:"author_link,omitempty"`
	AuthorIcon    string `json:"author_icon,omitempty"`

	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Text      string `json:"text"`

	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`

	Fields     []AttachmentField  `json:"fields,omitempty"`
	Actions    []AttachmentAction `json:"actions,omitempty"`
	MarkdownIn []string           `json:"mrkdwn_in,omitempty"`

	Footer     string `json:"footer,omitempty"`
	FooterIcon string `json:"footer_icon,omitempty"`

	Ts int64 `json:"ts,omitempty"`
}

type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconUrl     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia bool         `json:"unfurl_media,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

//func (attachment *Attachment) AddField(field Field) *Attachment {
//	attachment.Fields = append(attachment.Fields, &field)
//	return attachment
//}

func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("Incorrect token (redirection)")
}

func Send(webhookUrl string, proxy string, payload Payload) []error {
	request := gorequest.New().Proxy(proxy)
	resp, _, err := request.
		Post(webhookUrl).
		RedirectPolicy(redirectPolicyFunc).
		Send(payload).
		End()

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return []error{fmt.Errorf("Error sending msg. Status: %v", resp.Status)}
	}

	return nil
}

func PostMessage(webhookUrl string, payload Payload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(t))
	}

	return nil
}
