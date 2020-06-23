package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/owenjacob/hubblegraphql/graph/generated"
	"github.com/owenjacob/hubblegraphql/graph/model"
)

func (r *accountResolver) Balances(ctx context.Context, obj *model.Account) ([]*model.Balance, error) {
	accountBalances := []TrustLine{}
	err := r.DB.Table("trust_lines").Where("account_id = ?", obj.ID).Order("balance desc").Find(&accountBalances).Error
	if err != nil {
		return nil, err
	}

	if len(accountBalances) != 0 {
		balances := make([]*model.Balance, 0, len(accountBalances))
		for i := range accountBalances {
			authorized := false
			if accountBalances[i].Flags == 1 {
				authorized = true
			}
			balances = append(balances, &model.Balance{
				Balance:            strconv.FormatFloat((float64(accountBalances[i].Balance) / 10000000.0), byte('f'), 7, 64),
				BuyingLiabilities:  strconv.FormatFloat((float64(accountBalances[i].BuyingLiabilities) / 10000000.0), byte('f'), 7, 64),
				SellingLiabilities: strconv.FormatFloat((float64(accountBalances[i].SellingLiabilities) / 10000000.0), byte('f'), 7, 64),
				Limit:              strconv.FormatFloat((float64(accountBalances[i].TrustLineLimit) / 10000000.0), byte('f'), 7, 64),
				LastModifiedLedger: accountBalances[i].LastModifiedLedger,
				IsAuthorized:       authorized,
				AssetCode:          accountBalances[i].AssetCode,
				AssetIssuer:        accountBalances[i].AssetIssuer,
			})
		}
		return balances, nil
	}
	return nil, nil
}

func (r *accountResolver) Signers(ctx context.Context, obj *model.Account) ([]*model.Signer, error) {
	accountSigners := []AccountSigner{}
	err := r.DB.Table("accounts_signers").Where("account_id = ?", obj.ID).Order("weight asc").Find(&accountSigners).Error
	if err != nil {
		return nil, err
	}

	signers := make([]*model.Signer, 0, len(accountSigners))
	for i := range accountSigners {
		signers = append(signers, &model.Signer{
			Weight: accountSigners[i].Weight,
			Key:    accountSigners[i].Signer,
		})
	}
	return signers, nil
}

func (r *accountResolver) Data(ctx context.Context, obj *model.Account, name *string) ([]*model.Data, error) {
	accountData := []AccountData{}
	if name == nil {
		// If no argument is passed return all data
		err := r.DB.Table("accounts_data").Where("account_id = ?", obj.ID).Order("name").Find(&accountData).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := r.DB.Table("accounts_data").Where("account_id =? AND name = ?", obj.ID, name).Order("name").Find(&accountData).Error
		if err != nil {
			return nil, err
		}
	}

	if len(accountData) != 0 {
		data := make([]*model.Data, 0, len(accountData))
		for i := range accountData {
			b64Val, _ := base64.StdEncoding.DecodeString(accountData[i].Value)
			value := string(b64Val)
			data = append(data, &model.Data{
				Name:  &accountData[i].Name,
				Value: &value,
			})
		}
		return data, nil
	}
	return nil, nil
}

func (r *accountResolver) Transactions(ctx context.Context, obj *model.Account, limit *int, order *model.Order, filterBy *model.FilterBy) ([]*model.Transaction, error) {
	if *limit > maxSearchLimit {
		return nil, fmt.Errorf("Maximum limit is %d", maxSearchLimit)
	}

	// Concat order string
	IDorder := "history_transaction_participants.history_transaction_id " + order.String()

	historyTransactionParticipants := []HistoryTransactionParticipants{}

	if filterBy != nil {
		// If filtering options have been defined check them here
		if filterBy.Account != nil {
			// Pull DB by accounts
			return nil, errors.New("pending implementation")
		} else if filterBy.Ledger != nil && filterBy.Date != nil {
			// Cannot filter by ledger and date at the same time
			return nil, errors.New("Cannot filter by ledger and date")
		} else if filterBy.Ledger != nil {
			// Pull DB between ledgers
			return nil, errors.New("pending implementation")
		} else if filterBy.Date != nil {
			// Pull DB between dates
			return nil, errors.New("pending implementation")
		} else {
			// Something I didn't think of
			return nil, errors.New("Unexpected condition")
		}
	} else {
		// If no filtering arguments are specified return transactions according to the limit and order
		// SELECT history_transaction_participants.history_transaction_id FROM history_accounts INNER JOIN history_transaction_participants ON history_accounts.id = history_transaction_participants.history_account_id WHERE history_accounts.address = 'GBJTZPKVINL5MNHYLJA7FYQ4PECRFS74AHFCKYVBI7NUCFQUERBN2Y7N';
		err := r.DB.Table("history_accounts").Select("history_transaction_participants.history_transaction_id").Joins("INNER JOIN history_transaction_participants ON history_accounts.id = history_transaction_participants.history_account_id").Where("history_accounts.address = ?", obj.ID).Order(IDorder).Limit(*limit).Find(&historyTransactionParticipants).Error
		if err != nil {
			return nil, err
		}
	}

	if len(historyTransactionParticipants) != 0 {
		transactions := make([]*model.Transaction, 0, len(historyTransactionParticipants))

		for i := range historyTransactionParticipants {
			transaction := HistoryTransactions{}
			err := r.DB.Table("history_transactions").Where("id = ?", historyTransactionParticipants[i].HistoryTransactionID).First(&transaction).Error
			if err != nil {
				return nil, err
			}

			// Parse for pointers
			transactionCreatedAt := transaction.CreatedAt.Format("2006-01-02 15:04:05")
			transactionUpdatedAt := transaction.UpdatedAt.Format("2006-01-02 15:04:05")
			transactionID := strconv.FormatInt(transaction.ID, 10)
			transactionFeeCharged := strconv.FormatInt(transaction.FeeCharged, 10)
			transactionNewMaxFee := strconv.FormatInt(transaction.NewMaxFee, 10)
			transactionTimebounds := make([]*int, 0, len(transaction.TimeBounds))
			for j := range transaction.TimeBounds {
				timebound := int(transaction.TimeBounds[j])
				transactionTimebounds = append(transactionTimebounds, &timebound)
			}
			transactionInnerSignatures := make([]*string, 0, len(transaction.InnerSignatures))
			for j := range transaction.InnerSignatures {
				transactionInnerSignatures = append(transactionInnerSignatures, &transaction.InnerSignatures[j])
			}

			transactions = append(transactions, &model.Transaction{
				TransactionHash:      transaction.TransactionHash,
				LedgerSequence:       transaction.LedgerSequence,
				ApplicationOrder:     transaction.ApplicationOrder,
				Account:              transaction.Account,
				AccountSequence:      strconv.FormatInt(transaction.AccountSequence, 10),
				MaxFee:               strconv.FormatFloat((float64(transaction.MaxFee) / 10000000.0), byte('f'), 7, 64),
				OperationCount:       transaction.OperationCount,
				CreatedAt:            &transactionCreatedAt,
				UpdatedAt:            &transactionUpdatedAt,
				ID:                   &transactionID,
				TxEnvelope:           transaction.TxEnvelope,
				TxResult:             transaction.TxResult,
				TxMeta:               transaction.TxMeta,
				TxFeeMeta:            transaction.TxFeeMeta,
				Signatures:           transaction.Signatures,
				MemoType:             transaction.MemoType,
				Memo:                 &transaction.Memo,
				TimeBounds:           transactionTimebounds,
				Successful:           &transaction.Successful,
				FeeCharged:           &transactionFeeCharged,
				InnerTransactionHash: &transaction.InnerTransactionHash,
				FeeAccount:           &transaction.FeeAccount,
				InnerSignatures:      transactionInnerSignatures,
				NewMaxFee:            &transactionNewMaxFee,
			})
		}
		return transactions, nil
	}
	return nil, nil
}

func (r *ledgerResolver) Transactions(ctx context.Context, obj *model.Ledger, limit *int, order *model.Order) ([]*model.Transaction, error) {
	if *limit > maxTxnsLimit {
		return nil, fmt.Errorf("ledgers cannot contain more than %d transactions", maxTxnsLimit)
	}

	// Concat order string
	IDorder := "id " + order.String()

	ledgerTransactions := []HistoryTransactions{}
	err := r.DB.Table("history_transactions").Where("ledger_sequence = ?", obj.Sequence).Order(IDorder).Limit(*limit).Find(&ledgerTransactions).Error
	if err != nil {
		return nil, err
	}

	if len(ledgerTransactions) != 0 {
		transactions := make([]*model.Transaction, 0, len(ledgerTransactions))

		for i := range ledgerTransactions {
			// Parse for pointers
			transactionCreatedAt := ledgerTransactions[i].CreatedAt.Format("2006-01-02 15:04:05")
			transactionUpdatedAt := ledgerTransactions[i].UpdatedAt.Format("2006-01-02 15:04:05")
			transactionID := strconv.FormatInt(ledgerTransactions[i].ID, 10)
			transactionFeeCharged := strconv.FormatInt(ledgerTransactions[i].FeeCharged, 10)
			transactionNewMaxFee := strconv.FormatInt(ledgerTransactions[i].NewMaxFee, 10)
			transactionTimebounds := make([]*int, 0, len(ledgerTransactions[i].TimeBounds))
			for j := range ledgerTransactions[i].TimeBounds {
				timebound := int(ledgerTransactions[i].TimeBounds[j])
				transactionTimebounds = append(transactionTimebounds, &timebound)
			}
			transactionInnerSignatures := make([]*string, 0, len(ledgerTransactions[i].InnerSignatures))
			for j := range ledgerTransactions[i].InnerSignatures {
				transactionInnerSignatures = append(transactionInnerSignatures, &ledgerTransactions[i].InnerSignatures[j])
			}

			transactions = append(transactions, &model.Transaction{
				TransactionHash:      ledgerTransactions[i].TransactionHash,
				LedgerSequence:       ledgerTransactions[i].LedgerSequence,
				ApplicationOrder:     ledgerTransactions[i].ApplicationOrder,
				Account:              ledgerTransactions[i].Account,
				AccountSequence:      strconv.FormatInt(ledgerTransactions[i].AccountSequence, 10),
				MaxFee:               strconv.FormatFloat((float64(ledgerTransactions[i].MaxFee) / 10000000.0), byte('f'), 7, 64),
				OperationCount:       ledgerTransactions[i].OperationCount,
				CreatedAt:            &transactionCreatedAt,
				UpdatedAt:            &transactionUpdatedAt,
				ID:                   &transactionID,
				TxEnvelope:           ledgerTransactions[i].TxEnvelope,
				TxResult:             ledgerTransactions[i].TxResult,
				TxMeta:               ledgerTransactions[i].TxMeta,
				TxFeeMeta:            ledgerTransactions[i].TxFeeMeta,
				Signatures:           ledgerTransactions[i].Signatures,
				MemoType:             ledgerTransactions[i].MemoType,
				Memo:                 &ledgerTransactions[i].Memo,
				TimeBounds:           transactionTimebounds,
				Successful:           &ledgerTransactions[i].Successful,
				FeeCharged:           &transactionFeeCharged,
				InnerTransactionHash: &ledgerTransactions[i].InnerTransactionHash,
				FeeAccount:           &ledgerTransactions[i].FeeAccount,
				InnerSignatures:      transactionInnerSignatures,
				NewMaxFee:            &transactionNewMaxFee,
			})
		}
		return transactions, nil
	}
	return nil, nil
}

func (r *queryResolver) Account(ctx context.Context, pubKey string) (*model.Account, error) {
	account := Account{}
	notFound := r.DB.Table("accounts").Where("account_id = ?", pubKey).First(&account).RecordNotFound()
	if notFound {
		return nil, errors.New("Account not found")
	}

	return &model.Account{
		ID:              account.AccountID,
		Sequence:        strconv.FormatInt(account.SeqNum, 10),
		HomeDomain:      account.HomeDomain,
		NativeBalance:   strconv.FormatFloat((float64(account.Balance) / 10000000.0), byte('f'), 7, 64),
		MasterWeight:    account.MasterWeight,
		LowThreshold:    account.Low,
		MediumThreshold: account.Medium,
		HighThreshold:   account.High,
		Flags:           parseAccountFlags(account.Flags),
	}, nil
}

func (r *queryResolver) Transaction(ctx context.Context, hash string) (*model.Transaction, error) {
	transaction := HistoryTransactions{}
	notFound := r.DB.Table("history_transactions").Where("transaction_hash = ?", hash).First(&transaction).RecordNotFound()
	if notFound {
		return nil, errors.New("Transaction not found")
	}

	// Parse for pointers
	transactionCreatedAt := transaction.CreatedAt.Format("2006-01-02 15:04:05")
	transactionUpdatedAt := transaction.UpdatedAt.Format("2006-01-02 15:04:05")
	transactionID := strconv.FormatInt(transaction.ID, 10)
	transactionFeeCharged := strconv.FormatInt(transaction.FeeCharged, 10)
	transactionNewMaxFee := strconv.FormatInt(transaction.NewMaxFee, 10)
	transactionTimebounds := make([]*int, 0, len(transaction.TimeBounds))
	for i := range transaction.TimeBounds {
		timebound := int(transaction.TimeBounds[i])
		transactionTimebounds = append(transactionTimebounds, &timebound)
	}
	transactionInnerSignatures := make([]*string, 0, len(transaction.InnerSignatures))
	for i := range transaction.InnerSignatures {
		transactionInnerSignatures = append(transactionInnerSignatures, &transaction.InnerSignatures[i])
	}

	return &model.Transaction{
		TransactionHash:      transaction.TransactionHash,
		LedgerSequence:       transaction.LedgerSequence,
		ApplicationOrder:     transaction.ApplicationOrder,
		Account:              transaction.Account,
		AccountSequence:      strconv.FormatInt(transaction.AccountSequence, 10),
		MaxFee:               strconv.FormatFloat((float64(transaction.MaxFee) / 10000000.0), byte('f'), 7, 64),
		OperationCount:       transaction.OperationCount,
		CreatedAt:            &transactionCreatedAt,
		UpdatedAt:            &transactionUpdatedAt,
		ID:                   &transactionID,
		TxEnvelope:           transaction.TxEnvelope,
		TxResult:             transaction.TxResult,
		TxMeta:               transaction.TxMeta,
		TxFeeMeta:            transaction.TxFeeMeta,
		Signatures:           transaction.Signatures,
		MemoType:             transaction.MemoType,
		Memo:                 &transaction.Memo,
		TimeBounds:           transactionTimebounds,
		Successful:           &transaction.Successful,
		FeeCharged:           &transactionFeeCharged,
		InnerTransactionHash: &transaction.InnerTransactionHash,
		FeeAccount:           &transaction.FeeAccount,
		InnerSignatures:      transactionInnerSignatures,
		NewMaxFee:            &transactionNewMaxFee,
	}, nil
}

func (r *queryResolver) Ledger(ctx context.Context, number *int, limit *int, order *model.Order, filterBy *model.FilterBy) ([]*model.Ledger, error) {
	if *limit > maxSearchLimit {
		return nil, fmt.Errorf("Maximum limit is %d", maxSearchLimit)
	}

	// Concat order string
	IDorder := "id " + order.String()

	ledger := []HistoryLedgers{}

	if number != nil {
		// If a specific ledger is requested ignore everything else
		err := r.DB.Table("history_ledgers").Where("sequence = ?", number).First(&ledger).Error
		if err != nil {
			return nil, errors.New("Ledger not found")
		}
	} else if limit != nil && filterBy == nil {
		// If a specific limit is set and there are no filtering options pull the latest ledgers
		err := r.DB.Table("history_ledgers").Order(IDorder).Limit(*limit).Find(&ledger).Error
		if err != nil {
			return nil, err
		}
	} else if filterBy != nil {
		// If filtering options have been defined check them here
		if filterBy.Account != nil {
			// Cannot filter by accounts (maybe via super complex query later?)
			return nil, errors.New("Cannot filter ledgers by accounts")
		} else if filterBy.Ledger != nil && filterBy.Date != nil {
			// Cannot filter by ledger and date at the same time
			return nil, errors.New("Cannot filter by ledger and date")
		} else if filterBy.Ledger != nil {
			// Pull DB between ledgers
			return nil, errors.New("pending implementation")
		} else if filterBy.Date != nil {
			// Pull DB between dates
			return nil, errors.New("pending implementation")
		} else {
			// Something I didn't think of
			return nil, errors.New("Unexpected condition")
		}
	} else {
		// If no arguments are specified just return the latest ledger
		err := r.DB.Table("history_ledgers").Last(&ledger).Error
		if err != nil {
			return nil, err
		}
	}

	if len(ledger) != 0 {
		ledgers := make([]*model.Ledger, 0, len(ledger))

		for i := range ledger {
			ledgerCreatedAt := ledger[i].CreatedAt.Format("2006-01-02 15:04:05")
			ledgerUpdatedAt := ledger[i].UpdatedAt.Format("2006-01-02 15:04:05")
			ledgerID := int(ledger[i].ID)

			ledgers = append(ledgers, &model.Ledger{
				Sequence:                   ledger[i].Sequence,
				LedgerHash:                 ledger[i].LedgerHash,
				PreviousLedgerHash:         &ledger[i].PreviousLedgerHash,
				TransactionCount:           ledger[i].TransactionCount,
				OperationCount:             ledger[i].OperationCount,
				ClosedAt:                   ledger[i].ClosedAt.Format("2006-01-02 15:04:05"),
				CreatedAt:                  &ledgerCreatedAt,
				UpdatedAt:                  &ledgerUpdatedAt,
				ID:                         &ledgerID,
				ImporterVersion:            ledger[i].ImporterVersion,
				TotalCoins:                 strconv.FormatFloat((float64(ledger[i].TotalCoins) / 10000000.0), byte('f'), 7, 64),
				FeePool:                    strconv.FormatFloat((float64(ledger[i].FeePool) / 10000000.0), byte('f'), 7, 64),
				BaseFee:                    ledger[i].BaseFee,
				BaseReserve:                ledger[i].BaseReserve,
				MaxTxSetSize:               ledger[i].MaxTxSetSize,
				ProtocolVersion:            ledger[i].ProtocolVersion,
				LedgerHeader:               &ledger[i].LedgerHeader,
				SuccessfulTransactionCount: &ledger[i].SuccessfulTransactionCount,
				FailedTransactionCount:     &ledger[i].FailedTransactionCount,
			})
		}
		return ledgers, nil
	}
	return nil, errors.New("Query failed")
}

func (r *transactionResolver) Operations(ctx context.Context, obj *model.Transaction, limit *int, order *model.Order) ([]*model.Operation, error) {
	// Check max limit
	if *limit > maxOpsLimit {
		return nil, fmt.Errorf("transactions cannot contain more than %d operations", maxOpsLimit)
	}

	// Concat order string
	IDorder := "id " + order.String()

	transactionOperations := []HistoryOperations{}
	err := r.DB.Table("history_operations").Where("transaction_id = ?", obj.ID).Order(IDorder).Limit(*limit).Find(&transactionOperations).Error
	if err != nil {
		return nil, err
	}

	operations := make([]*model.Operation, 0, len(transactionOperations))
	for i := range transactionOperations {
		detail, _ := transactionOperations[i].Details.MarshalJSON()
		details := string(detail)
		operations = append(operations, &model.Operation{
			ID:               strconv.FormatInt(transactionOperations[i].ID, 10),
			TransactionID:    strconv.FormatInt(transactionOperations[i].TransactionID, 10),
			ApplicationOrder: transactionOperations[i].ApplicationOrder,
			Type:             transactionOperations[i].Type,
			Details:          &details,
			SourceAccount:    transactionOperations[i].SourceAccount,
		})
	}
	return operations, nil
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

// Ledger returns generated.LedgerResolver implementation.
func (r *Resolver) Ledger() generated.LedgerResolver { return &ledgerResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Transaction returns generated.TransactionResolver implementation.
func (r *Resolver) Transaction() generated.TransactionResolver { return &transactionResolver{r} }

type accountResolver struct{ *Resolver }
type ledgerResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type transactionResolver struct{ *Resolver }
