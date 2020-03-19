package data

import "testing"

func TestGetVersion(t *testing.T) {
	if s := getVersion("v1.4.4"); s != "1.4.4" {
		t.Errorf("expected %s got %s", "1.4.4", s)
	}

	if s := getVersion("v1.4.2.1"); s != "1.4.2.1" {
		t.Errorf("expected %s got %s", "1.4.2.1", s)
	}

	if s := getVersion("1.4.4"); s != "1.4.4" {
		t.Errorf("expected %s got %s", "1.4.4", s)
	}

	if s := getVersion("1.4.2.1"); s != "1.4.2.1" {
		t.Errorf("expected %s got %s", "1.4.2.1", s)
	}
}

func TestVersionCmp(t *testing.T) {
	if c := versionCmp("v1.4.4", "v1.4.2.1"); c != 1 {
		t.Errorf("expected %d got %d", 1, c)
	}

	if c := versionCmp("v1.4.2.1", "v1.4.4"); c != -1 {
		t.Errorf("expected %d got %d", -1, c)
	}

	if c := versionCmp("v1.4.4", "v1.4.4"); c != 0 {
		t.Errorf("expected %d got %d", 0, c)
	}

	if c := versionCmp("1.4.4", "1.4.2.1"); c != 1 {
		t.Errorf("expected %d got %d", 1, c)
	}

	if c := versionCmp("1.4.2.1", "1.4.4"); c != -1 {
		t.Errorf("expected %d got %d", -1, c)
	}

	if c := versionCmp("1.4.4", "1.4.4"); c != 0 {
		t.Errorf("expected %d got %d", 0, c)
	}

	if c := versionCmp("v1.4.4", "1.4.2.1"); c != 1 {
		t.Errorf("expected %d got %d", 1, c)
	}

	if c := versionCmp("v1.4.2.1", "1.4.4"); c != -1 {
		t.Errorf("expected %d got %d", -1, c)
	}

	if c := versionCmp("v1.4.4", "1.4.4"); c != 0 {
		t.Errorf("expected %d got %d", 0, c)
	}

	if c := versionCmp("1.4.4", "v1.4.2.1"); c != 1 {
		t.Errorf("expected %d got %d", 1, c)
	}

	if c := versionCmp("1.4.2.1", "v1.4.4"); c != -1 {
		t.Errorf("expected %d got %d", -1, c)
	}

	if c := versionCmp("1.4.4", "v1.4.4"); c != 0 {
		t.Errorf("expected %d got %d", 0, c)
	}
}
