package release

import (
	"context"
	"testing"
)

type assetUploaderStub struct{ request UploadAssetRequest }

func (s *assetUploaderStub) UploadAsset(_ context.Context, request UploadAssetRequest) (Asset, error) {
	s.request = request
	return Asset{ID: 11, Name: request.AssetName, Size: int64(len(request.Content))}, nil
}

type assetDeleterStub struct {
	called  bool
	request DeleteAssetRequest
}

func (s *assetDeleterStub) DeleteAsset(_ context.Context, request DeleteAssetRequest) error {
	s.called = true
	s.request = request
	return nil
}

func TestUploadAssetUseCaseForwardsRequest(t *testing.T) {
	stub := &assetUploaderStub{}
	result, err := NewUploadAssetUseCase(stub).Execute(context.Background(), UploadAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: []byte("payload")})
	if err != nil || result.ID != 11 || result.Name != "fj.tar.gz" || result.Size != 7 {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if stub.request.Tag != "v1.0.0" || string(stub.request.Content) != "payload" {
		t.Fatalf("unexpected request: %+v", stub.request)
	}
}

func TestUploadAssetUseCaseRejectsInvalidInput(t *testing.T) {
	content := []byte("payload")
	for _, request := range []UploadAssetRequest{
		{Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: content},
		{Owner: "alice", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: content},
		{Owner: "alice", Name: "project", AssetName: "fj.tar.gz", Content: content},
		{Owner: "alice", Name: "project", Tag: "  ", AssetName: "fj.tar.gz", Content: content},
		{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "  ", Content: content},
		{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"},
		{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: []byte{}},
	} {
		if _, err := NewUploadAssetUseCase(&assetUploaderStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestUploadAssetUseCaseRejectsNilUploader(t *testing.T) {
	if _, err := NewUploadAssetUseCase(nil).Execute(context.Background(), UploadAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: []byte("payload")}); err == nil {
		t.Fatal("expected internal error for nil uploader")
	}
}

func TestDeleteAssetUseCaseForwardsRequest(t *testing.T) {
	stub := &assetDeleterStub{}
	if err := NewDeleteAssetUseCase(stub).Execute(context.Background(), DeleteAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"}); err != nil {
		t.Fatal(err)
	}
	if !stub.called || stub.request.Tag != "v1.0.0" || stub.request.AssetName != "fj.tar.gz" {
		t.Fatalf("unexpected request: %+v", stub.request)
	}
}

func TestDeleteAssetUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []DeleteAssetRequest{
		{Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"},
		{Owner: "alice", Tag: "v1.0.0", AssetName: "fj.tar.gz"},
		{Owner: "alice", Name: "project", AssetName: "fj.tar.gz"},
		{Owner: "alice", Name: "project", Tag: "  ", AssetName: "fj.tar.gz"},
		{Owner: "alice", Name: "project", Tag: "v1.0.0"},
		{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "  "},
	} {
		stub := &assetDeleterStub{}
		if err := NewDeleteAssetUseCase(stub).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
		if stub.called {
			t.Fatalf("deleter must not be called for %+v", request)
		}
	}
}

func TestDeleteAssetUseCaseRejectsNilDeleter(t *testing.T) {
	if err := NewDeleteAssetUseCase(nil).Execute(context.Background(), DeleteAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"}); err == nil {
		t.Fatal("expected internal error for nil deleter")
	}
}
