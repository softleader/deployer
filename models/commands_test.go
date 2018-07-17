package models

import "testing"

func TestRemoveAfterLast(t *testing.T) {
	actual := removeAfterLast(`test-0-base_eureka.1.12345`, ".")
	expect := "test-0-base_eureka.1"
	if actual != expect {
		t.Errorf("should be: %s, but was: %s", expect, actual)
	}
	actual = removeAfterLast(`test-0-base_eureka112345`, ".")
	expect = "test-0-base_eureka112345"
	if actual != expect {
		t.Errorf("should be: %s, but was: %s", expect, actual)
	}
}
