package util

func IsHiddenStream(streamName string) bool {
	if streamName[len(streamName)-1] == '_' {
		// treat underscore-ending streams as hidden
		return true
	}
	return false
}
