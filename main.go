package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	syspath "path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type renamer struct{ err error }

var (
	rootCmd = &cobra.Command{
		Use:   "renamer",
		Short: "renamer",
		Long:  "Batch naming of files",
	}
)

func main() {
	// global flags
	rootCmd.PersistentFlags().StringP("dir", "d", "", "directory path (required)")
	rootCmd.MarkPersistentFlagRequired("dir")
	rootCmd.PersistentFlags().StringP("old", "o", "", "old string")
	rootCmd.PersistentFlags().StringP("new", "n", "", "new string")

	ins := renamer{}

	// add subcommand
	rootCmd.AddCommand(&cobra.Command{
		Use:   "replace",
		Short: "replace -d file-directory --old=find-string --new=new-string",
		Run:   ins.replace,
	})
	rootCmd.AddCommand(&cobra.Command{Use: "append", Short: "append -d . --new=xxx", Run: ins.append})
	rootCmd.AddCommand(&cobra.Command{Use: "forward", Short: "forward  -d . --new=xxx", Run: ins.forward})

	rootCmd.Execute()
}

func (r *renamer) sflags(cmd *cobra.Command, name string) string {
	if r.err != nil {
		log.Fatalf("get flag err: %v\n", r.err)
	}
	var v string
	v, r.err = cmd.Flags().GetString(name)
	return v
}

func (r *renamer) replace(cmd *cobra.Command, args []string) {
	old := r.sflags(cmd, "old")
	new := r.sflags(cmd, "new")

	replacer := strings.NewReplacer(old, new)
	walkDir(r.sflags(cmd, "dir"), func(path string, d fs.DirEntry) error {
		newpath := replacer.Replace(path)
		if newpath == path {
			return nil
		}
		return os.Rename(path, newpath)
	})
}

func (r *renamer) append(cmd *cobra.Command, args []string) {
	new := r.sflags(cmd, "new")
	if new == "" {
		r.error("--new flag required")
		return
	}

	walkDir(r.sflags(cmd, "dir"), func(path string, d fs.DirEntry) error {
		ext := syspath.Ext(d.Name())
		newPath := strings.Replace(path, ext, new+ext, -1)
		return os.Rename(path, newPath)
	})
}

func (r *renamer) forward(cmd *cobra.Command, args []string) {
	new := r.sflags(cmd, "new")
	if new == "" {
		r.error("--new flag required")
		return
	}

	err := walkDir(r.sflags(cmd, "dir"), func(path string, d fs.DirEntry) error {
		newPath := strings.Replace(path, d.Name(), new+d.Name(), -1)
		return os.Rename(path, newPath)
	})

	r.msgIfErr(err)
}

func (r *renamer) error(msg string) {
	r.msgIfErr(errors.New(msg))
}

func (r *renamer) msgIfErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func walkDir(root string, fn func(path string, d fs.DirEntry) error) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return fn(path, d)
	})

	if os.IsNotExist(err) {
		return fmt.Errorf("directory[%s] does not exist: %s", root, err.Error())
	}

	return err
}
