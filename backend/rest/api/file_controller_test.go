package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_parseListHandlerParams(t *testing.T) {
	testCases := map[string]struct {
		volume   string
		filename string
		groupBy  string
		sortBy   string
		orderBy  string
		offset   string
		limit    string
		want     *service.ListParams
		wantErr  bool
	}{
		"no user defined params sent": {
			volume:   "0",
			filename: "/foo",
			want: &service.ListParams{
				Volume:   123,
				Filename: "/foo",
				GroupBy:  service.DefaultGroupBy,
				SortBy:   service.DefaultSortBy,
				OrderBy:  service.DefaultOrderBy,
				Offset:   0,
				Limit:    -1,
			},
		},
		"with user defined params": {
			volume:   "123",
			filename: "/foo",
			groupBy:  "type",
			sortBy:   "name",
			orderBy:  "asc",
			offset:   "40",
			limit:    "20",
			want: &service.ListParams{
				Volume:   123,
				Filename: "/foo",
				GroupBy:  service.GroupByType,
				SortBy:   service.SortByName,
				OrderBy:  service.OrderByAsc,
				Offset:   40,
				Limit:    20,
			},
			wantErr: false,
		},
		"incorrect volume format": {
			volume:  "foo",
			wantErr: true,
		},
		"unsupported group_by": {
			volume:   "123",
			filename: "/foo",
			groupBy:  "bar",
			wantErr:  true,
		},
		"unsupported sort_by": {
			volume:   "123",
			filename: "/foo",
			sortBy:   "bar",
			wantErr:  true,
		},
		"unsupported order_by": {
			volume:   "123",
			filename: "/foo",
			orderBy:  "bar",
			wantErr:  true,
		},
		"offset parsing error": {
			volume:   "123",
			filename: "/foo",
			offset:   "bar",
			wantErr:  true,
		},
		"negative offset error": {
			volume:   "123",
			filename: "/foo",
			offset:   "-5",
			wantErr:  true,
		},
		"limit parsing error": {
			volume:   "123",
			filename: "/foo",
			limit:    "bar",
			wantErr:  true,
		},
		"zero limit error": {
			volume:   "123",
			filename: "/foo",
			limit:    "0",
			wantErr:  true,
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			engine := gin.New()
			var handlerCalled bool
			engine.GET("/files/:volume/*path", func(c *gin.Context) {
				handlerCalled = true
				got, err := parseListHandlerParams(c)
				if (err != nil) != tt.wantErr {
					t.Errorf("wantErr: %v, got: %+v", tt.want, err)
					return
				}
				assert.Equal(t, tt.want, got)
			})

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/files/%s/%s", tt.volume, tt.filename), nil)
			assert.NoError(t, err)
			queryParams := url.Values{}
			if tt.groupBy != "" {
				queryParams.Set("group_by", tt.groupBy)
			}
			if tt.sortBy != "" {
				queryParams.Set("sort_by", tt.sortBy)
			}
			if tt.orderBy != "" {
				queryParams.Set("order_by", tt.orderBy)
			}
			if tt.offset != "" {
				queryParams.Set("offset", tt.offset)
			}
			if tt.limit != "" {
				queryParams.Set("limit", tt.limit)
			}
			req.URL.RawQuery = queryParams.Encode()
			engine.ServeHTTP(recorder, req)
			assert.True(t, handlerCalled)
		})
	}
}
