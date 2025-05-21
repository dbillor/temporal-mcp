package auth

import (
    "encoding/json"
    "errors"
    "github.com/golang-jwt/jwt/v5"
)

type Verifier interface {
    Verify(tokenString string) (*jwt.RegisteredClaims, jwt.MapClaims, error)
}

type HMACVerifier struct {
    Key []byte
}

func (v *HMACVerifier) Verify(tok string) (*jwt.RegisteredClaims, jwt.MapClaims, error) {
    mc := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tok, mc, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return v.Key, nil
    })
    if err != nil {
        return nil, nil, err
    }
    if !token.Valid {
        return nil, nil, errors.New("invalid token")
    }
    rc := &jwt.RegisteredClaims{}
    b, _ := json.Marshal(mc)
    _ = json.Unmarshal(b, rc)
    return rc, mc, nil
}
