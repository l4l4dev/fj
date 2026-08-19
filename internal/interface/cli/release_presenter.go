package cli

import (
	"fmt"
	"io"
	"strings"

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

func (releasePresenter) PresentInspect(w io.Writer, detail applicationrelease.ReleaseDetail) error {
	state := "published"
	if detail.Draft {
		state = "draft"
	} else if detail.Prerelease {
		state = "prerelease"
	}
	notes := detail.Notes
	if notes == "" {
		notes = "-"
	}
	if _, err := fmt.Fprintf(w, "Release: %s\nTitle: %s\nState: %s\nNotes: %s\n", detail.TagName, detail.Title, state, notes); err != nil {
		return err
	}
	if len(detail.Assets) == 0 {
		_, err := fmt.Fprintln(w, "Assets: none")
		return err
	}
	if _, err := fmt.Fprintln(w, "Assets:"); err != nil {
		return err
	}
	for _, asset := range detail.Assets {
		if _, err := fmt.Fprintf(w, "- #%d %s (%d bytes)\n", asset.ID, asset.Name, asset.Size); err != nil {
			return err
		}
	}
	return nil
}

func (releasePresenter) PresentCreated(w io.Writer, detail applicationrelease.ReleaseDetail) error {
	_, err := fmt.Fprintf(w, "Release created as draft: %s\nTitle: %s\nPrerelease: %t\n", detail.TagName, detail.Title, detail.Prerelease)
	return err
}

func (releasePresenter) PresentUpdated(w io.Writer, detail applicationrelease.ReleaseDetail, fields []string) error {
	state := "published"
	if detail.Draft {
		state = "draft"
	} else if detail.Prerelease {
		state = "prerelease"
	}
	_, err := fmt.Fprintf(w, "Release updated: %s\nChanged fields: %s\nState: %s\n", detail.TagName, strings.Join(fields, ", "), state)
	return err
}

func (releasePresenter) PresentPublished(w io.Writer, detail applicationrelease.ReleaseDetail) error {
	_, err := fmt.Fprintf(w, "Release published: %s\nTitle: %s\n", detail.TagName, detail.Title)
	return err
}
