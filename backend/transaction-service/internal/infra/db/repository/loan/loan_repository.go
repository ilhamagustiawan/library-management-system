package loan

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

const loanReturnedRoutingKey = "transactions.loan.returned.v1"

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

func (r *Repository) CreatePending(ctx context.Context, loan *entity.Loan) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO member_loan_counters (member_id, active_count, updated_at)
		VALUES (?, 0, ?) ON DUPLICATE KEY UPDATE member_id = VALUES(member_id)`, loan.MemberID, loan.CreatedAt); err != nil {
		return err
	}
	var activeCount int
	if err := tx.GetContext(ctx, &activeCount, `SELECT active_count FROM member_loan_counters WHERE member_id = ? FOR UPDATE`, loan.MemberID); err != nil {
		return err
	}
	if activeCount >= 3 {
		return errs.ErrLoanLimit
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO loans (id, member_id, book_id, status, stock_sync_status, borrowed_at, due_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		loan.ID, loan.MemberID, loan.BookID, loan.Status, loan.StockSyncStatus,
		loan.BorrowedAt, loan.DueAt, loan.CreatedAt, loan.UpdatedAt)
	if duplicateKey(err) {
		return errs.ErrActiveLoan
	}
	if err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `UPDATE member_loan_counters SET active_count = active_count + 1, updated_at = ? WHERE member_id = ? AND active_count < 3`, loan.UpdatedAt, loan.MemberID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected != 1 {
		return errs.ErrLoanLimit
	}
	return tx.Commit()
}

func (r *Repository) Activate(ctx context.Context, loanID, transactionID string, now time.Time) (*entity.Loan, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE loans SET status = 'active', updated_at = ? WHERE id = ? AND status = 'pending_reservation'`, now, loanID)
	if err != nil {
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected != 1 {
		return nil, errs.ErrNotFound
	}
	var loan entity.Loan
	if err := tx.GetContext(ctx, &loan, `SELECT id, member_id, book_id, status, stock_sync_status, borrowed_at, due_at, returned_at, created_at, updated_at FROM loans WHERE id = ?`, loanID); err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO loan_transactions (id, loan_id, member_id, book_id, type, occurred_at) VALUES (?, ?, ?, ?, 'borrow', ?)`, transactionID, loan.ID, loan.MemberID, loan.BookID, now); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *Repository) CancelPending(ctx context.Context, loanID string, now time.Time) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var loan entity.Loan
	if err := tx.GetContext(ctx, &loan, `SELECT id, member_id, status FROM loans WHERE id = ? FOR UPDATE`, loanID); errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNotFound
	} else if err != nil {
		return err
	}
	if loan.Status == entity.LoanCancelled {
		return tx.Commit()
	}
	if loan.Status != entity.LoanPendingReservation {
		return fmt.Errorf("cannot cancel loan in %s state", loan.Status)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE loans SET status = 'cancelled', updated_at = ? WHERE id = ?`, now, loanID); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `UPDATE member_loan_counters SET active_count = active_count - 1, updated_at = ? WHERE member_id = ? AND active_count > 0`, now, loan.MemberID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected != 1 {
		return fmt.Errorf("member active-loan counter is inconsistent")
	}
	return tx.Commit()
}

func (r *Repository) Return(ctx context.Context, command repository.ReturnCommand) (*entity.Loan, bool, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()
	loan, err := getLoan(ctx, tx, command.LoanID, true)
	if err != nil {
		return nil, false, err
	}
	if !command.AllowAnyMember && loan.MemberID != command.MemberID {
		return nil, false, errs.ErrForbidden
	}
	if loan.Status == entity.LoanReturned {
		if err := tx.Commit(); err != nil {
			return nil, false, err
		}
		return loan, true, nil
	}
	if loan.Status != entity.LoanActive {
		return nil, false, errs.ErrNotFound
	}
	overdue := calculateOverdueDays(loan.DueAt, command.ReturnedAt)
	if overdue > 0 {
		loan.Fine = &entity.Fine{
			ID: command.FineID, LoanID: loan.ID, MemberID: loan.MemberID, OverdueDays: overdue,
			DailyRateMinor: command.DailyFineMinor, TotalAmountMinor: int64(overdue) * command.DailyFineMinor,
			Currency: "IDR", Status: entity.FineUnpaid, AssessedAt: command.ReturnedAt,
		}
		if _, err := tx.NamedExecContext(ctx, `
			INSERT INTO fines (id, loan_id, member_id, overdue_days, daily_rate_minor, total_amount_minor, currency, status, assessed_at)
			VALUES (:fine_id, :loan_id, :fine_member_id, :overdue_days, :daily_rate_minor, :total_amount_minor, :currency, :fine_status, :assessed_at)`, loan.Fine); err != nil {
			return nil, false, err
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE loans SET status = 'returned', stock_sync_status = 'pending', returned_at = ?, updated_at = ? WHERE id = ?`, command.ReturnedAt, command.ReturnedAt, loan.ID); err != nil {
		return nil, false, err
	}
	result, err := tx.ExecContext(ctx, `UPDATE member_loan_counters SET active_count = active_count - 1, updated_at = ? WHERE member_id = ? AND active_count > 0`, command.ReturnedAt, loan.MemberID)
	if err != nil {
		return nil, false, err
	}
	if affected, _ := result.RowsAffected(); affected != 1 {
		return nil, false, fmt.Errorf("member active-loan counter is inconsistent")
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO loan_transactions (id, loan_id, member_id, book_id, type, occurred_at) VALUES (?, ?, ?, ?, 'return', ?)`, command.TransactionID, loan.ID, loan.MemberID, loan.BookID, command.ReturnedAt); err != nil {
		return nil, false, err
	}
	event := entity.Event[entity.LoanReturnedData]{EventID: command.EventID, Type: entity.LoanReturnedEventType, OccurredAt: command.ReturnedAt, Data: entity.LoanReturnedData{LoanID: loan.ID, BookID: loan.BookID, MemberID: loan.MemberID, ReturnedAt: command.ReturnedAt}}
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, false, fmt.Errorf("marshal loan-returned event: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO outbox_events (id, event_type, routing_key, payload, next_attempt_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`, command.EventID, event.Type, loanReturnedRoutingKey, payload, command.ReturnedAt, command.ReturnedAt); err != nil {
		return nil, false, err
	}
	if err := tx.Commit(); err != nil {
		return nil, false, err
	}
	loan.Status, loan.StockSyncStatus = entity.LoanReturned, entity.StockSyncPending
	loan.ReturnedAt, loan.UpdatedAt = &command.ReturnedAt, command.ReturnedAt
	return loan, false, nil
}

func (r *Repository) ListTransactions(ctx context.Context, filter repository.PageFilter) (repository.TransactionPage, error) {
	where, args := "", []any{}
	if filter.MemberID != "" {
		where, args = " WHERE t.member_id = ?", append(args, filter.MemberID)
	}
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM loan_transactions t`+where, args...); err != nil {
		return repository.TransactionPage{}, err
	}
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	rows, err := r.db.QueryxContext(ctx, `
		SELECT t.id, t.loan_id, t.member_id, t.book_id, t.type, t.occurred_at,
		       f.id, f.overdue_days, f.daily_rate_minor, f.total_amount_minor, f.currency, f.status, f.assessed_at
		FROM loan_transactions t LEFT JOIN fines f ON f.loan_id = t.loan_id AND t.type = 'return'`+where+`
		ORDER BY t.occurred_at DESC, t.id DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return repository.TransactionPage{}, err
	}
	defer rows.Close()
	items := make([]*entity.Transaction, 0)
	for rows.Next() {
		var item entity.Transaction
		var fineID, currency, status sql.NullString
		var overdue sql.NullInt64
		var dailyRate, totalAmount sql.NullInt64
		var assessed sql.NullTime
		if err := rows.Scan(&item.ID, &item.LoanID, &item.MemberID, &item.BookID, &item.Type, &item.OccurredAt, &fineID, &overdue, &dailyRate, &totalAmount, &currency, &status, &assessed); err != nil {
			return repository.TransactionPage{}, err
		}
		if fineID.Valid {
			item.Fine = &entity.Fine{ID: fineID.String, LoanID: item.LoanID, MemberID: item.MemberID, OverdueDays: int(overdue.Int64), DailyRateMinor: dailyRate.Int64, TotalAmountMinor: totalAmount.Int64, Currency: currency.String, Status: entity.FineStatus(status.String), AssessedAt: assessed.Time}
		}
		items = append(items, &item)
	}
	return repository.TransactionPage{Items: items, TotalItems: total}, rows.Err()
}

func (r *Repository) StockSyncStatus(ctx context.Context, loanID string) (entity.StockSyncStatus, error) {
	var status entity.StockSyncStatus
	err := r.db.GetContext(ctx, &status, `SELECT stock_sync_status FROM loans WHERE id = ?`, loanID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errs.ErrNotFound
	}
	return status, err
}

func (r *Repository) ConfirmStock(ctx context.Context, eventID, causationID, loanID, bookID string, processedAt time.Time) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `INSERT IGNORE INTO inbox_events (event_id, event_type, processed_at) VALUES (?, ?, ?)`, eventID, entity.BookStockUpdatedEventType, processedAt)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return tx.Commit()
	}
	result, err = tx.ExecContext(ctx, `UPDATE loans SET stock_sync_status = 'confirmed', updated_at = GREATEST(updated_at, ?) WHERE id = ? AND book_id = ? AND status = 'returned' AND stock_sync_status = 'pending' AND EXISTS (SELECT 1 FROM outbox_events WHERE id = ? AND event_type = ?)`, processedAt, loanID, bookID, causationID, entity.LoanReturnedEventType)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected != 1 {
		return errs.ErrNotFound
	}
	return tx.Commit()
}

func (r *Repository) PendingOutbox(ctx context.Context, limit int) ([]repository.OutboxMessage, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT id, event_type, routing_key, payload, attempts FROM outbox_events WHERE status = 'pending' AND next_attempt_at <= NOW(6) ORDER BY created_at LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]repository.OutboxMessage, 0)
	for rows.Next() {
		var item repository.OutboxMessage
		if err := rows.Scan(&item.ID, &item.EventType, &item.RoutingKey, &item.Payload, &item.Attempts); err != nil {
			return nil, err
		}
		messages = append(messages, item)
	}
	return messages, rows.Err()
}

func (r *Repository) MarkOutboxPublished(ctx context.Context, id string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE outbox_events SET status = 'published', published_at = ?, last_error = NULL WHERE id = ?`, now, id)
	return err
}

func (r *Repository) MarkOutboxFailed(ctx context.Context, id, message string, now time.Time) error {
	message = strings.TrimSpace(message)
	if len(message) > 1000 {
		message = message[:1000]
	}
	_, err := r.db.ExecContext(ctx, `UPDATE outbox_events SET attempts = attempts + 1, last_error = ?, next_attempt_at = DATE_ADD(?, INTERVAL LEAST(POW(2, attempts + 1), 60) SECOND) WHERE id = ?`, message, now, id)
	return err
}

func getLoan(ctx context.Context, tx *sqlx.Tx, id string, lock bool) (*entity.Loan, error) {
	query := `SELECT id, member_id, book_id, status, stock_sync_status, borrowed_at, due_at, returned_at, created_at, updated_at FROM loans WHERE id = ?`
	if lock {
		query += " FOR UPDATE"
	}
	var loan entity.Loan
	if err := tx.GetContext(ctx, &loan, query, id); errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	if loan.Status == entity.LoanReturned {
		var fine entity.Fine
		err := tx.GetContext(ctx, &fine, `SELECT id AS fine_id, loan_id, member_id AS fine_member_id, overdue_days, daily_rate_minor, total_amount_minor, currency, status AS fine_status, assessed_at FROM fines WHERE loan_id = ?`, id)
		if err == nil {
			loan.Fine = &fine
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	return &loan, nil
}

func calculateOverdueDays(dueAt, returnedAt time.Time) int {
	if !returnedAt.After(dueAt) {
		return 0
	}
	return int(math.Ceil(returnedAt.Sub(dueAt).Hours() / 24))
}

func duplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
