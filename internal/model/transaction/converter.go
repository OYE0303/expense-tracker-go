package transaction

import (
	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/model/icon"
	"github.com/OYE0303/expense-tracker-go/internal/model/maincateg"
	"github.com/OYE0303/expense-tracker-go/internal/model/subcateg"
)

func cvtToDomainTransaction(t Transaction, m maincateg.MainCateg, s subcateg.SubCateg, i icon.Icon) domain.Transaction {
	return domain.Transaction{
		ID:     t.ID,
		UserID: t.UserID,
		Price:  t.Price,
		Note:   t.Note,
		Date:   t.Date,
		MainCateg: domain.MainCateg{
			ID:   m.ID,
			Name: m.Name,
			Type: domain.CvtToMainCategType(m.Type),
			Icon: domain.Icon{
				ID:  i.ID,
				URL: i.URL,
			},
		},
		SubCateg: domain.SubCateg{
			ID:          s.ID,
			Name:        s.Name,
			MainCategID: m.ID, // use m.ID because in the get query, we don't reterive the subCateg.MainCategID
		},
	}
}

func cvtToModelTransaction(t *domain.Transaction) *Transaction {
	return &Transaction{
		UserID:      t.UserID,
		MainCategID: t.MainCateg.ID,
		SubCategID:  t.SubCateg.ID,
		Price:       t.Price,
		Note:        t.Note,
		Date:        t.Date,
	}
}
