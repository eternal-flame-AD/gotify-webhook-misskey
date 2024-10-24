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
		Version:     "1.0.0",
		Author:      "eternal-flame-ad <yume@yumechi.jp>",
		Website:     "https://github.com/eternal-flame-ad/gotify-webhook-misskey",
		Description: "Webhook Bridge for Misskey",
		License:     "Apache-2.0",
		Name:        "eternal-flame-ad/gotify-webhook-misskey",
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

		var payload WebhookPayload[NoteRelatedWebhookPayloadBody]

		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		post := payload.Body

		title := fmt.Sprintf("Misskey %s from %s", payload.Type, post.User.Name)

		var url string
		if payload.Server != "" {
			url = strings.TrimRight(payload.Server, "/") + "/notes/" + post.ID
		}

		var message bytes.Buffer

		message.WriteString(fmt.Sprintf("User: %s (%s)\n\n", escapeMarkdown(post.User.Name), escapeMarkdown(post.User.Username)))

		if post.Cw != nil {
			message.WriteString(fmt.Sprintf("CW: %s\n\n", escapeMarkdown(*post.Cw)))
		} else if post.Text != nil {
			message.WriteString(escapeMarkdown(*post.Text))
		} else {
			message.WriteString("<missing content>")
		}

		if post.Reply != nil {
			message.WriteString("\n\n---\n\n")

			message.WriteString(fmt.Sprintf("Parent: %s\n\n", escapeMarkdown(post.Reply.User.Name)))

			if post.Reply.Cw != nil {
				message.WriteString(fmt.Sprintf("CW: %s\n\n", escapeMarkdown(*post.Reply.Cw)))
			} else if post.Reply.Text != nil {
				message.WriteString(escapeMarkdown(*post.Reply.Text))
			} else {
				message.WriteString("<missing content>")
			}

			message.WriteString("\n\n---\n\n")
		}

		if post.Renote != nil {
			message.WriteString("\n\n---\n\n")

			message.WriteString(fmt.Sprintf("Renote of: %s\n\n", escapeMarkdown(post.Renote.User.Name)))

			if post.Renote.Cw != nil {
				message.WriteString(fmt.Sprintf("CW: %s\n\n", escapeMarkdown(*post.Renote.Cw)))
			} else if post.Renote.Text != nil {
				message.WriteString(escapeMarkdown(*post.Renote.Text))
			} else {
				message.WriteString("<missing content>")
			}

			message.WriteString("\n\n---\n\n")
		}

		msg := plugin.Message{
			Title:    title,
			Message:  message.String(),
			Priority: 0,
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
	})

	g.HEAD("/push/misskey/:slug/follow", func(ctx *gin.Context) {
		ctx.SetAccepted("application/json")
		ctx.Status(200)
	})

	g.GET("/push/misskey/:slug/follow", func(ctx *gin.Context) {
		ctx.JSON(405, gin.H{"error": "Method Not Allowed"})
	})

	g.POST("/push/misskey/:slug/follow", func(ctx *gin.Context) {

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

		var payload WebhookPayload[WebhookUser]

		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		user := payload.Body

		title := fmt.Sprintf("Misskey Follow from %s", user.Name)

		var url string
		if payload.Server != "" {
			url = strings.TrimRight(payload.Server, "/") + "/users/" + user.ID
		}

		var message bytes.Buffer

		message.WriteString(fmt.Sprintf("User: %s (%s)\n\n", escapeMarkdown(user.Name), escapeMarkdown(user.Username)))

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

		title := fmt.Sprintf("Misskey Abuse Report from %s", report.ID)

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
			Priority: 0,
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
