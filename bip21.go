package bip21

import (
	"bytes"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

var (
	// ErrInvalidUrn is returned when urn invalid
	ErrInvalidUrn = errors.New("invalid urn")
	// ErrInvalidAmount is returned when amount is invalid
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrNegativeAmount is returned when amount is negative
	ErrNegativeAmount = errors.New("amount can not be negative")
)

// URIResources is the input value when building URI
type URIResources struct {
	UrnScheme string
	Address   string
	Amount    float64
	Label     string
	Message   string
	Params    map[string]string
}

// BuildURI uses URIResources to generate a BIP21 standard URI.
func (u URIResources) BuildURI() (string, error) {
	if u.UrnScheme != "bitcoin" {
		return "", ErrInvalidUrn
	}

	b := &bytes.Buffer{}

	if _, err := b.WriteString(u.UrnScheme); err != nil {
		return "", err
	}

	if err := b.WriteByte(':'); err != nil {
		return "", err
	}

	if _, err := b.WriteString(u.Address); err != nil {
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
		if err := b.WriteByte('?'); err != nil {
			return "", err
		}
		if _, err := b.WriteString(ps.Encode()); err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

// Parse parses the BIP21 standard URI and returns URIResources.
func Parse(uri string) (*URIResources, error) {
	s := strings.Split(uri, ":")
	if s[0] != "bitcoin" || len(s) != 2 {
		return nil, ErrInvalidUrn
	}
	u := &URIResources{
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
	}

	return uri[len(urn)+1 : i]
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
