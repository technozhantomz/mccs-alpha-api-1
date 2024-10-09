package balancecheck

import (
	"time"

	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/pg"
	"github.com/technoshantoms/mccs-alpha-api/internal/pkg/email"
	"github.com/technoshantoms/mccs-alpha-api/util/l"
	"go.uber.org/zap"
)

// Run will check whether the last past 5 hours the sum of the balance in the posting table is zero.
func Run() {
	to := time.Now()
	from := to.Add(-5 * time.Hour)

	postings, err := pg.Posting.FindInRange(from, to)
	if err != nil {
		l.Logger.Error("checking balance failed", zap.Error(err))
		return
	}

	var sum float64
	for _, p := range postings {
		sum += p.Amount
	}

	if sum != 0.0 {
		email.Balance.SendNonZeroBalanceEmail(&email.NonZeroBalanceEmail{
			From: from,
			To:   to,
		})
	}
}
