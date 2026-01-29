package repository

import "context"

type SSHAuditRepository interface {
	Log(ctx context.Context, userID, connectionID int64, remoteIP, event, detail string) error
}
