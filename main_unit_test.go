package main

import "testing"

func TestExtractRecordName(t *testing.T) {
	tests := []struct {
		name string
		fqdn string
		zone string
		want string
	}{
		{
			name: "subdomain record",
			fqdn: "_acme-challenge.www.example.com.",
			zone: "example.com.",
			want: "_acme-challenge.www",
		},
		{
			name: "zone apex returns at",
			fqdn: "example.com.",
			zone: "example.com.",
			want: "@",
		},
		{
			name: "mismatched zone does not partially trim",
			fqdn: "_acme-challenge.sample.com",
			zone: "ample.com",
			want: "_acme-challenge.sample.com",
		},
		{
			name: "empty zone returns fqdn",
			fqdn: "_acme-challenge.example.com.",
			zone: "",
			want: "_acme-challenge.example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractRecordName(tc.fqdn, tc.zone)
			if got != tc.want {
				t.Fatalf("extractRecordName(%q, %q) = %q, want %q", tc.fqdn, tc.zone, got, tc.want)
			}
		})
	}
}
