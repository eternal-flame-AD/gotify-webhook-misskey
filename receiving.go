package main

import (
	"encoding/json"
	"fmt"
)

type WebhookUser struct {
	ID             string `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	UsernameLower  string `json:"usernameLower,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowersCount int    `json:"followersCount,omitempty"`
	FollowingCount int    `json:"followingCount,omitempty"`
	NotesCount     int    `json:"notesCount,omitempty"`
}

func (u *WebhookUser) UnmarshalJSON(b []byte) error {
	var id string
	if err := json.Unmarshal(b, &id); err == nil {
		u.ID = id
		u.Name = fmt.Sprintf("User %s", id)
		u.UsernameLower = id
		return nil
	}

	return json.Unmarshal(b, u)
}

type NoteRelatedWebhookPayload struct {
	ID   string  `json:"id,omitempty"`
	Text *string `json:"text,omitempty"`

	Mentions *[]string `json:"mentions,omitempty"`

	User   WebhookUser `json:"user,omitempty"`
	UserId string      `json:"userId,omitempty"`

	Renote   *string `json:"renote,omitempty"`
	RenoteId *string `json:"renoteId,omitempty"`

	Reply   *string `json:"reply,omitempty"`
	ReplyId *string `json:"replyId,omitempty"`

	Name *string `json:"name,omitempty"`
	Cw   *string `json:"cw,omitempty"`

	Visibility string `json:"visibility,omitempty"`

	RenoteCount  int  `json:"renoteCount,omitempty"`
	RepliesCount int  `json:"repliesCount,omitempty"`
	ClippedCount *int `json:"clippedCount,omitempty"`

	IsHidden *bool   `json:"isHidden,omitempty"`
	ThreadId *string `json:"threadId,omitempty"`
}

type AbuseReportWebhookPayload struct {
	ID string `json:"id,omitempty"`

	TargetUserId string      `json:"targetUserId,omitempty"`
	TargetUser   WebhookUser `json:"targetUser,omitempty"`

	ReporterId string      `json:"reporterId,omitempty"`
	Reporter   WebhookUser `json:"reporter,omitempty"`

	Comment *string `json:"comment,omitempty"`
}
