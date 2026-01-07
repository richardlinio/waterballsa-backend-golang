package middleware

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// Security creates and returns a security headers middleware with hardcoded configuration
func Security() gin.HandlerFunc {
	return secure.New(secure.Config{
		// SSL redirect is disabled because reverse proxy (nginx/Zeabur) handles HTTPS
		SSLRedirect: false,
		// SSLProxyHeaders is used to detect if the request came through a reverse proxy with HTTPS
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},

		// HSTS (HTTP Strict Transport Security) - Forces HTTPS for 1 year
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,

		// Frame Options - Prevents clickjacking attacks by denying the page from being displayed in iframes
		FrameDeny: true,

		// Content Type Options - Prevents MIME type sniffing
		ContentTypeNosniff: true,

		// XSS Filter - Enables browser's XSS filter (legacy, but still useful for older browsers)
		BrowserXssFilter: true,

		// Content Security Policy - Relaxed setting for development flexibility
		// Allows inline scripts/styles and images from any source
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval'; img-src * data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval'",

		// Referrer Policy - Controls how much referrer information is sent with requests
		ReferrerPolicy: "strict-origin-when-cross-origin",

		// IE No Open - Prevents Internet Explorer from executing downloads in the site's context
		IENoOpen: true,

		// Development mode is disabled to ensure security headers are always set
		IsDevelopment: false,
	})
}
