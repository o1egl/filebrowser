package hash

import "github.com/google/wire"

var Set = wire.NewSet(
	NewHasher,
	wire.Bind(new(Hasher), new(*HasherImpl)),
)
