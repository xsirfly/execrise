package command

import (
	"exercise/database"
	"fmt"

	"github.com/docker/docker/api/types/mount"
)

type Context interface {
	GenCmd(chapter *database.Chapter) []string
	GenMounts(base string) []mount.Mount
	GetImage() string
	GetWorkDir() string
}

type BaseContext struct {
	Image      string
	WorkDir    string
	RunCommand string
	RunArgs    []string
}

func (ctx BaseContext) GetImage() string {
	return ctx.Image
}

func (ctx BaseContext) GetWorkDir() string {
	return ctx.WorkDir
}

type GradleContext struct {
	BaseContext
}

func (ctx GradleContext) GenCmd(chapter *database.Chapter) []string {
	var cmd []string
	cmd = append(cmd, ctx.RunCommand)
	for _, arg := range ctx.RunArgs {
		cmd = append(cmd, arg)
	}
	cmd = append(cmd, fmt.Sprintf(":%s:test", chapter.CodeLocation))
	return cmd
}

func (ctx GradleContext) GenMounts(base string) []mount.Mount {
	return []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: base,
			Target: "/usr/src/app",
		},
		{
			Type:   mount.TypeBind,
			Source: "/Users/xsir/.gradle",
			Target: "/home/gradle/.gradle",
		},
	}
}

var contextMap map[string]Context

func init() {
	contextMap = make(map[string]Context, 0)
	contextMap["gradle"] = GradleContext{
		BaseContext: BaseContext{
			Image:      "xsirfly/kylx:gradle",
			WorkDir:    "/usr/src/app",
			RunCommand: "gradle",
			RunArgs:    []string{"-i", "--no-daemon"},
		},
	}
}

func GetRunContextByBuild(projectKey string) (Context, bool) {
	c, ok := contextMap[projectKey]
	return c, ok
}
