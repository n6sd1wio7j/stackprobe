package region

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("svc-a", "us-east"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected region to be found")
	}
	if got != "us-east" {
		t.Errorf("got %q, want %q", got, "us-east")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Error("expected false for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", "eu-west"); err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestSet_EmptyRegionReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("svc-a", ""); err == nil {
		t.Error("expected error for empty region name")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("svc-a", "ap-south")
	s.Delete("svc-a")
	_, ok := s.Get("svc-a")
	if ok {
		t.Error("expected entry to be deleted")
	}
}

func TestFilter_ReturnsMatchingServices(t *testing.T) {
	s := New()
	_ = s.Set("svc-a", "us-east")
	_ = s.Set("svc-b", "us-east")
	_ = s.Set("svc-c", "eu-west")

	got := s.Filter("us-east")
	if len(got) != 2 {
		t.Errorf("expected 2 services, got %d", len(got))
	}
}

func TestFilter_UnknownRegionReturnsEmpty(t *testing.T) {
	s := New()
	_ = s.Set("svc-a", "us-east")
	got := s.Filter("ap-south")
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestByRegion_GroupsCorrectly(t *testing.T) {
	s := New()
	_ = s.Set("svc-a", "us-east")
	_ = s.Set("svc-b", "eu-west")
	_ = s.Set("svc-c", "us-east")

	groups := s.ByRegion()
	if len(groups["us-east"]) != 2 {
		t.Errorf("expected 2 services in us-east, got %d", len(groups["us-east"]))
	}
	if len(groups["eu-west"]) != 1 {
		t.Errorf("expected 1 service in eu-west, got %d", len(groups["eu-west"]))
	}
}
