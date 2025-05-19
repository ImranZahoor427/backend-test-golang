package service_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/internal/service"
	"github.com/imranzahoor/banking-ledger/test/mocks"
)

var _ = Describe("AccountService", func() {
	var (
		mockCtrl      *gomock.Controller
		mockRepo      *mocks.MockAccountRepository
		accountSvc    *service.AccountService
		ctx           context.Context
		sampleAccount *model.Account
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockAccountRepository(mockCtrl)
		accountSvc = service.NewAccountService(mockRepo)
		ctx = context.TODO()

		sampleAccount = &model.Account{
			ID:        uuid.New(),
			OwnerName: "Alice",
			Balance:   0,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should create an account successfully", func() {
		mockRepo.EXPECT().
			CreateAccount(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		_, err := accountSvc.CreateAccount(ctx, sampleAccount.OwnerName, sampleAccount.Balance)
		Expect(err).To(BeNil())
	})

	It("should fetch an account by ID", func() {
		mockRepo.EXPECT().
			GetAccountByID(gomock.Any(), sampleAccount.ID.String()).
			Return(sampleAccount, nil).
			Times(1)

		acc, err := accountSvc.GetAccountByID(ctx, sampleAccount.ID.String())
		Expect(err).To(BeNil())
		Expect(acc).To(Equal(sampleAccount))
	})
})
