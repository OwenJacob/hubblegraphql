package graph

//go:generate go run github.com/99designs/gqlgen
import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"github.com/owenjacob/hubblegraphql/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *gorm.DB
}

// Limit when search for unconstrained data
var maxSearchLimit int = 500

// Network constrained number of transactions per ledger
var maxTxnsLimit int = 105

// Network constrained number of operations per transaction
var maxOpsLimit int = 100

func parseAccountFlags(flags int) *model.Flags {
	switch flags {
	case 1:
		return &model.Flags{
			AuthImmutable: false,
			AuthRequired:  true,
			AuthRevocable: false,
		}
	case 3:
		return &model.Flags{
			AuthImmutable: false,
			AuthRequired:  true,
			AuthRevocable: true,
		}
	default:
		return &model.Flags{
			AuthImmutable: false,
			AuthRequired:  false,
			AuthRevocable: false,
		}
	}
}

type Account struct {
	AccountID          string `gorm:"column:account_id"`
	Balance            int64  `gorm:"column:balance"`
	BuyingLiabilities  int64  `gorm:"column:buying_liabilities"`
	SellingLiabilities int64  `gorm:"column:selling_liabilities"`
	SeqNum             int64  `gorm:"column:sequence_number"`
	Subentries         int    `gorm:"column:num_subentries"`
	InflationDest      string `gorm:"column:inflation_destination"`
	Flags              int    `gorm:"column:flags"`
	HomeDomain         string `gorm:"column:home_domain"`
	MasterWeight       int    `gorm:"column:master_weight"`
	Low                int    `gorm:"column:threshold_low"`
	Medium             int    `gorm:"column:threshold_medium"`
	High               int    `gorm:"column:threshold_high"`
	LastModified       int    `gorm:"column:last_modified_ledger"`
}

type AccountData struct {
	LedgerKey    string `gorm:"column:ledger_key"`
	AccountID    string `gorm:"column:account_id"`
	Name         string `gorm:"column:name"`
	Value        string `gorm:"column:value"`
	LastModified int    `gorm:"column:last_modified_ledger"`
}

type AccountSigner struct {
	AccountID string `gorm:"column:account_id"`
	Signer    string `gorm:"column:signer"`
	Weight    int    `gorm:"column:weight"`
}

type HistoryAccounts struct {
	ID      int64  `gorm:"column:id"`
	Address string `gorm:"column:address"`
}

type HistoryLedgers struct {
	Sequence                   int       `gorm:"column:sequence"`
	LedgerHash                 string    `gorm:"column:ledger_hash"`
	PreviousLedgerHash         string    `gorm:"column:previous_ledger_hash"`
	TransactionCount           int       `gorm:"column:transaction_count"`
	OperationCount             int       `gorm:"column: operation_count"`
	ClosedAt                   time.Time `gorm:"column:closed_at"`
	CreatedAt                  time.Time `gorm:"column:created_at"`
	UpdatedAt                  time.Time `gorm:"column:updated_at"`
	ID                         int64     `gorm:"column:id"`
	ImporterVersion            int       `gorm:"column:importer_version"`
	TotalCoins                 int64     `gorm:"column:total_coins"`
	FeePool                    int64     `gorm:"column:fee_pool"`
	BaseFee                    int       `gorm:"column:base_fee"`
	BaseReserve                int       `gorm:"column:base_reserve"`
	MaxTxSetSize               int       `gorm:"column:max_tx_set_size"`
	ProtocolVersion            int       `gorm:"column:protocol_version"`
	LedgerHeader               string    `gorm:"column:ledger_header"`
	SuccessfulTransactionCount int       `gorm:"column:successful_transaction_count"`
	FailedTransactionCount     int       `gorm:"column:failed_transaction_count"`
}

type HistoryOperationParticipants struct {
	HistoryOperationID int64 `gorm:"column:history_operation_id"`
	HistoryAccountID   int64 `gorm:"column:history_account_id"`
}

type HistoryOperations struct {
	ID               int64          `gorm:"column:id"`
	TransactionID    int64          `gorm:"column:transaction_id"`
	ApplicationOrder int            `gorm:"column:application_order"`
	Type             int            `gorm:"column:type"`
	Details          postgres.Jsonb `gorm:"column:details"`
	SourceAccount    string         `gorm:"column:source_account"`
}

type HistoryTransactionParticipants struct {
	HistoryTransactionID int64 `gorm:"column:history_transaction_id"`
	HistoryAccountID     int64 `gorm:"column:history_account_id"`
}

type HistoryTransactions struct {
	TransactionHash      string         `gorm:"column:transaction_hash"`
	LedgerSequence       int            `gorm:"column:ledger_sequence"`
	ApplicationOrder     int            `gorm:"column:application_order"`
	Account              string         `gorm:"column:account"`
	AccountSequence      int64          `gorm:"column:account_sequence"`
	MaxFee               int64          `gorm:"column:max_fee"`
	OperationCount       int            `gorm:"column:operation_count"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
	ID                   int64          `gorm:"column:id"`
	TxEnvelope           string         `gorm:"column:tx_envelope"`
	TxResult             string         `gorm:"column:tx_result"`
	TxMeta               string         `gorm:"column:tx_meta"`
	TxFeeMeta            string         `gorm:"column:tx_fee_meta"`
	Signatures           pq.StringArray `gorm:"column:signatures"`
	MemoType             string         `gorm:"column:memo_type"`
	Memo                 string         `gorm:"column:memo"`
	TimeBounds           []uint8        `gorm:"column:time_bounds"`
	Successful           bool           `gorm:"column:successful"`
	FeeCharged           int64          `gorm:"column:fee_charged"`
	InnerTransactionHash string         `gorm:"column:inner_transaction_hash"`
	FeeAccount           string         `gorm:"column:fee_account"`
	InnerSignatures      pq.StringArray `gorm:"column:inner_signatures"`
	NewMaxFee            int64          `gorm:"column:new_max_fee"`
}

type TrustLine struct {
	LedgerKey          string `gorm:"column:ledger_key"`
	AccountID          string `gorm:"column:account_id"`
	AssetType          int    `gorm:"column:asset_type"`
	AssetIssuer        string `gorm:"column:asset_issuer"`
	AssetCode          string `gorm:"column:asset_code"`
	Balance            int64  `gorm:"column:balance"`
	TrustLineLimit     int64  `gorm:"column:trust_line_limit"`
	BuyingLiabilities  int64  `gorm:"column:buying_liabilities"`
	SellingLiabilities int64  `gorm:"column:selling_liabilities"`
	Flags              int    `gorm:"column:flags"`
	LastModifiedLedger int    `gorm:"column:last_modified_ledger"`
}
