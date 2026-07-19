package main

import "github.com/ilhamagustiawan/library-management-system/backend/book-service/cmd"

// @title Library Management Book API
// @version 1.0
// @description Catalog, inventory, and atomic stock reservation service.
// @schemes http https
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() { cmd.Execute() }
