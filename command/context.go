package command

type Context struct {
	RunCommand     string
	RunArgs        []string
}

var contextMap map[string]*Context

func Init(testTarget string) {
	contextMap = make(map[string]*Context)

	contextMap["gradleJava"] = &Context{
		RunCommand: "gradle",
		RunArgs: []string{
			testTarget + ":test",
		},
	}
}

func GetRunContextByProjectKey(projectKey string) (*Context, bool) {
	c, ok := contextMap[projectKey]
	return c, ok
}
