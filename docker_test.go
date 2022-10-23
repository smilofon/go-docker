package docker

import "testing"

func TestDocker(t *testing.T) {
	ctx, err := SetUp("alpine:latest", []string{}, []string{"tail", "-f", "/dev/null"})
	if err != nil {
		t.Fatal(err)
	}
	err = TearDown(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
