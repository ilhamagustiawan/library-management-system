package main

import "github.com/ilhamagustiawan/library-management-system/backend/auth-service/cmd"

// @title Library Management Auth API
// @version 1.0
// @description OAuth 2.0 authorization and user authentication service for the library management system.
// @description This service uses JWT access tokens, Authorization Code with PKCE, and HttpOnly session cookies.
// @schemes http https
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Prefix access tokens with "Bearer ".
// @securityDefinitions.basic BasicAuth
func main() {
	cmd.Execute()
}
