package main

import "net/http"

// swaggerUIHandler serves a minimal Swagger UI page (via CDN assets) pointed
// at /openapi.json, giving the same "/docs" JWT-testable experience as the
// Python FastAPI version without vendoring swagger-ui's static assets.
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
  <title>TaskFlow API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
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
