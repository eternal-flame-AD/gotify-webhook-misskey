package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  "github.com/eternal-flame-ad/gotify-webhook-misskey",
		Version:     "0.1.1",
		Author:      "eternal-flame-ad <yume@yumechi.jp>",
		Website:     "https://github.com/eternal-flame-ad/gotify-webhook-misskey",
		Description: "Webhook Bridge for Misskey",
		License:     "Apache-2.0",
		Name:        "Gotify Misskey Webhook Bridge",
	}
}

// MisskeyHookPlugin is the gotify plugin instance.
type MisskeyHookPlugin struct {
	basePath   string
	enabled    bool
	config     *Config
	msgHandler plugin.MessageHandler
}

func (c *MisskeyHookPlugin) SetMessageHandler(handler plugin.MessageHandler) {
	c.msgHandler = handler
}

// Enable enables the plugin.
func (c *MisskeyHookPlugin) Enable() error {
	c.enabled = true
	return nil
}

// Disable disables the plugin.
func (c *MisskeyHookPlugin) Disable() error {
	c.enabled = false
	return nil
}

func (c *MisskeyHookPlugin) DefaultConfig() any {
	conf := CreateDefaultConfig()
	return &conf
}

func (c *MisskeyHookPlugin) ValidateAndSetConfig(input any) error {
	config := input.(*Config)
	if err := config.Validate(); err != nil {
		return err
	}
	c.config = config
	return nil
}

func escapeMarkdown(s string) string {
	return strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
	).Replace(s)
}

func truncateString(s string, length int) string {
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}

// RegisterWebhook implements plugin.Webhooker.
func (c *MisskeyHookPlugin) RegisterWebhook(basePath string, g *gin.RouterGroup) {
	c.basePath = basePath
	g.HEAD("/push/misskey/:slug", func(ctx *gin.Context) {
		ctx.SetAccepted("application/json")
		ctx.Status(200)
	})
	g.GET("/push/misskey/:slug", func(ctx *gin.Context) {
		ctx.JSON(405, gin.H{"error": "Method Not Allowed"})
	})
	g.POST("/push/misskey/:slug", func(ctx *gin.Context) {

		secret := ctx.GetHeader("X-Misskey-Hook-Secret")

		if secret == "" {
			ctx.JSON(400, gin.H{"error": "Missing secret"})
			return
		}

		src := c.config.GetSource(ctx.Param("slug"))

		if src == nil || src.Secret == DummySecret || src.Secret != secret {
			ctx.JSON(404, gin.H{"error": "Source not found or secret mismatch"})
			return
		}

		var payload WebhookPayload[UserPayload]

		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		if payload.Body.Note != nil {
			post := payload.Body.Note

			var content string
			if payload.Body.Reaction != nil {
				content = payload.Body.Reaction.Reaction
			} else {
				content = post.Text
			}

			title := fmt.Sprintf("[%s] [%s] %s: %s", payload.Type, src.Name, post.User.UserNameFull(), truncateString(content, 50))

			var url string
			if payload.Server != "" {
				url = strings.TrimRight(payload.Server, "/") + "/notes/" + post.ID
			}

			var message bytes.Buffer

			if payload.Body.Reaction != nil {
				message.WriteString(fmt.Sprintf("Reaction from %s (%s): %s\n\n", escapeMarkdown(payload.Body.Reaction.User.Name),
					escapeMarkdown(payload.Body.Reaction.User.Username), escapeMarkdown(payload.Body.Reaction.Reaction)))
			}

			message.WriteString(fmt.Sprintf("Post User: %s (%s)\n\n", escapeMarkdown(post.User.Name), escapeMarkdown(post.User.Username)))

			if post.Cw != nil {
				message.WriteString(fmt.Sprintf("CW: %s\n\n", escapeMarkdown(*post.Cw)))
			} else {
				message.WriteString(escapeMarkdown(post.Text))
			}

			if post.Reply != nil {
				message.WriteString("\n\n---\n\n")

				message.WriteString(fmt.Sprintf("Parent: %s\n\n", escapeMarkdown(post.Reply.User.Name)))

				if post.Reply.Cw != nil {
					message.WriteString(fmt.Sprintf("CW: %s\n\n", escapeMarkdown(*post.Reply.Cw)))
				} else {
					message.WriteString(escapeMarkdown(post.Reply.Text))
				}

				message.WriteString("\n\n---\n\n")
			}

			if post.Renote != nil {
				message.WriteString("\n\n---\n\n")

				message.WriteString(fmt.Sprintf("Renote of: %s\n\n", escapeMarkdown(post.Renote.User.Name)))

				if post.Renote.Cw != nil {
					message.WriteString(fmt.Sprintf("CW: %s\n\n", escapeMarkdown(*post.Renote.Cw)))
				} else {
					message.WriteString(escapeMarkdown(post.Renote.Text))
				}

				message.WriteString("\n\n---\n\n")
			}

			msg := plugin.Message{
				Title:    title,
				Message:  message.String(),
				Priority: src.Priority,
				Extras: map[string]interface{}{
					"misskey::payload": map[string]interface{}{
						"note": post,
					},
					"client::display": map[string]interface{}{
						"contentType": "text/markdown",
					},
				},
			}

			if url != "" {
				msg.Extras["client::notification"] = map[string]interface{}{
					"click": map[string]interface{}{
						"url": url,
					},
				}
			}

			if err := c.msgHandler.SendMessage(msg); err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to send message"})
				return
			}
		} else if payload.Body.User != nil {
			user := payload.Body.User

			title := fmt.Sprintf("[%s] [%s] %s", payload.Type, src.Name, user.UserNameFull())

			var url string

			if payload.Server != "" {
				url = strings.TrimRight(payload.Server, "/") + "/" + user.UserNameFull()
			}

			var message bytes.Buffer

			message.WriteString(fmt.Sprintf("User: %s (%s)\n\n", escapeMarkdown(user.Name), escapeMarkdown(user.Username)))

			message.WriteString(fmt.Sprintf("Followers: %d\n", user.FollowersCount))
			message.WriteString(fmt.Sprintf("Following: %d\n", user.FollowingCount))
			message.WriteString(fmt.Sprintf("Notes: %d\n", user.NotesCount))

			msg := plugin.Message{
				Title:    title,
				Message:  message.String(),
				Priority: 0,
				Extras: map[string]interface{}{
					"misskey::payload": map[string]interface{}{
						"user": user,
					},
					"client::display": map[string]interface{}{
						"contentType": "text/markdown",
					},
				},
			}

			if url != "" {
				msg.Extras["client::notification"] = map[string]interface{}{
					"click": map[string]interface{}{
						"url": url,
					},
				}
			}

			if err := c.msgHandler.SendMessage(msg); err != nil {
				ctx.JSON(500, gin.H{"error": "Failed to send message"})
				return
			}
		} else {
			ctx.JSON(400, gin.H{"error": "Invalid payload, unknown body"})
			return
		}
	})

	g.HEAD("/push/misskey/:slug/abuse", func(ctx *gin.Context) {
		ctx.SetAccepted("application/json")
		ctx.Status(200)
	})
	g.GET("/push/misskey/:slug/abuse", func(ctx *gin.Context) {
		ctx.JSON(405, gin.H{"error": "Method Not Allowed"})
	})
	g.POST("/push/misskey/:slug/abuse", func(ctx *gin.Context) {

		secret := ctx.GetHeader("X-Misskey-Hook-Secret")

		if secret == "" {
			ctx.JSON(400, gin.H{"error": "Missing secret"})
			return
		}

		src := c.config.GetSource(ctx.Param("slug"))

		if src == nil || src.Secret == DummySecret || src.Secret != secret {
			ctx.JSON(404, gin.H{"error": "Source not found or secret mismatch"})
			return
		}

		var payload WebhookPayload[AbuseReportWebhookPayloadBody]

		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		report := payload.Body

		title := fmt.Sprintf("[abuse] [%s] %s", src.Name, report.ID)

		var url string
		if payload.Server != "" {
			url = strings.TrimRight(payload.Server, "/") + "/admin/abuses"
		}

		var message bytes.Buffer

		message.WriteString(fmt.Sprintf("User: %s\n\n", escapeMarkdown(report.ID)))

		if report.Comment == nil {
			message.WriteString("<no comment>")
		} else if *report.Comment == "" {
			message.WriteString("<empty comment>")
		} else {
			message.WriteString(escapeMarkdown(*report.Comment))
		}

		msg := plugin.Message{
			Title:    title,
			Message:  message.String(),
			Priority: src.Priority,
			Extras: map[string]interface{}{
				"misskey::payload": map[string]interface{}{
					"abuse_report": report,
				},
				"client::display": map[string]interface{}{
					"contentType": "text/markdown",
				},
			},
		}

		if url != "" {
			msg.Extras["client::notification"] = map[string]interface{}{
				"click": map[string]interface{}{
					"url": url,
				},
			}
		}

		if err := c.msgHandler.SendMessage(msg); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to send message"})
			return
		}
	})
}

var displayTemplate = template.Must(
	template.New("display").Parse(
		strings.TrimSpace(
			`
# Misskey Webhook Plugin

Enabled: {{ .State.Enabled }}

## Sources:

{{ range .Sources }}

### **{{ .Name }}** ({{ .Slug }})

- Secret: {{ .Secret }}
- URL: [{{ .URL }}]({{ .URL }}) (Append /abuse to receive abuse reports)

{{ end }}
		`)))

func (c *MisskeyHookPlugin) GetDisplay(location *url.URL) string {
	type displayData struct {
		State   struct{ Enabled bool }
		Sources []struct {
			Name   string
			Slug   string
			Secret string
			URL    string
		}
	}

	loc := &url.URL{
		Path: c.basePath,
	}
	if location != nil {
		// If the server location can be determined, make the URL absolute
		loc.Scheme = location.Scheme
		loc.Host = location.Host
	}

	data := displayData{}

	data.State.Enabled = c.enabled

	for _, source := range c.config.Sources {
		thisLoc := loc.ResolveReference(&url.URL{
			Path: fmt.Sprintf("push/misskey/%s", source.Slug),
		})
		data.Sources = append(data.Sources, struct {
			Name   string
			Slug   string
			Secret string
			URL    string
		}{
			Name:   source.Name,
			Slug:   source.Slug,
			Secret: source.Secret,
			URL:    thisLoc.String(),
		})
	}

	var write bytes.Buffer

	err := displayTemplate.Execute(&write, data)

	if err != nil {
		return err.Error()
	}

	return write.String()
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	conf := CreateDefaultConfig()
	return &MisskeyHookPlugin{
		config:  &conf,
		enabled: false,
	}
}

func main() {
	panic("this should be built as go plugin")
}
