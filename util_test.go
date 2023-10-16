package main

import "testing"

func TestToByteSize(t *testing.T) {
	v, err := toByteSize("1")
	if err != nil {
		t.Fatalf("error parsing value: %v", err)
	}
	if v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	v, err = toByteSize("2k")
	if err != nil {
		t.Fatalf("error parsing value: %v", err)
	}
	if v != 2048 {
		t.Errorf("expected 2048, got %d", v)
	}

	v, err = toByteSize("2M")
	if err != nil {
		t.Fatalf("error parsing value: %v", err)
	}
	if v != int64(2)*1024*1024 {
		t.Errorf("expected %d, got %d", int64(2)*1024*1024, v)
	}

	v, err = toByteSize("2g")
	if err != nil {
		t.Fatalf("error parsing value: %v", err)
	}
	if v != int64(2)*1024*1024*1024 {
		t.Errorf("expected %d, got %d", int64(2)*1024*1024*1024, v)
	}
}
