# gotify-website-misskey

Gotify Misskey Webhook Bridge (Alpha stage)

Receive Misskey Webhooks and forward them to Gotify. Useful for when you do not want to use the PWA or have multiple accounts on different instances.

![Screenshot](assets/screenshot.png)

## Usage

1. Enable the plugin, set the webhook slugs:

```yaml
sources:
- slug: yumechi-no-kuni
  name: mi.yumechi.jp
  secret: xxxx
  priority: 7
- slug: misskey-io
  name: misskey.io
  secret: xxxx
  priority: 5
```

1. Copy the Webhook URL and secret to Misskey (Settings -> Webhook) Set the WebHook URL to receive any of 'replied to', 'renoted', 'mentioned', 'followed' (reactions have bugs upstream and nothing is received, a server-side patch is required for this to work). See https://forge.yumechi.jp/yume/yumechi-no-kuni/commit/18d5587e5c559d45eb77fba1681239d31debfbb4#diff-b5c2505e496b6b3f76c0b67ddd74c9380b9725a5

2. Done! The on-click URL for notifications will also be populated when there is an User or Note associated with the notification.

## Notes

- The mockup webhook call payload before 2024.11.0 is incorrect, the actual payload is different. This is an upstream issue. Please test with a real event for now for the message to be generated correctly.

## License

Apache 2.0
