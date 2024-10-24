package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testReplyPayload = `{"server":"https://mi.yumechi.jp","hookId":"9zr9b1z0cq2o01rf","userId":"9vaew1esmfme0001","eventId":"1bbc767c-6060-4982-9aa9-164321e89c4a","createdAt":1729798889614,"type":"reply","body":{"id":"dummy-reply-1","createdAt":"2024-10-24T19:41:29.613Z","deletedAt":null,"text":"This is a dummy note for testing purposes.","cw":null,"userId":"dummy-user-1","user":{"id":"dummy-user-1","name":"DummyUser1","username":"dummy1","host":null,"avatarUrl":null,"avatarBlurhash":null,"avatarDecorations":[],"isBot":false,"isCat":true,"emojis":[],"onlineStatus":"active","badgeRoles":[]},"replyId":"dummy-note-1","renoteId":null,"isHidden":false,"visibility":"public","mentions":[],"visibleUserIds":[],"fileIds":[],"files":[],"tags":[],"poll":null,"emojis":[],"channelId":null,"channel":null,"localOnly":true,"reactionAcceptance":"likeOnly","reactionEmojis":{},"reactions":{},"reactionCount":0,"renoteCount":10,"repliesCount":5,"reactionAndUserPairCache":[],"clippedCount":0,"reply":{"id":"dummy-note-1","createdAt":"2024-10-24T19:41:29.614Z","deletedAt":null,"text":"This is a dummy note for testing purposes.","cw":null,"userId":"dummy-user-1","user":{"id":"dummy-user-1","name":"DummyUser1","username":"dummy1","host":null,"avatarUrl":null,"avatarBlurhash":null,"avatarDecorations":[],"isBot":false,"isCat":true,"emojis":[],"onlineStatus":"active","badgeRoles":[]},"replyId":null,"renoteId":null,"isHidden":false,"visibility":"public","mentions":[],"visibleUserIds":[],"fileIds":[],"files":[],"tags":[],"poll":null,"emojis":[],"channelId":null,"channel":null,"localOnly":true,"reactionAcceptance":"likeOnly","reactionEmojis":{},"reactions":{},"reactionCount":0,"renoteCount":10,"repliesCount":5,"reactionAndUserPairCache":[]},"renote":null,"myReaction":null}}`

const testRenotePayload = `{"server":"https://mi.yumechi.jp","hookId":"9zrajfd0cq2o01rz","userId":"9vaew1esmfme0001","eventId":"49186f25-e235-43eb-a138-6cef98fe1528","createdAt":1729799956439,"type":"renote","body":{"id":"dummy-renote-1","createdAt":"2024-10-24T19:59:16.438Z","deletedAt":null,"text":null,"cw":null,"userId":"dummy-user-2","user":{"id":"dummy-user-2","name":"DummyUser2","username":"dummy2","host":null,"avatarUrl":null,"avatarBlurhash":null,"avatarDecorations":[],"isBot":false,"isCat":true,"emojis":[],"onlineStatus":"active","badgeRoles":[]},"replyId":null,"renoteId":"dummy-note-1","isHidden":false,"visibility":"public","mentions":[],"visibleUserIds":[],"fileIds":[],"files":[],"tags":[],"poll":null,"emojis":[],"channelId":null,"channel":null,"localOnly":true,"reactionAcceptance":"likeOnly","reactionEmojis":{},"reactions":{},"reactionCount":0,"renoteCount":10,"repliesCount":5,"reactionAndUserPairCache":[],"clippedCount":0,"reply":null,"renote":{"id":"dummy-note-1","createdAt":"2024-10-24T19:59:16.439Z","deletedAt":null,"text":"This is a dummy note for testing purposes.","cw":null,"userId":"dummy-user-1","user":{"id":"dummy-user-1","name":"DummyUser1","username":"dummy1","host":null,"avatarUrl":null,"avatarBlurhash":null,"avatarDecorations":[],"isBot":false,"isCat":true,"emojis":[],"onlineStatus":"active","badgeRoles":[]},"replyId":null,"renoteId":null,"isHidden":false,"visibility":"public","mentions":[],"visibleUserIds":[],"fileIds":[],"files":[],"tags":[],"poll":null,"emojis":[],"channelId":null,"channel":null,"localOnly":true,"reactionAcceptance":"likeOnly","reactionEmojis":{},"reactions":{},"reactionCount":0,"renoteCount":10,"repliesCount":5,"reactionAndUserPairCache":[],"clippedCount":0,"reply":null,"renote":null,"myReaction":null},"myReaction":null}}`

func TestReplyPayload(t *testing.T) {
	var payload WebhookPayload[NoteRelatedWebhookPayloadBody]

	err := json.Unmarshal([]byte(testReplyPayload), &payload)

	assert.Nil(t, err)

	assert.Equal(t, "https://mi.yumechi.jp", payload.Server)

	assert.Equal(t, "reply", payload.Type)

	assert.Equal(t, "9zr9b1z0cq2o01rf", payload.HookID)

	assert.Equal(t, "9vaew1esmfme0001", payload.UserID)

	assert.Equal(t, "1bbc767c-6060-4982-9aa9-164321e89c4a", payload.EventID)
}

func TestRenotePayload(t *testing.T) {
	var payload WebhookPayload[NoteRelatedWebhookPayloadBody]

	err := json.Unmarshal([]byte(testRenotePayload), &payload)

	assert.Nil(t, err)

	assert.Equal(t, "https://mi.yumechi.jp", payload.Server)

	assert.Equal(t, "renote", payload.Type)

	assert.Equal(t, "9zrajfd0cq2o01rz", payload.HookID)

	assert.Equal(t, "9vaew1esmfme0001", payload.UserID)

	assert.Equal(t, "49186f25-e235-43eb-a138-6cef98fe1528", payload.EventID)
}
