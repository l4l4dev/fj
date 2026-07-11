package cli

import (
	"fmt"
	"io"
	"strings"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type issuePresenter struct{}

func (issuePresenter) PresentUpdated(w io.Writer, detail applicationissue.IssueDetail, fields []string) error {
	body := detail.Body
	if body == "" {
		body = "-"
	}
	_, err := fmt.Fprintf(w, "Issue: #%d\nChanged fields: %s\nTitle: %s\nState: %s\nBody: %s\n", detail.Number, strings.Join(fields, ", "), detail.Title, detail.State, body)
	return err
}

func (issuePresenter) PresentInspect(w io.Writer, detail applicationissue.IssueDetail) error {
	body := detail.Body
	if body == "" {
		body = "-"
	}
	_, err := fmt.Fprintf(w, "Issue: #%d\nTitle: %s\nState: %s\nBody: %s\n", detail.Number, detail.Title, detail.State, body)
	return err
}

func (issuePresenter) PresentList(w io.Writer, page applicationissue.Page) error {
	if _, err := fmt.Fprintln(w, "Issues:"); err != nil {
		return err
	}
	if len(page.Issues) == 0 {
		if _, err := fmt.Fprintln(w, "No issues found."); err != nil {
			return err
		}
	} else {
		for _, item := range page.Issues {
			if _, err := fmt.Fprintf(w, "- #%d %s [%s]\n", item.Number, item.Title, item.State); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintf(w, "\nPage: %d\nLimit: %d\nMore pages: %t\n", page.Page, page.Limit, page.MorePages)
	return err
}
