package configurator

// default const values for application.
var defaults = map[ConfigKey]string{ //nolint: exhaustive // not all keys have defaults
	DatabaseFilenameKey: "wherehouse.db",
}
