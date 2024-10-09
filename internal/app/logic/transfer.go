package logic

import (
	"errors"
	"fmt"
	"math"

	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/es"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/pg"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/types"
)

type transfer struct{}

var Transfer = &transfer{}

func (t *transfer) Search(req *types.SearchTransferReq) (*types.SearchTransferRespond, error) {
	transfers, err := pg.Journal.Search(req)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

// POST /transfers
// POST /admin/transfers

func (t *transfer) CheckBalance(payer, payee string, amount float64) error {
	from, err := pg.Account.FindByAccountNumber(payer)
	if err != nil {
		return err
	}
	to, err := pg.Account.FindByAccountNumber(payee)
	if err != nil {
		return err
	}

	exceed, err := BalanceLimit.IsExceedLimit(from.AccountNumber, from.Balance-amount)
	if err != nil {
		return err
	}
	if exceed {
		amount, err := t.maxNegativeBalanceCanBeTransferred(from)
		if err != nil {
			return err
		}
		return errors.New("Sender will exceed its credit limit." + " The maximum amount that can be sent is: " + fmt.Sprintf("%.2f", amount))
	}

	exceed, err = BalanceLimit.IsExceedLimit(to.AccountNumber, to.Balance+amount)
	if err != nil {
		return err
	}
	if exceed {
		amount, err := t.maxPositiveBalanceCanBeTransferred(to)
		if err != nil {
			return err
		}
		return errors.New("Receiver will exceed its maximum balance limit." + " The maximum amount that can be received is: " + fmt.Sprintf("%.2f", amount))
	}

	return nil
}

// PATCH /transfers/{transferID}

func (t *transfer) FindByID(transferID string) (*types.Journal, error) {
	journal, err := pg.Journal.FindByID(transferID)
	if err != nil {
		return nil, err
	}
	return journal, nil
}

// POST /transfers

func (t *transfer) Propose(req *types.TransferReq) (*types.Journal, error) {
	journal, err := pg.Journal.Propose(req)
	if err != nil {
		return nil, err
	}
	err = es.Journal.Create(journal)
	if err != nil {
		return nil, err
	}
	return journal, nil
}

func (t *transfer) maxPositiveBalanceCanBeTransferred(a *types.Account) (float64, error) {
	maxPosBal, err := BalanceLimit.GetMaxPosBalance(a.AccountNumber)
	if err != nil {
		return 0, err
	}
	if a.Balance >= 0 {
		return maxPosBal - a.Balance, nil
	}
	return math.Abs(a.Balance) + maxPosBal, nil
}

func (t *transfer) maxNegativeBalanceCanBeTransferred(a *types.Account) (float64, error) {
	maxNegBal, err := BalanceLimit.GetMaxNegBalance(a.AccountNumber)
	if err != nil {
		return 0, err
	}
	if a.Balance >= 0 {
		return a.Balance + maxNegBal, nil
	}
	return maxNegBal - math.Abs(a.Balance), nil
}

// PATCH /transfers/{transferID}

func (t *transfer) Accept(j *types.Journal) (*types.Journal, error) {
	updated, err := pg.Journal.Accept(j)
	if err != nil {
		return nil, err
	}
	err = es.Journal.Update(updated)
	if err != nil {
		return nil, err
	}
	err = t.updateESEntityBalances(updated)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (t *transfer) updateESEntityBalances(j *types.Journal) error {
	from, err := pg.Account.FindByAccountNumber(j.FromAccountNumber)
	if err != nil {
		return err
	}
	err = es.Entity.UpdateBalance(from.AccountNumber, from.Balance)
	if err != nil {
		return err
	}
	to, err := pg.Account.FindByAccountNumber(j.ToAccountNumber)
	if err != nil {
		return err
	}
	err = es.Entity.UpdateBalance(to.AccountNumber, to.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (t *transfer) Cancel(transferID string, reason string) (*types.Journal, error) {
	canceled, err := pg.Journal.Cancel(transferID, reason)
	if err != nil {
		return nil, err
	}
	err = es.Journal.Update(canceled)
	if err != nil {
		return nil, err
	}
	return canceled, nil
}

// GET /user/entities

func (t *transfer) GetPendingTransfers(accountNumber string) ([]*types.TransferRespond, error) {
	journals, err := pg.Journal.GetPending(accountNumber)
	if err != nil {
		return nil, err
	}
	return types.NewJournalsToTransfersRespond(journals, accountNumber), nil
}

// GET /admin/transfers/{transferID}

func (t *transfer) AdminGetTransfer(transferID string) (*types.Journal, error) {
	journal, err := pg.Journal.FindByID(transferID)
	if err != nil {
		return nil, err
	}
	return journal, nil
}

// POST /admin/transfers

func (t *transfer) Create(req *types.AdminTransferReq) (*types.Journal, error) {
	created, err := pg.Journal.Create(req)
	if err != nil {
		return nil, err
	}
	err = es.Journal.Create(created)
	if err != nil {
		return nil, err
	}
	err = t.updateESEntityBalances(created)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// GET /admin/transfers

func (t *transfer) AdminSearch(req *types.AdminSearchTransferReq) (*types.AdminSearchTransferRespond, error) {
	result, err := es.Journal.AdminSearch(req)
	if err != nil {
		return nil, err
	}
	journals, err := pg.Journal.FindByIDs(result.IDs)
	if err != nil {
		return nil, err
	}
	return &types.AdminSearchTransferRespond{
		Transfers:       types.NewJournalsToAdminTransfersRespond(journals),
		NumberOfResults: result.NumberOfResults,
		TotalPages:      result.TotalPages,
	}, nil
}

// GET /admin/entities/{entityID}

func (t *transfer) AdminGetPendingTransfers(accountNumber string) ([]*types.AdminTransferRespond, error) {
	journals, err := pg.Journal.GetPending(accountNumber)
	if err != nil {
		return nil, err
	}
	return types.NewJournalsToAdminTransfersRespond(journals), nil
}
