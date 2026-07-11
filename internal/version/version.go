package version

// Value is set to "dev" for local builds and may be overridden with -ldflags.
var Value = "dev"

func Current() string {
	if Value == "" {
		return "dev"
	}
	return Value
}
