package app

const (
	errMsg  = "error"
	respMsg = "response"

	noDBErr = "No database provided"

	appStart = "Starting application"

	httpServerStart = "HTTP Server started"
	httpServerFatal = "HTTP server failed"

	grpcServerStart = "gRPC Server started"
	grpcServerFatal = "GRPC server failed"

	grpcGetMerchRepoReadError = "GetMerch repo read error"
	grpcGetMerchStreamError   = "GetMerch stream send error"
	grpcGetMerchSuccess       = "GetMerch success"

	grpcPostMerchBatchError = "PostMerch batch save error"

	grpcEOF = "gRPC EOF"

	grpcReceiveError   = "gRPC receive error"
	grpcReceiveSuccess = "gRPC receive success"

	notificationServiceError   = "Notification service start failed"
	notificationServiceSuccess = "Notification service started"
)
