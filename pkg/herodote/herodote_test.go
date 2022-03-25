package herodote

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ViBiOh/httputils/v4/pkg/request"
)

func TestFlags(t *testing.T) {
	cases := map[string]struct {
		want string
	}{
		"simple": {
			"Usage of simple:\n  -httpSecret string\n    \t[herodote] HTTP Secret Key for Update {SIMPLE_HTTP_SECRET}\n",
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			fs := flag.NewFlagSet(intention, flag.ContinueOnError)
			Flags(fs, "")

			var writer strings.Builder
			fs.SetOutput(&writer)
			fs.Usage()

			result := writer.String()

			if result != tc.want {
				t.Errorf("Flags() = `%s`, want `%s`", result, tc.want)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	postWithToken := httptest.NewRequest(http.MethodPost, "/", nil)
	postWithToken.Header.Add("Authorization", "testing")

	cases := map[string]struct {
		instance   App
		request    *http.Request
		want       string
		wantStatus int
		wantHeader http.Header
	}{
		"simple": {
			App{},
			httptest.NewRequest(http.MethodGet, "/", nil),
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
			http.Header{},
		},
		"post invalid token": {
			App{secret: "testing"},
			httptest.NewRequest(http.MethodPost, "/", nil),
			fmt.Sprintf("%s\n", ErrAuthentificationFailed.Error()),
			http.StatusUnauthorized,
			http.Header{},
		},
		"post valid": {
			App{secret: "testing"},
			postWithToken,
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
			http.Header{},
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
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

	cases := map[string]struct {
		args    args
		wantErr error
	}{
		"empty": {
			args{
				raw: "",
			},
			nil,
		},
		"invalid format": {
			args{
				raw: "2020-31-08",
			},
			errors.New(`unable to parse date: parsing time "2020-31-08": month out of range`),
		},
		"valid": {
			args{
				raw: "2020-08-31",
			},
			nil,
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
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
