package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func isRootDomain(name string) bool {
	n := []rune(name)
	var count int
	for _, v := range n {
		if v == '.' {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func randString(length int) (string, error) {
	wr := errWrapper("error making random string")
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", wr(err, "error gathering token entropy")
	} else {
		return base64.RawURLEncoding.EncodeToString(b), nil
	}
}

func mustRandString(numBytes int) string {
	if s, err := randString(numBytes); err != nil {
		log.Fatal(err)
		return ""
	} else {
		return s
	}
}

type ErrWrapFunc func(error, ...string) error

func subWrapper(w ErrWrapFunc, outerPrefix string) ErrWrapFunc {
	return func(err error, prefixes ...string) error {
		return w(err, append(prefixes, outerPrefix)...)
	}
}

func errWrapper(outerPrefix string) ErrWrapFunc {
	return func(err error, prefixes ...string) error {
		for _, v := range prefixes {
			err = errors.Wrap(err, v)
		}
		return errors.Wrap(err, outerPrefix)
	}
}

func backoffLogNW() func(err error, delay time.Duration) {
	return func(err error, delay time.Duration) {
		log.Println(err)
	}
}

func backoffLog(w ErrWrapFunc) func(err error, delay time.Duration) {
	return func(err error, delay time.Duration) {
		log.Println(w(err))
	}
}
