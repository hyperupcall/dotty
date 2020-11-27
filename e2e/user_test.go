package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/eankeen/dotty/actions"
	"github.com/eankeen/dotty/config"
	"github.com/eankeen/dotty/internal/util"
)

func TestUserApply(t *testing.T) {
	// tests; path relative to destDir
	userTests := []struct {
		dir string
		fn  func(*testing.T, string, string)
	}{
		// basic; single file
		{
			dir: "user-1-file",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureSymlink(t, srcDir, destDir, "bar")
			},
		},
		// deeply nested subdirectories (file)
		{
			dir: "user-2-file-subdirs",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureSymlink(t, srcDir, destDir, "one")
				ensureDir(t, destDir, "subdir-1")
				ensureSymlink(t, srcDir, destDir, "subdir-1/two")
				ensureDir(t, destDir, "subdir-2/subdir-2-1")
				ensureSymlink(t, srcDir, destDir, "subdir-2/subdir-2-1/three")
			},
		},
		// test matching full path to root
		{
			dir: "user-3-file-root",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureDir(t, destDir, "afolder")
				ensureSymlink(t, srcDir, destDir, "afolder/afile")
			},
		},
		// basic; single folder
		{
			dir: "user-4-folder",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureDir(t, destDir, "afolder")
				ensureSymlink(t, srcDir, destDir, "afolder")
			},
		},
		// deeply nested subdirs (folders)
		{
			dir: "user-5-folder-subdirs",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureDir(t, destDir, "sub-1")
				ensureDir(t, destDir, "sub-1/sub-2")
				ensureDir(t, destDir, "sub-1/sub-2/sub-3")
				ensureSymlink(t, srcDir, destDir, "sub-1/sub-2/sub-3/afolder")
			},
		},
		// test matching full path to root (folder)
		{
			dir: "user-6-folder-root",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureDir(t, destDir, "afolder")
				ensureSymlink(t, srcDir, destDir, "afolder/asubfolder")
			},
		},
		// test defaults of dotty.toml
		{
			dir: "user-7-defaults",
			fn: func(t *testing.T, srcDir, destDir string) {
				defaultDestDir := filepath.Join(testDir(), "test-user", "user-7-defaults", "dest-HOME")

				ensureDir(t, defaultDestDir, "")
				ensureSymlink(t, srcDir, defaultDestDir, "bar")
			},
		},
		// test adding permission to file
		{
			dir: "user-8-permission",
			fn: func(t *testing.T, srcDir, destDir string) {
				ensureDir(t, destDir, "")
				ensureSymlink(t, srcDir, destDir, "bar")
				ensureMode(t, destDir, "bar", "-rw-------")
			},
		},
	}

	for _, userTest := range userTests {
		fmt.Printf("--- ON: '%s' ---\n", userTest.dir)

		dotfilesDir := filepath.Join(testDir(), "test-user", userTest.dir)
		dottyCfg := config.DottyCfg(dotfilesDir)

		srcDir := config.Src(dotfilesDir, dottyCfg, "user")
		destDir := config.Dest(dotfilesDir, dottyCfg, "user")

		homeDir, err := os.UserHomeDir()
		util.HandleError(err)

		if destDir == homeDir {
			fmt.Println("FAIL NOW")
			t.FailNow()
		}

		err = os.RemoveAll(destDir)
		util.HandleFsError(err)

		actions.Apply(dotfilesDir, srcDir, destDir)
		userTest.fn(t, srcDir, destDir)
	}

}
