package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/util/gtag"

	"hotgo/internal/library/hggen/internal/utility/allyes"
	"hotgo/internal/library/hggen/internal/utility/mlog"
)

var (
	Init = cInit{}
)

type cInit struct {
	g.Meta `name:"init" brief:"{cInitBrief}" eg:"{cInitEg}"`
}

const (
	cInitRepoPrefix = `github.com/gogf/`
	cInitMonoRepo   = `template-mono`
	cInitSingleRepo = `template-single`
	cInitBrief      = `create and initialize an empty GoFrame project`
	cInitEg         = `
gf init my-project
gf init my-mono-repo -m
`
	cInitNameBrief = `
name for the project. It will create a folder with NAME in current directory.
The NAME will also be the module name for the project.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cInitBrief`:     cInitBrief,
		`cInitEg`:        cInitEg,
		`cInitNameBrief`: cInitNameBrief,
	})
}

type cInitInput struct {
	g.Meta `name:"init"`
	Name   string `name:"NAME" arg:"true" v:"required" brief:"{cInitNameBrief}"`
	Mono   bool   `name:"mono"   short:"m" brief:"initialize a mono-repo instead a single-repo" orphan:"true"`
	Update bool   `name:"update" short:"u" brief:"update to the latest goframe version" orphan:"true"`
}

type cInitOutput struct{}

func (c cInit) Index(ctx context.Context, in cInitInput) (out *cInitOutput, err error) {
	if !gfile.IsEmpty(in.Name) && !allyes.Check() {
		s := gcmd.Scanf(`the folder "%s" is not empty, files might be overwrote, continue? [y/n]: `, in.Name)
		if strings.EqualFold(s, "n") {
			return
		}
	}
	mlog.Print("initializing...")

	// Create project folder and files.
	var (
		templateRepoName string
	)
	if in.Mono {
		templateRepoName = cInitMonoRepo
	} else {
		templateRepoName = cInitSingleRepo
	}
	err = gres.Export(templateRepoName, in.Name, gres.ExportOption{
		RemovePrefix: templateRepoName,
	})
	if err != nil {
		return
	}

	// Replace template name to project name.
	err = gfile.ReplaceDir(
		cInitRepoPrefix+templateRepoName,
		gfile.Basename(gfile.RealPath(in.Name)),
		in.Name,
		"*",
		true,
	)
	if err != nil {
		return
	}

	// Update the GoFrame version.
	if in.Update {
		mlog.Print("update goframe...")
		// go get -u github.com/gogf/gf/v2@latest
		updateCommand := `go get -u github.com/gogf/gf/v2@latest`
		if in.Name != "." {
			updateCommand = fmt.Sprintf(`cd %s && %s`, in.Name, updateCommand)
		}
		if err = gproc.ShellRun(ctx, updateCommand); err != nil {
			mlog.Fatal(err)
		}
		// go mod tidy
		gomModTidyCommand := `go mod tidy`
		if in.Name != "." {
			gomModTidyCommand = fmt.Sprintf(`cd %s && %s`, in.Name, gomModTidyCommand)
		}
		if err = gproc.ShellRun(ctx, gomModTidyCommand); err != nil {
			mlog.Fatal(err)
		}
	}

	mlog.Print("initialization done! ")
	if !in.Mono {
		enjoyCommand := `gf run main.go`
		if in.Name != "." {
			enjoyCommand = fmt.Sprintf(`cd %s && %s`, in.Name, enjoyCommand)
		}
		mlog.Printf(`you can now run "%s" to start your journey, enjoy!`, enjoyCommand)
	}
	return
}
