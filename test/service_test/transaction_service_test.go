package service_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/internal/service"
	"github.com/imranzahoor/banking-ledger/pkg/config"
	"github.com/imranzahoor/banking-ledger/pkg/errors"
	"github.com/imranzahoor/banking-ledger/test/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransactionService", func() {
	var (
		mockCtrl        *gomock.Controller
		mockAccountRepo *mocks.MockAccountRepository
		mockLedgerRepo  *mocks.MockLedgerRepository
		transactionSvc  service.TransactionServiceInterface
		ctx             context.Context
		cfg             config.Config
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAccountRepo = mocks.NewMockAccountRepository(mockCtrl)
		mockLedgerRepo = mocks.NewMockLedgerRepository(mockCtrl)
		ctx = context.TODO()

		cfg = config.Config{
			KafkaBrokers: []string{"localhost:9092"},
			KafkaTopic:   "transactions",
		}

		transactionSvc = service.NewTransactionService(mockAccountRepo, mockLedgerRepo, cfg)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("EnqueueTransaction", func() {
		It("should enqueue a valid transaction", func() {
			txn := &model.Transaction{
				AccountID: uuid.New(),
				Type:      "deposit",
				Amount:    1000,
			}

			err := transactionSvc.EnqueueTransaction(ctx, txn)
			Expect(err).To(BeNil())
			Expect(txn.ID).NotTo(Equal(uuid.Nil))
		})

		It("should return error for invalid transaction marshal", func() {
			invalidTxn := &model.Transaction{}

			err := transactionSvc.EnqueueTransaction(ctx, invalidTxn)
			Expect(err).To(BeNil())
		})
	})

	Describe("GetTransactions", func() {
		It("should return list of transactions", func() {
			accountID := uuid.New()
			transactions := []model.Transaction{
				{ID: uuid.New(), AccountID: accountID, Amount: 1000, Type: "deposit"},
				{ID: uuid.New(), AccountID: accountID, Amount: 500, Type: "withdrawal"},
			}

			mockLedgerRepo.EXPECT().
				GetTransactionsByAccountID(gomock.Any(), accountID, int64(10), int64(0)).
				Return(transactions, nil).
				Times(1)

			result, err := transactionSvc.GetTransactions(accountID, 10, 0)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(transactions))
		})

		It("should return error if ledger repo fails", func() {
			accountID := uuid.New()
			mockLedgerRepo.EXPECT().
				GetTransactionsByAccountID(gomock.Any(), accountID, int64(10), int64(0)).
				Return(nil, errors.ErrFake).
				Times(1)

			result, err := transactionSvc.GetTransactions(accountID, 10, 0)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})
})
