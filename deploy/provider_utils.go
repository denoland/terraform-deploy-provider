// Copyright 2021 Deno Land Inc. All rights reserved. MIT License.
package deploy

import (
	"os"
	"testing"
)

var testToken string = os.Getenv("DEPLOY_TOKEN")

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DEPLOY_TOKEN"); v == "" {
		t.Fatal("DEPLOY_TOKEN must be set for acceptance tests")
	}
}
