# protoc-gen-go-protojson

Generate protobuf messages which implements [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) and [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler) interfaces.

Various reasons which makes this desirable:
- With [Go implement Protocol Buffers moving to opaque API](https://go.dev/blog/protobuf-opaque), `encoding/json` no longer _works out of the box_.
- The Protobuf maintainers has long stated that reflection output of `encoding/json` is not [canonical Protobuf JSON encoding](https://go.dev/blog/protobuf-opaque#reflection).
- Having Protobuf message serializable via `encoding/json` makes integrating `JSONB` fields with [sqlc](https://sqlc.dev) much smoother and transparent.

> [!WARNING]
> Consider this alpha software. Use at your own risk.

> [!TIP]
> There is a more popular project [mfridman/protoc-gen-go-json](https://github.com/mfridman/protoc-gen-go-json), the maintained fork of [mitchellh/protoc-gen-go-json](https://github.com/mitchellh/protoc-gen-go-json), which has the same goal.
> I just wanted the fun route.

## Acknowledgements

I managed to finish the first _usable_ version of this project in an evening thanks to a tutorial by @ericchiang: [Protobuf generators for fun and profit](https://ericchiang.github.io/post/protoc-plugins/).
