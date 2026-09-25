package vikunjabot

import "embed"

//go:embed localizations/*.lang
var LocalizationPack embed.FS
