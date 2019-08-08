package bip21

import (
	"bytes"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrInvalidUrn     = errors.New("invalid URN.")
	ErrInvalidAmount  = errors.New("invalid amount.")
	ErrNegativeAmount = errors.New("amount can not be negative.")
)

type UriResources struct {
	UrnScheme string
	Address   string
	Amount    float64
	Label     string
	Message   string
	Params    map[string]string
}

func (u UriResources) BuildUri() (string, error) {

	if u.UrnScheme != "bitcoin" {
		return "", ErrInvalidUrn
	}

	b := &bytes.Buffer{}

	_, err := b.WriteString(u.UrnScheme)
	if err != nil {
		return "", err
	}

	err = b.WriteByte(':')
	if err != nil {
		return "", err
	}

	_, err = b.WriteString(u.Address)
	if err != nil {
		return "", err
	}

	if u.Amount < 0 {
		return "", ErrNegativeAmount
	}

	ps := url.Values{}

	if u.Amount != 0 {
		ps.Add("amount", strconv.FormatFloat(u.Amount, 'f', -1, 64))
	}

	if len(u.Label) != 0 {
		ps.Add("label", u.Label)
	}

	if len(u.Message) != 0 {
		ps.Add("message", u.Message)
	}

	if u.Params != nil {
		for k, v := range u.Params {
			ps.Add(k, v)
		}
	}

	if len(ps) != 0 {
		err = b.WriteByte('?')
		if err != nil {
			return "", err
		}
		_, err = b.WriteString(ps.Encode())
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

func Parse(uri string) (*UriResources, error) {
	s := strings.Split(uri, ":")
	if s[0] != "bitcoin" || len(s) != 2 {
		return nil, ErrInvalidUrn
	}
	u := &UriResources{
		UrnScheme: "bitcoin",
		Params:    make(map[string]string),
	}

	u.Address = parseAddress(uri, u.UrnScheme)
	if strings.Index(uri, "?") == -1 {
		return u, nil
	}

	p := parseParams(uri, u.UrnScheme, u.Address)

	if v, ok := p["amount"]; ok {
		a, err := parseAmount(v)
		if err != nil {
			return nil, err
		}
		u.Amount = a
		delete(p, "amount")
	}

	if v, ok := p["label"]; ok {
		u.Label = v
		delete(p, "label")
	}

	if v, ok := p["message"]; ok {
		u.Message = v
		delete(p, "message")
	}

	u.Params = p

	return u, nil
}

func parseAmount(amount string) (float64, error) {

	a, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return -1, ErrInvalidAmount
	}
	if a < 0 {
		return -1, ErrNegativeAmount
	}

	return a, nil
}

func parseAddress(uri string, urn string) string {

	i := strings.Index(uri, "?")
	if i == -1 {
		return uri[len(urn)+1:]
	} else {
		return uri[len(urn)+1 : i]
	}
}

func parseParams(uri string, urn string, address string) map[string]string {

	ps := make(map[string]string)

	qp := uri[len(urn)+1+len(address)+1:]

	query := strings.Split(qp, "&")

	for _, q := range query {
		p := strings.Split(q, "=")
		if len(p) < 2 {
			continue
		}
		ps[strings.ToLower(p[0])] = p[1]
	}

	return ps
}
