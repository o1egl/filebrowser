package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseListHandlerParams(t *testing.T) {
	hasher := hash.NewHasher("secret")
	hashed123, err := hasher.EncodeInt64(123)
	require.NoError(t, err)

	testCases := map[string]struct {
		volume   string
		filename string
		groupBy  string
		sortBy   string
		orderBy  string
		offset   string
		limit    string
		want     *filebrowser.ListParams
		wantErr  bool
		errMsg   string
	}{
		"no user defined params sent": {
			volume:   homeVolumeID,
			filename: "foo",
			want: &filebrowser.ListParams{
				Volume:   filebrowser.HomeVolumeID,
				Filename: "/foo",
				GroupBy:  filebrowser.DefaultGroupBy,
				SortBy:   filebrowser.DefaultSortBy,
				OrderBy:  filebrowser.DefaultOrderBy,
				Offset:   0,
				Limit:    -1,
			},
		},
		"with user defined params": {
			volume:   hashed123,
			filename: "foo",
			groupBy:  "type",
			sortBy:   "name",
			orderBy:  "asc",
			offset:   "40",
			limit:    "20",
			want: &filebrowser.ListParams{
				Volume:   123,
				Filename: "/foo",
				GroupBy:  filebrowser.GroupByType,
				SortBy:   filebrowser.SortByName,
				OrderBy:  filebrowser.OrderByAsc,
				Offset:   40,
				Limit:    20,
			},
		},
		"incorrect volume format": {
			volume:  "foo",
			wantErr: true,
			errMsg:  "incorrect volume id",
		},
		"unsupported group_by": {
			volume:   homeVolumeID,
			filename: "/foo",
			groupBy:  "bar",
			wantErr:  true,
			errMsg:   "incorrect group_by param",
		},
		"unsupported sort_by": {
			volume:   homeVolumeID,
			filename: "/foo",
			sortBy:   "bar",
			wantErr:  true,
			errMsg:   "incorrect sort_by param",
		},
		"unsupported order_by": {
			volume:   homeVolumeID,
			filename: "/foo",
			orderBy:  "bar",
			wantErr:  true,
			errMsg:   "incorrect order_by param",
		},
		"offset parsing error": {
			volume:   homeVolumeID,
			filename: "/foo",
			offset:   "bar",
			wantErr:  true,
			errMsg:   "incorrect offset param",
		},
		"negative offset error": {
			volume:   homeVolumeID,
			filename: "/foo",
			offset:   "-5",
			wantErr:  true,
			errMsg:   "offset must be positive",
		},
		"limit parsing error": {
			volume:   homeVolumeID,
			filename: "/foo",
			limit:    "bar",
			wantErr:  true,
			errMsg:   "incorrect limit param",
		},
		"zero limit error": {
			volume:   homeVolumeID,
			filename: "/foo",
			limit:    "0",
			wantErr:  true,
			errMsg:   "limit must be greater than 0",
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
				fc := newFileController(nil, hasher)
				got, err := fc.parseListHandlerParams(c)
				if (err != nil) != tt.wantErr {
					t.Errorf("wantErr: %v, got: %+v", tt.want, err)
					return
				}
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
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
