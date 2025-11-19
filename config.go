package fault

var maxDepthStackTrace = 32

func SetMaxDepthStackTrace(depth int) {
	maxDepthStackTrace = depth
}

func GetMaxDepthStackTrace() int {
	return maxDepthStackTrace
}
