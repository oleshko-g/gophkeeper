package main

import (
	"bytes"
	"context"
	"io"
)

func ReadAllCtx(ctx context.Context, r io.Reader) ([]byte, error) {
	resCh := make(chan []byte, 1)
	errCh := make(chan error, 1)

	func() {
		buf := make([]byte, 2048)
		res := bytes.Buffer{}
		for {
			n, err := r.Read(buf)

			if n > 0 {
				res.Write(buf[:n])
			}

			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}
		}

		resCh <- res.Bytes()
		res.Reset()
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return nil, err
	}
}
