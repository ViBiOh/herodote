package herodote

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ViBiOh/herodote/pkg/store"
	"github.com/ViBiOh/httputils/v3/pkg/request"
)

func TestFlags(t *testing.T) {
	var cases = []struct {
		intention string
		want      string
	}{
		{
			"simple",
			"Usage of simple:\n  -httpSecret string\n    \t[herodote] HTTP Secret Key for Update {SIMPLE_HTTP_SECRET}\n",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			fs := flag.NewFlagSet(testCase.intention, flag.ContinueOnError)
			Flags(fs, "")

			var writer strings.Builder
			fs.SetOutput(&writer)
			fs.Usage()

			result := writer.String()

			if result != testCase.want {
				t.Errorf("Flags() = `%s`, want `%s`", result, testCase.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	emptyString := ""
	secretString := "testing"

	type args struct {
		config Config
		store  store.App
	}

	var cases = []struct {
		intention string
		args      args
		want      App
		wantErr   error
	}{
		{
			"empty param",
			args{
				config: Config{secret: &emptyString},
			},
			nil,
			errors.New("http secret is required"),
		},
		{
			"empty databse",
			args{
				config: Config{secret: &secretString},
			},
			nil,
			errors.New("store is required"),
		},
		{
			"valid",
			args{
				config: Config{secret: &secretString},
				store:  store.New(nil),
			},
			app{
				secret: "testing",
				store:  store.New(nil),
				colors: make(map[string]string),
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.intention, func(t *testing.T) {
			got, gotErr := New(tc.args.config, tc.args.store)

			failed := false

			if tc.wantErr == nil && gotErr != nil {
				failed = true
			} else if tc.wantErr != nil && gotErr == nil {
				failed = true
			} else if tc.wantErr != nil && !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
				failed = true
			} else if !reflect.DeepEqual(got, tc.want) {
				failed = true
			}

			if failed {
				t.Errorf("New() = (%+v, `%s`), want (%+v, `%s`)", got, gotErr, tc.want, tc.wantErr)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	postWithToken := httptest.NewRequest(http.MethodPost, "/", nil)
	postWithToken.Header.Add("Authorization", "testing")

	var cases = []struct {
		intention  string
		instance   app
		request    *http.Request
		want       string
		wantStatus int
		wantHeader http.Header
	}{
		{
			"simple",
			app{},
			httptest.NewRequest(http.MethodGet, "/", nil),
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
			http.Header{},
		},
		{
			"post invalid token",
			app{secret: "testing"},
			httptest.NewRequest(http.MethodPost, "/", nil),
			fmt.Sprintf("%s\n", ErrAuthentificationFailed.Error()),
			http.StatusUnauthorized,
			http.Header{},
		},
		{
			"post valid",
			app{secret: "testing"},
			postWithToken,
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
			http.Header{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.intention, func(t *testing.T) {
			writer := httptest.NewRecorder()
			tc.instance.Handler().ServeHTTP(writer, tc.request)

			if got := writer.Code; got != tc.wantStatus {
				t.Errorf("Handler = %d, want %d", got, tc.wantStatus)
			}

			if got, _ := request.ReadBodyResponse(writer.Result()); string(got) != tc.want {
				t.Errorf("Handler = `%s`, want `%s`", string(got), tc.want)
			}

			for key := range tc.wantHeader {
				want := tc.wantHeader.Get(key)
				if got := writer.Header().Get(key); got != want {
					t.Errorf("`%s` Header = `%s`, want `%s`", key, got, want)
				}
			}
		})
	}
}

func TestCheckDate(t *testing.T) {
	type args struct {
		raw string
	}

	var cases = []struct {
		intention string
		args      args
		wantErr   error
	}{
		{
			"empty",
			args{
				raw: "",
			},
			nil,
		},
		{
			"invalid format",
			args{
				raw: "2020-31-08",
			},
			errors.New(`unable to parse date: parsing time "2020-31-08": month out of range`),
		},
		{
			"valid",
			args{
				raw: "2020-08-31",
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.intention, func(t *testing.T) {
			gotErr := checkDate(tc.args.raw)

			failed := false

			if tc.wantErr == nil && gotErr != nil {
				failed = true
			} else if tc.wantErr != nil && gotErr == nil {
				failed = true
			} else if tc.wantErr != nil && !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
				failed = true
			}

			if failed {
				t.Errorf("checkDate() = `%s`, want `%s`", gotErr, tc.wantErr)
			}
		})
	}
}
