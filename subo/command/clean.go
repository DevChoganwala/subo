package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/context"
)

//CleanCmd  removes all of the target/.build folders for Runnables and deletes the .wasm files.
func CleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove build folders and .wasm files",
		Long:  "remove all of target/.build folders and deletes .wasm files",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				cwd = "$HOME"
			}

			bctx, err := context.CurrentBuildContext(cwd)
			if err != nil {
				return errors.Wrap(err, "failed to get CurrentBuildContext")
			}

			if len(bctx.Runnables) == 0 {
				return errors.New("🚫 no runnables found in current directory (no .runnable yaml files found)")
			}
			logStart(fmt.Sprintf("cleaning in %s", bctx.Cwd))

			//regex to find .wasm files
			fileRegEx, _ := regexp.Compile("^.+\\.wasm")

			for i := 0; i < len(bctx.Runnables); i++ {
				//delete target or .build folder
				dirs, _ := ioutil.ReadDir(".")
				for _, dir := range dirs {
					if dir.IsDir() {
						if dir.Name() == "target" || dir.Name() == ".build" {
							if rErr := os.RemoveAll(dir.Name()); rErr != nil {
								logInfo(rErr.Error())
							}
							logDone(fmt.Sprintf("removed %s", dir.Name()))
						}
					} else {
						if fileRegEx.MatchString(dir.Name()) {
							os.Remove(dir.Name())
						}
					}
				}
			}

			logDone("cleaned")
			return nil
		},
	}

	return cmd
}
