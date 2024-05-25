package auth

import "net/http"

type LoginPassword struct {
}

func (l *LoginPassword) Auth(w http.ResponseWriter, r *http.Request) {

}
