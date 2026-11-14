package jsxtranspiler

import "testing"

func TestCompileAllowsCommentInsideOpeningTag(t *testing.T) {
	source := "<input\n      // Optional: live search when typing\n      value={query}\n    />"
	got, err := Compile(source, "")
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}
	want := `d("input", { "value": query })`
	if got != want {
		t.Errorf("Compile() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestCompileAllowsBlockCommentInsideOpeningTag(t *testing.T) {
	source := "<div /* spacing */ class=\"x\" />"
	got, err := Compile(source, "")
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}
	want := `d("div", { "class": "x" })`
	if got != want {
		t.Errorf("Compile() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
