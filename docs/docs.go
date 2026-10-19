package docs

import "embed"

// Content embeds the component documentation tree so the server can serve it
// under /showcase/docs. Demo pages link to docs/<name>.md relative to
// /showcase/<name>, so that href must resolve to a live URL rather than a
// repository path.
//
//go:embed *.md
var Content embed.FS
