package api

import (
	"fmt"
	"net/http"
	"net/url"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func swaggerHandler(serverURI *url.URL) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL(serverURI.JoinPath("/api/swagger/doc.json").String()),
		httpSwagger.BeforeScript(`const UrlMutatorPlugin = (system) => ({
  rootInjects: {
    setScheme: (scheme) => {
      const jsonSpec = system.getState().toJSON().spec.json;
      const schemes = Array.isArray(scheme) ? scheme : [scheme];
      const newJsonSpec = Object.assign({}, jsonSpec, { schemes });

      return system.specActions.updateJsonSpec(newJsonSpec);
    },
    setHost: (host) => {
      const jsonSpec = system.getState().toJSON().spec.json;
      const newJsonSpec = Object.assign({}, jsonSpec, { host });

      return system.specActions.updateJsonSpec(newJsonSpec);
    },
    setBasePath: (basePath) => {
      const jsonSpec = system.getState().toJSON().spec.json;
      const newJsonSpec = Object.assign({}, jsonSpec, { basePath });

      return system.specActions.updateJsonSpec(newJsonSpec);
    }
  }
});`),
		httpSwagger.Plugins([]string{"UrlMutatorPlugin"}),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": fmt.Sprintf(`() => {
    window.ui.setScheme('%s');
    window.ui.setHost('%s');
    window.ui.setBasePath('%s');
  }`, serverURI.Scheme, serverURI.Host, serverURI.Path),
		}),
	)
}
