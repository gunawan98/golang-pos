package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type PurchaseRepository interface {
	AddPurchase(ctx context.Context, tx *sql.Tx, purchase domain.Purchase) domain.Purchase
}