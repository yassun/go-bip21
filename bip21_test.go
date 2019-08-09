package bip21_test

import (
	"reflect"
	"testing"

	"github.com/yassun/go-bip21"
)

func TestParse_Success(t *testing.T) {
	type uriParseTest struct {
		given string
		exp   *bip21.URIResources
	}

	uriParseTests := []uriParseTest{
		// Just the address
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params:    make(map[string]string),
			},
		},
		// Just the address
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params:    make(map[string]string),
			},
		},

		// Address with name
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?Label=Luke-Jr",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "Luke-Jr",
				Message:   "",
				Params:    make(map[string]string),
			},
		},

		// Request 20.30 BTC to "Luke-Jr"
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=20.3&label=Luke-Jr",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    20.3,
				Label:     "Luke-Jr",
				Message:   "",
				Params:    make(map[string]string),
			},
		},

		// Request 50 BTC with message
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=50&label=Luke-Jr&message=Donation%20for%20project%20xyz",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    50,
				Label:     "Luke-Jr",
				Message:   "Donation%20for%20project%20xyz",
				Params:    make(map[string]string),
			},
		},

		// Some future version that has variables which are (currently) not understood and required and thus invalid:
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?req-somethingyoudontunderstand=50&req-somethingelseyoudontget=999",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params: map[string]string{
					"req-somethingyoudontunderstand": "50",
					"req-somethingelseyoudontget":    "999",
				},
			},
		},

		// Some future version that has variables which are (currently) not understood but not required and thus valid:
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?somethingyoudontunderstand=50&somethingelseyoudontget=999",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params: map[string]string{
					"somethingyoudontunderstand": "50",
					"somethingelseyoudontget":    "999",
				},
			},
		},
	}

	for _, tt := range uriParseTests {
		t.Run("", func(t *testing.T) {
			got, err := bip21.Parse(tt.given)
			if err != nil {
				t.Errorf("Parse(%q) returned error %s", tt.given, err)
				return
			}
			if got.UrnScheme != tt.exp.UrnScheme {
				t.Errorf("%+v = UrnScheme: %q; exp %q", tt.given, got.UrnScheme, tt.exp.UrnScheme)
			}
			if got.Address != tt.exp.Address {
				t.Errorf("%+v = Address: %q; exp %q", tt.given, got.Address, tt.exp.Address)
			}
			if got.Amount != tt.exp.Amount {
				t.Errorf("%+v = Amount: %f; exp %f", tt.given, got.Amount, tt.exp.Amount)
			}
			if got.Label != tt.exp.Label {
				t.Errorf("%+v = Label: %q; exp %q", tt.given, got.Label, tt.exp.Label)
			}
			if got.Message != tt.exp.Message {
				t.Errorf("%+v = Message: %q; exp %q", tt.given, got.Message, tt.exp.Message)
			}
			if !reflect.DeepEqual(got.Params, tt.exp.Params) {
				t.Errorf("%+v = Params: %+v; exp %+v", tt.given, got.Params, tt.exp.Params)
			}
		})
	}
}

func TestParse_Fail(t *testing.T) {
	type uriParseTest struct {
		given string
		exp   error
	}

	uriParseTests := []uriParseTest{
		{
			"xxxx:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
			bip21.ErrInvalidUrn,
		},
		{
			"bitcoin",
			bip21.ErrInvalidUrn,
		},
		{
			"175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
			bip21.ErrInvalidUrn,
		},
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=",
			bip21.ErrInvalidAmount,
		},
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=hoge",
			bip21.ErrInvalidAmount,
		},
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=-1",
			bip21.ErrNegativeAmount,
		},
	}
	for _, tt := range uriParseTests {
		t.Run("", func(t *testing.T) {
			_, err := bip21.Parse(tt.given)
			if err != tt.exp {
				t.Errorf("Parse(%s) : exp error(%s), got error(%s)", tt.given, tt.exp, err)
			}
		})
	}
}

func TestBuildURI_Success(t *testing.T) {
	type buildURITest struct {
		exp   string
		given *bip21.URIResources
	}

	buildURITests := []buildURITest{
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params:    make(map[string]string),
			},
		},

		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?label=Luke-Jr",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "Luke-Jr",
				Message:   "",
				Params:    make(map[string]string),
			},
		},

		// Request 20.30 BTC to "Luke-Jr"
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=20.3&label=Luke-Jr",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    20.3,
				Label:     "Luke-Jr",
				Message:   "",
				Params:    make(map[string]string),
			},
		},

		// Request 50 BTC with message
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=50&label=Luke-Jr&message=Donation+for+project+xyz+%F0%9F%92%B0",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    50,
				Label:     "Luke-Jr",
				Message:   "Donation for project xyz 💰",
				Params:    make(map[string]string),
			},
		},

		// Some future version that has variables which are (currently) not understood and required and thus invalid:
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?req-somethingyoudontunderstand=50",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params: map[string]string{
					"req-somethingyoudontunderstand": "50",
				},
			},
		},

		// Some future version that has variables which are (currently) not understood but not required and thus valid:
		{
			"bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?somethingelseyoudontget=999",
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params: map[string]string{
					"somethingelseyoudontget": "999",
				},
			},
		},
	}

	for _, tt := range buildURITests {
		t.Run("", func(t *testing.T) {
			got, err := tt.given.BuildURI()
			if err != nil {
				t.Errorf("BuildURI(%+v) returned error %s", tt.given, err)
				return
			}
			if got != tt.exp {
				t.Errorf("%+v = uri: %q; exp %q", tt.given, got, tt.exp)
			}
		})
	}
}

func TestBuildURI_Fail(t *testing.T) {
	type buildURITest struct {
		given *bip21.URIResources
		exp   error
	}

	buildURITests := []buildURITest{
		{
			&bip21.URIResources{
				UrnScheme: "",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params:    make(map[string]string),
			},
			bip21.ErrInvalidUrn,
		},
		{
			&bip21.URIResources{
				UrnScheme: "xxxxx",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    0,
				Label:     "",
				Message:   "",
				Params:    make(map[string]string),
			},
			bip21.ErrInvalidUrn,
		},
		{
			&bip21.URIResources{
				UrnScheme: "bitcoin",
				Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
				Amount:    -1,
				Label:     "",
				Message:   "",
				Params:    make(map[string]string),
			},
			bip21.ErrNegativeAmount,
		},
	}

	for _, tt := range buildURITests {
		t.Run("", func(t *testing.T) {
			_, err := tt.given.BuildURI()
			if err != tt.exp {
				t.Errorf("%#v : exp error(%s), got error(%s)", tt.given, tt.exp, err)
			}
		})
	}
}
