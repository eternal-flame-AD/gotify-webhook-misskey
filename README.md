# gotify-website-misskey

Gotify Misskey Webhook Bridge

![Screenshot](assets/screenshot.png)

## Usage

1. Enable the plugin, set the webhook slugs:

```yaml
sources:
- slug: yumechi-no-kuni
  name: mi.yumechi.jp
  secret: xxxx
- slug: misskey-io
  name: misskey.io
  secret: xxxx
```

2. Copy the Webhook URL and secret to Misskey (Settings -> Webhook)

  a. Set the base URL to receive any of 'replied to', 'renoted', 'mentioned', 'followed' (reactions seem to have bugs upstream and nothing is received)

3. Done!

## License

Apache 2.0
