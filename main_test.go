// spserve - Serve files to current network with ease.
// Copyright (C) 2023 Safal Piya

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import "testing"

func Test_getRootPathCleaned(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Valid directory",
			args:    args{
				root: ".",
			},
			want:    ".",
			wantErr: false,
		},
		{
			name:    "Regular file",
			args:    args{
				root: "./go.mod",
			},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid location",
			args:    args{
				root: "foo/bar",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRootPathCleaned(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRootPathCleaned() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRootPathCleaned() = %v, want %v", got, tt.want)
			}
		})
	}
}
