package api

// @title	Platform APIs
func PlatformRoutesV1(rt *Routes) {
	platformRoutesV1(rt)
	platformPrivateRoutesV1(rt)
}

// platformRoutesV1 configures and defines the open routes for platform specific to version 1.
//
//	@title			Golang Base Project Platform Publice API
//	@description	Service to manage the Auth of end user
//
// platformRoutesV1 Open routes
//
//	@BasePath		/platform/api/v1
func platformRoutesV1(rt *Routes) {
	serviceRoutes := rt.router.Group("/platform/api/v1")

	// Retrieve the logs controller from the provided Routes struct
	batchLogs := rt.apiCtl.BatchLogsCtl

	serviceRoutes.POST("/ingest", batchLogs.CreateBatchLogs)
	serviceRoutes.GET("/query", batchLogs.ListBatchLogs)

}

// platformPrivateRoutesV1 configures and defines the private routes for platform specific to version 1.
//
//	@title			Golang Base Project Platform Private API
//	@description	Service to manage the Auth of end user
//
// platformPrivateRoutesV1 Private routes
//
//	@BasePath		/platform/api/v1
func platformPrivateRoutesV1(rt *Routes) {
}
