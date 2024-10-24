package main

import (
	"time"
)

type WebhookUser struct {
	ID             string `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	UsernameLower  string `json:"usernameLower,omitempty"`
	Name           string `json:"name,omitempty"`
	Host           string `json:"host,omitempty"`
	FollowersCount int    `json:"followersCount,omitempty"`
	FollowingCount int    `json:"followingCount,omitempty"`
	NotesCount     int    `json:"notesCount,omitempty"`
}

func (u *WebhookUser) UserNameFull() string {
	if u.Host != "" {
		return "@" + u.Username + "@" + u.Host
	}
	return "@" + u.Username
}

type WebhookPayload[T any] struct {
	Server  string `json:"server,omitempty"`
	Type    string `json:"type,omitempty"`
	HookID  string `json:"hookId,omitempty"`
	UserID  string `json:"userId,omitempty"`
	EventID string `json:"eventId,omitempty"`

	CreatedAt uint64 `json:"createdAt,omitempty"`

	Body T `json:"body,omitempty"`
}

type UserPayload struct {
	User *WebhookUser                   `json:"user,omitempty"`
	Note *NoteRelatedWebhookPayloadBody `json:"note,omitempty"`
}

func (p *WebhookPayload[T]) CreatedAtUnix() int64 {
	return int64(p.CreatedAt / 1000)
}

func (p *WebhookPayload[T]) CreatedAtDate() time.Time {
	return time.UnixMilli(int64(p.CreatedAt))
}

type NoteRelatedWebhookPayloadBody struct {
	ID   string `json:"id,omitempty"`
	Text string `json:"text,omitempty"`

	Mentions *[]string `json:"mentions,omitempty"`

	User   WebhookUser `json:"user,omitempty"`
	UserId string      `json:"userId,omitempty"`

	Renote   *NoteRelatedWebhookPayloadBody `json:"renote,omitempty"`
	RenoteId *string                        `json:"renoteId,omitempty"`

	Reply   *NoteRelatedWebhookPayloadBody `json:"reply,omitempty"`
	ReplyId *string                        `json:"replyId,omitempty"`

	Name *string `json:"name,omitempty"`
	Cw   *string `json:"cw,omitempty"`

	Visibility string `json:"visibility,omitempty"`

	RenoteCount  int  `json:"renoteCount,omitempty"`
	RepliesCount int  `json:"repliesCount,omitempty"`
	ClippedCount *int `json:"clippedCount,omitempty"`

	IsHidden *bool   `json:"isHidden,omitempty"`
	ThreadId *string `json:"threadId,omitempty"`
}

type AbuseReportWebhookPayloadBody struct {
	ID string `json:"id,omitempty"`

	TargetUserId string      `json:"targetUserId,omitempty"`
	TargetUser   WebhookUser `json:"targetUser,omitempty"`

	ReporterId string      `json:"reporterId,omitempty"`
	Reporter   WebhookUser `json:"reporter,omitempty"`

	Comment *string `json:"comment,omitempty"`
}
