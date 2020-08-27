package errors

import "errors"

var (
	ErrServExists   = errors.New("The service already exists")
	ErrServNotExist = errors.New("Service does not exist")

	ErrUnknownMethod                 = errors.New("etcdserver: unknown method")
	ErrStopped                       = errors.New("etcdserver: server stopped")
	ErrCanceled                      = errors.New("etcdserver: request cancelled")
	ErrTimeout                       = errors.New("etcdserver: request timed out")
	ErrTimeoutDueToLeaderFail        = errors.New("etcdserver: request timed out, possibly due to previous leader failure")
	ErrTimeoutDueToConnectionLost    = errors.New("etcdserver: request timed out, possibly due to connection lost")
	ErrTimeoutLeaderTransfer         = errors.New("etcdserver: request timed out, leader transfer took too long")
	ErrLeaderChanged                 = errors.New("etcdserver: leader changed")
	ErrNotEnoughStartedMembers       = errors.New("etcdserver: re-configuration failed due to not enough started members")
	ErrLearnerNotReady               = errors.New("etcdserver: can only promote a learner member which is in sync with leader")
	ErrNoLeader                      = errors.New("etcdserver: no leader")
	ErrNotLeader                     = errors.New("etcdserver: not leader")
	ErrRequestTooLarge               = errors.New("etcdserver: request is too large")
	ErrNoSpace                       = errors.New("etcdserver: no space")
	ErrTooManyRequests               = errors.New("etcdserver: too many requests")
	ErrUnhealthy                     = errors.New("etcdserver: unhealthy cluster")
	ErrKeyNotFound                   = errors.New("etcdserver: key not found")
	ErrCorrupt                       = errors.New("etcdserver: corrupt cluster")
	ErrBadLeaderTransferee           = errors.New("etcdserver: bad leader transferee")
	ErrClusterVersionUnavailable     = errors.New("etcdserver: cluster version not found during downgrade")
	ErrWrongDowngradeVersionFormat   = errors.New("etcdserver: wrong downgrade target version format")
	ErrInvalidDowngradeTargetVersion = errors.New("etcdserver: invalid downgrade target version")
	ErrDowngradeInProcess            = errors.New("etcdserver: cluster has a downgrade job in progress")
	ErrNoInflightDowngrade           = errors.New("etcdserver: no inflight downgrade job")
)
