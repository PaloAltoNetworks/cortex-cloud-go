// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package platform

//import (
//	"context"
//	//"fmt"
//	//"strconv"
//	"testing"
//	//"time"
//
//	//"github.com/PaloAltoNetworks/cortex-cloud-go/types"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestAccRolesList(t *testing.T) {
//	client := setupAcceptanceTest(t)
//	ctx := context.Background()
//	//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
//
//	// Read
//	listResp, err := client.ListRoles(ctx, []string{"IT Admin"})
//	if err != nil {
//		t.Fatalf("error fetching roles: %s", err.Error())
//	}
//
//	// Check
//	assert.NotNil(t, listResp)
//}
