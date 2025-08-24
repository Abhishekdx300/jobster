package constants

const (
	TagRemote   = "remote"
	TagStartup  = "startup"
	TagPending  = "pending"
	TagFinished = "finished"
)

var allowedJobTags = map[string]struct{}{
	TagRemote:  {},
	TagStartup: {},
}
var allowedApplyTags = map[string]struct{}{
	TagPending:  {},
	TagFinished: {},
}

func IsValidJobTag(tag string) bool {
	_, ok := allowedJobTags[tag]
	return ok
}
func IsValidApplyTag(tag string) bool {
	_, ok := allowedApplyTags[tag]
	return ok
}
