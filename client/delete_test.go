package main

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pb "github.com/oleshko-g/gophkeeper/api/v1"
	grpc "github.com/oleshko-g/gophkeeper/internal/grpc/client"
	"github.com/spf13/cobra"
)

// TODO:
//   - [*] mock client
//   - [] temp file storage
//   - [] mock input
func Test_deleteRunE(t *testing.T) {
	testApp := app
	_ = testApp
	app.Client = &grpc.Client{
		KeeperClient: &grpc.KeeperClientMock{
			PurgeSecretFunc: func(ctx context.Context, _ *pb.PurgeSecretRequest, _ ...grpc.CallOption) (*pb.PurgeSecretResponse, error) {
				for {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					default:
						return &pb.PurgeSecretResponse{}, nil
					}
				}
			},
		},
	}

	type result struct {
		Err error
	}
	tests := []struct {
		name  string
		input io.Reader
		delay time.Duration
		want  result
		got   result
	}{
		{
			name:  "deadline exceeded while reading input",
			delay: 2 * time.Second,
			input: strings.NewReader("019e8499-ef8e-7eff-80f4-a2a63e351893"),
			want:  result{Err: context.DeadlineExceeded},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
			defer cancel()
			cmd := &cobra.Command{}
			cmd.SetContext(ctx)
			cmd.SetIn(v.input)
			v.got.Err = testApp.deleteRunE(cmd, nil)
			if !cmp.Equal(v.want, v.got) {
				cmp.Diff(v.want, v.got)
				t.Fail()
			}
		})
	}
}
