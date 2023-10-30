package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"sort"
	"strings"
)

func VerifyBodySigMiddleware(next http.Handler, secret string) http.Handler {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		parameters := new(SignParameters)
		parameters.ReadFrom(r.Body)
		if !parameters.Verify(secret) {
			Error(w, NewPaymentError(10, "Несовпадение вычисленной и переданной подписи", true))
			return
		}
		r.Body = parameters
		next.ServeHTTP(w, r)
	}
	return handler
}

type KeyValuePair struct {
	Key   string
	Value string
}

func (k *KeyValuePair) Parse(s string) bool {
	params := strings.Split(s, "=")
	if len(params) != 2 {
		return false
	}
	k.Key = params[0]
	k.Value = params[1]
	return true
}

func (k KeyValuePair) String() string {
	return k.Key + "=" + k.Value
}

func (k KeyValuePair) IsSig() bool {
	return k.Key == "sig"
}

type KeyValuePairs []KeyValuePair

func (p KeyValuePairs) Len() int           { return len(p) }
func (p KeyValuePairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p KeyValuePairs) Less(i, j int) bool { return p[i].Key < p[j].Key }

func (k KeyValuePairs) Bytes() []byte {
	buf := new(bytes.Buffer)
	for _, keyValue := range k {
		buf.WriteString(keyValue.String())
	}
	return buf.Bytes()
}

type SignParameters struct {
	buf           bytes.Buffer
	KeyValuePairs KeyValuePairs
	Sig           string
}

func (s *SignParameters) Read(data []byte) (n int, err error) {
	return s.buf.Read(data)
}

func (s *SignParameters) Close() error {
	return nil
}

func (p *SignParameters) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = p.buf.ReadFrom(r)
	s := p.buf.String()
	params := strings.Split(s, "&")
	for _, keyValue := range params {
		var keyValuePair KeyValuePair
		ok := keyValuePair.Parse(keyValue)
		if !ok {
			continue
		}
		if keyValuePair.IsSig() {
			p.Sig = keyValuePair.Value
			continue
		}
		p.KeyValuePairs = append(p.KeyValuePairs, keyValuePair)
	}
	sort.Sort(p.KeyValuePairs)
	return
}

func (p SignParameters) Verify(secret string) bool {
	data := p.KeyValuePairs.Bytes()
	data = append(data, secret...)
	sig := md5.Sum(data)
	return p.Sig == hex.EncodeToString(sig[:])
}
