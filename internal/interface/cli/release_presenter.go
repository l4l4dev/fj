package cli

import (
	"fmt"
	"io"

	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
)

type releasePresenter struct{}

func (releasePresenter) Present(w io.Writer, releases []applicationrelease.Release) error {
	if _, err := fmt.Fprintln(w, "Releases:"); err != nil {
		return err
	}
	if len(releases) == 0 {
		_, err := fmt.Fprintln(w, "No releases found.")
		return err
	}
	for _, release := range releases {
		state := "published"
		if release.Draft {
			state = "draft"
		} else if release.Prerelease {
			state = "prerelease"
		}
		if _, err := fmt.Fprintf(w, "%s %s [%s]\n", release.TagName, release.Title, state); err != nil {
			return err
		}
	}
	return nil
}
