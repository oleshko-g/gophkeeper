package main

import (
	"bytes"
	"context"
	"fmt"
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
		fmt.Println("received ctx.Done")
		return nil, ctx.Err()
	case res := <-resCh:
		fmt.Printf("received res %v\n", res)
		return res, nil
	case err := <-errCh:
		fmt.Printf("received err %v\n", err)
		return nil, err
	}
}
