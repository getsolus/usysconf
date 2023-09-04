package cli

import (
	"log/slog"
	"os"

	"gitlab.com/slxh/go/powerline"
)

// Level contains the application log level.
var Level slog.LevelVar

var colors = map[slog.Level]powerline.ColorScheme{
	slog.LevelDebug: {
		Time:    powerline.NewColor(99, powerline.ColorBlack),
		Level:   powerline.NewColor(powerline.ColorBlack, 99),
		Message: powerline.NewColor(99, powerline.ColorDefault),
		Attrs:   powerline.NewColor(powerline.DefaultAttrColor, powerline.ColorDefault),
	},
	slog.LevelInfo: {
		Time:    powerline.NewColor(45, powerline.ColorBlack),
		Level:   powerline.NewColor(powerline.ColorBlack, 45),
		Message: powerline.NewColor(45, powerline.ColorDefault),
		Attrs:   powerline.NewColor(powerline.DefaultAttrColor, powerline.ColorDefault),
	},
	slog.LevelWarn: {
		Time:    powerline.NewColor(220, powerline.ColorBlack),
		Level:   powerline.NewColor(powerline.ColorBlack, 220),
		Message: powerline.NewColor(220, powerline.ColorDefault),
		Attrs:   powerline.NewColor(powerline.DefaultAttrColor, powerline.ColorDefault),
	},
	slog.LevelError: {
		Time:    powerline.NewColor(208, powerline.ColorBlack),
		Level:   powerline.NewColor(powerline.ColorBlack, 208),
		Message: powerline.NewColor(208, powerline.ColorDefault),
		Attrs:   powerline.NewColor(powerline.DefaultAttrColor, powerline.ColorDefault),
	},
}

func ConfigureLogger() {
	slog.SetDefault(slog.New(powerline.NewHandler(os.Stderr, &powerline.HandlerOptions{
		Level:  &Level,
		Colors: colors,
	})))
}
