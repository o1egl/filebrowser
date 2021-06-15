package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_parseListHandlerParams(t *testing.T) {
	testCases := map[string]struct {
		filename string
		groupBy  string
		sortBy   string
		orderBy  string
		offset   string
		limit    string
		want     *ListHandlerParams
		wantErr  bool
	}{
		"no user defined params sent": {
			filename: "/foo",
			want: &ListHandlerParams{
				Filename: "/foo",
				GroupBy:  defaultGroupBy,
				SortBy:   defaultSortBy,
				OrderBy:  defaultOrderBy,
				Offset:   0,
				Limit:    -1,
			},
		},
		"with user defined params": {
			filename: "/foo",
			groupBy:  "type",
			sortBy:   "name",
			orderBy:  "asc",
			offset:   "40",
			limit:    "20",
			want: &ListHandlerParams{
				Filename: "/foo",
				GroupBy:  GroupByType,
				SortBy:   SortByName,
				OrderBy:  OrderByAsc,
				Offset:   40,
				Limit:    20,
			},
			wantErr: false,
		},
		"unsupported group_by": {
			filename: "/foo",
			groupBy:  "bar",
			wantErr:  true,
		},
		"unsupported sort_by": {
			filename: "/foo",
			sortBy:   "bar",
			wantErr:  true,
		},
		"unsupported order_by": {
			filename: "/foo",
			orderBy:  "bar",
			wantErr:  true,
		},
		"offset parsing error": {
			filename: "/foo",
			offset:   "bar",
			wantErr:  true,
		},
		"negative offset error": {
			filename: "/foo",
			offset:   "-5",
			wantErr:  true,
		},
		"limit parsing error": {
			filename: "/foo",
			limit:    "bar",
			wantErr:  true,
		},
		"zero limit error": {
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
			engine.GET("/files/*path", func(c *gin.Context) {
				handlerCalled = true
				got, err := parseListHandlerParams(c)
				if (err != nil) != tt.wantErr {
					t.Errorf("wantErr: %v, got: %+v", tt.want, err)
					return
				}
				assert.Equal(t, tt.want, got)
			})

			req, err := http.NewRequest(http.MethodGet, "/files"+tt.filename, nil)
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
