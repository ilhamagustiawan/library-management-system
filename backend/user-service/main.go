package main

import "github.com/ilhamagustiawan/library-management-system/backend/user-service/cmd"

// @title Library Management User API
// @version 1.0
// @description Member registration and profile service.
// @schemes http https
// @BasePath /
func main() { cmd.Execute() }
