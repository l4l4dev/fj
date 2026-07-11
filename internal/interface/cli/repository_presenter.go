package cli

import (
	"fmt"
	"io"
	"strings"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
)

type repositoryPresenter struct{}

func (repositoryPresenter) PresentRepository(w io.Writer, r applicationrepository.Repository) error {
	d := r.Description
	if d == "" {
		d = "-"
	}
	b := r.DefaultBranch
	if b == "" {
		b = "-"
	}
	_, err := fmt.Fprintf(w, "Repository: %s/%s\nDescription: %s\nPrivate: %t\nArchived: %t\nDefault branch: %s\n", r.Owner, r.Name, d, r.Private, r.Archived, b)
	return err
}
func (repositoryPresenter) PresentUpdated(w io.Writer, r applicationrepository.Repository, fields []string) error {
	d := r.Description
	if d == "" {
		d = "-"
	}
	b := r.DefaultBranch
	if b == "" {
		b = "-"
	}
	_, err := fmt.Fprintf(w, "Repository: %s/%s\nChanged fields: %s\nDescription: %s\nPrivate: %t\nArchived: %t\nDefault branch: %s\n", r.Owner, r.Name, strings.Join(fields, ", "), d, r.Private, r.Archived, b)
	return err
}
func (repositoryPresenter) PresentArchive(w io.Writer, r applicationrepository.Repository) error {
	_, err := fmt.Fprintf(w, "Repository: %s/%s\nArchived: %t\n", r.Owner, r.Name, r.Archived)
	return err
}
func (repositoryPresenter) PresentList(w io.Writer, rs []applicationrepository.Repository) error {
	if len(rs) == 0 {
		_, err := fmt.Fprintln(w, "No repositories found.")
		return err
	}
	for _, r := range rs {
		if _, err := fmt.Fprintf(w, "%s/%s\n", r.Owner, r.Name); err != nil {
			return err
		}
	}
	return nil
}
func (repositoryPresenter) PresentAccess(w io.Writer, a applicationrepository.RepositoryAccess) error {
	if _, err := fmt.Fprintf(w, "Repository: %s/%s\nCollaborators:\n", a.Owner, a.Name); err != nil {
		return err
	}
	if len(a.Collaborators) == 0 {
		_, err := fmt.Fprintln(w, "No collaborators found.")
		return err
	}
	for _, c := range a.Collaborators {
		if _, err := fmt.Fprintf(w, "- %s: %s\n", c.Username, c.Permission); err != nil {
			return err
		}
	}
	return nil
}
