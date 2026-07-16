package main

import (
	"embed"
	"net/http"
)

//go:embed swaggerui/swagger-ui.css swaggerui/swagger-ui-bundle.js
var swaggerUIAssets embed.FS

// swaggerUIFileHandler serves the vendored swagger-ui-dist assets embedded
// in the binary at build time, so /docs works fully offline (no unpkg CDN).
var swaggerUIFileHandler = http.FileServer(http.FS(swaggerUIAssets))

// swaggerUIHandler serves a minimal Swagger UI page pointed at /openapi.json,
// using the locally embedded swagger-ui assets instead of a CDN.
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
  <title>TaskFlow API Docs</title>
  <link rel="stylesheet" href="/docs/swaggerui/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="/docs/swaggerui/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/openapi.json',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
      });
    };
  </script>
</body>
</html>`))
}
