package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectIndexes(t *testing.T) {
	for _, c := range []struct {
		label   string
		headers []string
		pName   string
		expRes  *columnIndexes
	}{
		{
			label:   "only_subprojects",
			headers: []string{"p1", "p2:subproj1"},
			pName:   "p2",
			expRes:  nil,
		},
		{
			label:   "project_with_subprojects",
			headers: []string{"p1", "p2:subproj1", "p2", "p2:subproj3"},
			pName:   "p2",
			expRes: &columnIndexes{
				proj: 2,
				subProjs: []pair[string, int]{
					{"subproj1", 1},
					{"subproj3", 3},
				}},
		},
		{
			label:   "project_with_the_same_prefix",
			headers: []string{"p1", "p2project", "p2"},
			pName:   "p2",
			expRes:  &columnIndexes{proj: 2},
		},
		{
			label:   "has_single_project",
			headers: []string{"p1", "p2", "p3", "p3:sub1", "p1:sub1"},
			pName:   "p2",
			expRes:  &columnIndexes{proj: 1},
		},
		{
			label:   "project_is_not_here",
			headers: []string{"p1", "p2"},
			pName:   "p3",
		},
		{
			label:   "empty_project_name",
			headers: []string{"p1", "p2"},
			pName:   "",
		},
		{
			label:   "empty_headers",
			headers: nil,
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			result := projectIndexes(c.headers, c.pName)
			require.Equal(t, c.expRes, result)
		})
	}
}
