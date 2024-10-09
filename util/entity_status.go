package util

import (
	"github.com/technoshantoms/mccs-alpha-api/global/constant"
)

func IsValidStatus(status string) bool {
	if status == constant.Entity.Pending ||
		status == constant.Entity.Rejected ||
		status == constant.Entity.Accepted ||
		status == constant.Trading.Pending ||
		status == constant.Trading.Accepted ||
		status == constant.Trading.Rejected {
		return true
	}
	return false
}

// IsAcceptedStatus checks if the entity status is accpeted.
func IsAcceptedStatus(status string) bool {
	if status == constant.Entity.Accepted ||
		status == constant.Trading.Pending ||
		status == constant.Trading.Accepted ||
		status == constant.Trading.Rejected {
		return true
	}
	return false
}

func IsTradingAccepted(status string) bool {
	if status == constant.Trading.Accepted {
		return true
	}
	return false
}
