package cli

import (
	"fmt"
	"io"
	"strings"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type issuePresenter struct{}

func (issuePresenter) PresentLabelAdded(w io.Writer, number int, label applicationissue.Label) error {
	_, err := fmt.Fprintf(w, "Issue: #%d\nLabel added: %s\n", number, label.Name)
	return err
}

func (issuePresenter) PresentLabelRemoved(w io.Writer, number int, label applicationissue.Label) error {
	_, err := fmt.Fprintf(w, "Issue: #%d\nLabel removed: %s\n", number, label.Name)
	return err
}

func (issuePresenter) PresentComments(w io.Writer, comments []applicationissue.Comment) error {
	if _, err := fmt.Fprintln(w, "Comments:"); err != nil {
		return err
	}
	if len(comments) == 0 {
		_, err := fmt.Fprintln(w, "No comments found.")
		return err
	}
	for _, comment := range comments {
		if _, err := fmt.Fprintf(w, "- #%d %s\n", comment.ID, comment.Body); err != nil {
			return err
		}
	}
	return nil
}

func (issuePresenter) PresentComment(w io.Writer, comment applicationissue.Comment) error {
	_, err := fmt.Fprintf(w, "Comment:\n#%d %s\n", comment.ID, comment.Body)
	return err
}

func (issuePresenter) PresentState(w io.Writer, detail applicationissue.IssueDetail) error {
	_, err := fmt.Fprintf(w, "Issue: #%d\nState: %s\n", detail.Number, detail.State)
	return err
}

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
