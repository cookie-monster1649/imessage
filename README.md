# mautrix-imessage
A Matrix-iMessage puppeting bridge. The bridge can run on a Mac to bridge
iMessage or an Android phone to bridge SMS. All features are available when
using a Mac with SIP disabled, while a normal Mac can be used for basic
bridging. A [websocket proxy](https://github.com/mautrix/wsproxy)
is required to receive appservice events from the homeserver.

## About this fork

This is a fork of [mautrix/imessage](https://github.com/mautrix/imessage) carrying
two fixes to the **BlueBubbles** connector. Everything else is unchanged, and each
fix also lives on its own single-commit branch so it can be submitted upstream
unchanged.

**Outgoing messages go to the service the recipient actually supports**
([#2](https://github.com/cookie-monster1649/imessage/pull/2), branch
`bluebubbles-sms-service-routing`)

With `bridge.disable_sms_portals` enabled, portal GUIDs are rewritten from
`SMS;-;<address>` to `iMessage;-;<address>` so both services share one Matrix room.
That is correct for receiving, but on send it asks Messages to deliver an iMessage
to addresses that cannot receive one — so every message to a non-iMessage contact
failed with the opaque `could not send message`. This asks BlueBubbles what the
address actually supports and routes to SMS only when it genuinely cannot receive
iMessages. Answers are cached per address for an hour, email addresses are skipped
(SMS cannot deliver to one), and it fails open.

**Backfill queries are paginated**
([#1](https://github.com/cookie-monster1649/imessage/pull/1), branch
`bluebubbles-backfill-pagination`)

`GetMessagesWithLimit` and `GetMessagesBeforeWithLimit` called `queryChatMessages`
with `paginate=false`. BlueBubbles caps a page at 1000 messages, so every chat
backfilled at most ~1000 messages regardless of `backfill.initial_limit`. These now
match the two callers that already paginated correctly.

### Note on Docker images

Upstream's `Dockerfile.ci` pins `alpine:3.19`, which is fine. If you build your own
image on a newer base, be aware that **FFmpeg 8.1 and later refuse to mux CAF/Opus**
([`3fc7e39eb8`](https://github.com/FFmpeg/FFmpeg/commit/3fc7e39eb8)), which is the
format outgoing voice messages are converted to. Sends fail with
`ffmpeg error: exit status 176` (`AVERROR_PATCHWELCOME`). Use a base image with
FFmpeg 8.0 or older.

## Documentation
All setup and usage instructions are located on
[docs.mau.fi](https://docs.mau.fi/bridges/go/imessage/index.html):

* Bridge setup:
  [macOS](https://docs.mau.fi/bridges/go/imessage/mac/setup.html),
  [macOS (without SIP)](https://docs.mau.fi/bridges/go/imessage/mac-nosip/setup.html),
  [Android SMS](https://docs.mau.fi/bridges/go/imessage/android/setup.html)

### Features & Roadmap
[ROADMAP.md](https://github.com/mautrix/imessage/blob/master/ROADMAP.md)
contains a general overview of what is supported by the bridge.

## Discussion
Matrix room: [#imessage:maunium.net](https://matrix.to/#/#imessage:maunium.net)
