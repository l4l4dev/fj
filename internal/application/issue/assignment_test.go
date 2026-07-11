package issue

import (
	"context"
	"testing"
)

type assignmentStub struct {
	assign   AssignRequest
	unassign UnassignRequest
}

func (s *assignmentStub) Assign(_ context.Context, r AssignRequest) (Assignment, error) {
	s.assign = r
	return Assignment{Username: r.Username}, nil
}
func (s *assignmentStub) Unassign(_ context.Context, r UnassignRequest) error {
	s.unassign = r
	return nil
}

func TestAssignmentUseCasesDelegateAndValidate(t *testing.T) {
	s := &assignmentStub{}
	if _, err := NewAssignUseCase(s).Execute(context.Background(), AssignRequest{Owner: "alice", Name: "project", Number: 2, Username: "bob"}); err != nil || s.assign.Username != "bob" {
		t.Fatalf("assign failed: %+v %v", s.assign, err)
	}
	if err := NewUnassignUseCase(s).Execute(context.Background(), UnassignRequest{Owner: "alice", Name: "project", Number: 2}); err != nil {
		t.Fatal(err)
	}
	for _, r := range []AssignRequest{{Owner: "", Name: "project", Number: 1, Username: "bob"}, {Owner: "alice", Name: "project", Number: 0, Username: "bob"}, {Owner: "alice", Name: "project", Number: 1, Username: "none"}} {
		if _, err := NewAssignUseCase(s).Execute(context.Background(), r); err == nil {
			t.Fatalf("expected validation: %+v", r)
		}
	}
}
