package hash

import (
	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/bcrypt"
)

const hashIDMinLength = 30

type Hasher interface {
	EncodeInt64(int64) (string, error)
	DecodeInt64(string) (int64, error)
	Password(password string) (string, error)
	CheckPassword(password, hash string) bool
}

type HasherImpl struct {
	enc *hashids.HashID
}

func NewHasher(secret string) *HasherImpl {
	hd := hashids.NewData()
	hd.Salt = secret
	hd.MinLength = hashIDMinLength
	h, _ := hashids.NewWithData(hd)
	return &HasherImpl{enc: h}
}

func (h *HasherImpl) EncodeInt64(i int64) (string, error) {
	return h.enc.EncodeInt64([]int64{i})
}

func (h *HasherImpl) DecodeInt64(s string) (int64, error) {
	nums, err := h.enc.DecodeInt64WithError(s)
	if err != nil {
		return 0, err
	}
	return nums[0], nil
}

func (h *HasherImpl) Password(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *HasherImpl) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
