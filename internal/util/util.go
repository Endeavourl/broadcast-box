package util

func IsHiddenStream(streamName string) bool {
	// treat underscore-ending streams as hidden
	return streamName[len(streamName)-1] == '_'
}
