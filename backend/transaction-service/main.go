package main

import "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/cmd"

// @title Library Management Transaction API
// @version 1.0
// @description Loan, return, fine assessment, and transaction-history service.
// @description Kong authenticates human tokens; this service rechecks trusted identity headers and endpoint scopes.
// @schemes http https
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Prefix access tokens with "Bearer ".
func main() { cmd.Execute() }
