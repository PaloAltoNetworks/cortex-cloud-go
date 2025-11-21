// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package version

import "testing"

func BenchmarkUserAgent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UserAgent("platform", "1.0.0")
	}
}

func BenchmarkUserAgentWithCustom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UserAgentWithCustom("platform", "1.0.0", "custom/1.0")
	}
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Info()
	}
}
