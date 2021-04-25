package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

			for _, r := range bctx.Runnables {
				//delete target or .build folder
				dirs, _ := ioutil.ReadDir(r.Fullpath)

				for _, dir := range dirs {
					fullPath := filepath.Join(r.Fullpath, dir.Name())
					if dir.IsDir() {
						if dir.Name() == "target" || dir.Name() == ".build" {
							if rErr := os.RemoveAll(fullPath); rErr != nil {
								logInfo(errors.Wrap(rErr, "🚫 failed to RemoveAll").Error())
								continue
							}

							logDone(fmt.Sprintf("removed %s", dir.Name()))
						}
					} else {
						if strings.HasSuffix(dir.Name(), ".wasm") || strings.HasSuffix(dir.Name(), ".wasm.zip") {
							logInfo(dir.Name())
							if err := os.Remove(fullPath); err != nil {
								logInfo(errors.Wrap(err, "🚫 failed to Remove").Error())
								continue
							}

							logDone(fmt.Sprintf("removed %s", dir.Name()))
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
