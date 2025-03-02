package error

import "errors"

var (
	//Database sentinels errors
	ErrFailedToConnectDb         = errors.New("failed to connect to database")
	ErrFailedToCloseDbConnection = errors.New("failed to close database connection")
	ErrFailedToGetDBInstance     = errors.New("failed to get database instance")

	ErrEnvFileNotFound = errors.New("failed to find env file, please check the path")
	ErrFailedToLoadEnv = errors.New("failed to load env keys into config struct")
)
