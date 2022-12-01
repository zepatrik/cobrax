// Copyright © 2022 Patrik Neu
// SPDX-License-Identifier: Apache-2.0

package cobrax

import (
	"errors"
	"github.com/spf13/cobra"
	"testing"
)

func TestExecuteRootIntegration(t *testing.T) {
	lastExitCode := -1
	osExit = func(code int) {
		lastExitCode = code
	}

	for _, tc := range []struct {
		desc             string
		expectedExitCode int
		err              error
	}{
		{
			desc:             "no error",
			expectedExitCode: 0,
		},
		{
			desc:             "ErrNoPrintButFail",
			err:              ErrNoPrintButFail,
			expectedExitCode: 1,
		},
		{
			desc:             "WithExitCode",
			err:              WithExitCode(errors.New("foo"), 2),
			expectedExitCode: 2,
		},
		{
			desc:             "ErrNoPrintButFail with WithExitCode",
			err:              WithExitCode(ErrNoPrintButFail, 3),
			expectedExitCode: 3,
		},
		{
			desc:             "WithExitCode 0",
			err:              WithExitCode(errors.New("foo"), 0),
			expectedExitCode: 0,
		},
		{
			desc:             "other error",
			err:              errors.New("some error"),
			expectedExitCode: 1,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ExecuteRootCommand(&cobra.Command{RunE: func(*cobra.Command, []string) error { return tc.err }})
			if lastExitCode != tc.expectedExitCode {
				t.Errorf("Expected exit code %d, got %d", tc.expectedExitCode, lastExitCode)
				t.FailNow()
			}
		})
	}
}
