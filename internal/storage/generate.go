//go:build generate

package storage

//go:generate sh -c "rm -f ./*_mock.go"
//go:generate moq --out=depositor_mock.go . Depositor
//go:generate moq --out=keeper_mock.go . Keeper
