// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/boundary/internal/db"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/kms"
	"github.com/hashicorp/boundary/internal/oplog"
	"github.com/hashicorp/go-dbw"
)

// CreateAccount inserts an Account, a, into the repository and returns a
// new Account containing its PublicId. a is not changed. a must contain a
// valid LdapMethodId and ScopeId. a must not contain a PublicId. The PublicId
// is generated and assigned by this method. a must contain a valid LoginName.
// a.LoginName must be unique for an a.AuthMethod.
//
// Both a.Name and a.Description are optional. If a.Name is set, it must be
// unique within a.AuthMethodId.
func (r *Repository) CreateAccount(ctx context.Context, a *Account, _ ...Option) (*Account, error) {
	const op = "ldap.(Repository).CreateAccount"
	switch {
	case a == nil:
		return nil, errors.New(ctx, errors.InvalidParameter, op, "missing account")
	case a.Account == nil:
		return nil, errors.New(ctx, errors.InvalidParameter, op, "missing embedded account")
	case a.PublicId != "":
		return nil, errors.New(ctx, errors.InvalidParameter, op, "public id must be empty")
	default:
		if err := a.validate(ctx, op); err != nil {
			return nil, err // err already wrapped
		}
	}
	id, err := newAccountId(ctx, a.AuthMethodId, a.LoginName)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op)
	}
	a.PublicId = id

	oplogWrapper, err := r.kms.GetWrapper(ctx, a.ScopeId, kms.KeyPurposeOplog)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg("unable to get oplog wrapper"), errors.WithCode(errors.Encrypt))
	}

	md, err := a.oplog(ctx, oplog.OpType_OP_TYPE_CREATE)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg("unable to generate account oplog metadata"))
	}
	var newAccount *Account
	_, err = r.writer.DoTx(ctx, db.StdRetryCnt, db.ExpBackoff{},
		func(_ db.Reader, w db.Writer) error {
			newAccount = a.clone()
			if err := w.Create(ctx, newAccount, db.WithOplog(oplogWrapper, md)); err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		switch {
		case errors.IsUniqueError(err):
			return nil, errors.New(ctx, errors.NotUnique, op, fmt.Sprintf(
				"in auth method %s: name %q already exists or login name %q already exists for auth method %q in scope %s",
				a.AuthMethodId, a.Name, a.LoginName, a.AuthMethodId, a.ScopeId))
		default:
			return nil, errors.Wrap(ctx, err, op, errors.WithMsg(a.AuthMethodId))
		}
	}
	return newAccount, nil
}

// LookupAccount will look up an account in the repository.  If the account is not
// found, it will return nil, nil.  All options are ignored.
func (r *Repository) LookupAccount(ctx context.Context, withPublicId string, _ ...Option) (*Account, error) {
	const op = "ldap.(Repository).LookupAccount"
	if withPublicId == "" {
		return nil, errors.New(ctx, errors.InvalidPublicId, op, "missing public id")
	}
	a := AllocAccount()
	a.PublicId = withPublicId
	if err := r.reader.LookupByPublicId(ctx, a); err != nil {
		switch {
		case errors.IsNotFoundError(err):
			return nil, nil
		default:
			return nil, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("failed for %s", withPublicId)))
		}
	}
	return a, nil
}

// ListAccounts in an auth method and supports the following options:
//   - WithLimit
//   - WithStartPageAfterItem
func (r *Repository) ListAccounts(ctx context.Context, withAuthMethodId string, opt ...Option) ([]*Account, error) {
	const op = "ldap.(Repository).ListAccounts"
	if withAuthMethodId == "" {
		return nil, errors.New(ctx, errors.InvalidParameter, op, "missing auth method id")
	}
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op)
	}

	limit := r.defaultLimit
	if opts.withLimit != 0 {
		// non-zero signals an override of the default limit for the repo.
		limit = opts.withLimit
	}

	var args []any
	args = append(args, sql.Named("with_auth_method_id", withAuthMethodId))
	whereClause := "auth_method_id = @with_auth_method_id"

	// Ordering and pagination are tightly coupled.
	// We order by update_time ascending so that new
	// and updated items appear at the end of the pagination.
	// We need to further order by public_id to distinguish items
	// with identical update times.
	withOrder := "update_time asc, public_id asc"
	if opts.withStartPageAfterItem != nil {
		// Now that the order is defined, we can use a simple where
		// clause to only include items updated since the specified
		// start of the page. We use greater than or equal for the update
		// time as there may be items with identical update_times. We
		// then use PublicId as a tiebreaker.
		args = append(args,
			sql.Named("after_item_update_time", opts.withStartPageAfterItem.GetUpdateTime()),
			sql.Named("after_item_id", opts.withStartPageAfterItem.GetPublicId()),
		)
		whereClause = "(" + whereClause + ") and (update_time > @after_item_update_time or (update_time = @after_item_update_time and public_id > @after_item_id))"
	}

	var accts []*Account
	err = r.reader.SearchWhere(ctx, &accts, whereClause, args, db.WithLimit(limit), db.WithOrder(withOrder))
	if err != nil {
		return nil, errors.Wrap(ctx, err, op)
	}
	return accts, nil
}

// listDeletedIds lists the public IDs of any accounts deleted since the timestamp provided.
func (r *Repository) listDeletedIds(ctx context.Context, since time.Time) ([]string, time.Time, error) {
	const op = "account.(Repository).listDeletedIds"
	var deletedAccounts []*deletedAccount
	var transactionTimestamp time.Time
	if _, err := r.writer.DoTx(ctx, db.StdRetryCnt, db.ExpBackoff{}, func(r db.Reader, w db.Writer) error {
		if err := r.SearchWhere(ctx, &deletedAccounts, "delete_time >= ?", []any{since}); err != nil {
			return errors.Wrap(ctx, err, op, errors.WithMsg("failed to query deleted accounts"))
		}
		var err error
		transactionTimestamp, err = r.Now(ctx)
		if err != nil {
			return errors.Wrap(ctx, err, op, errors.WithMsg("failed to query transaction timestamp"))
		}
		return nil
	}); err != nil {
		return nil, time.Time{}, err
	}
	var accountIds []string
	for _, acct := range deletedAccounts {
		accountIds = append(accountIds, acct.PublicId)
	}
	return accountIds, transactionTimestamp, nil
}

// estimatedCount returns and estimate of the total number of items in the accounts table.
func (r *Repository) estimatedCount(ctx context.Context) (int, error) {
	const op = "account.(Repository).estimatedCount"
	rows, err := r.reader.Query(ctx, estimateCountAccounts, nil)
	if err != nil {
		return 0, errors.Wrap(ctx, err, op, errors.WithMsg("failed to query total accounts"))
	}
	var count int
	for rows.Next() {
		if err := r.reader.ScanRows(ctx, rows, &count); err != nil {
			return 0, errors.Wrap(ctx, err, op, errors.WithMsg("failed to query total accounts"))
		}
	}
	return count, nil
}

// DeleteAccount deletes the account for the provided id from the repository returning a count of the
// number of records deleted.  All options are ignored.
func (r *Repository) DeleteAccount(ctx context.Context, withPublicId string, _ ...Option) (int, error) {
	const op = "ldap.(Repository).DeleteAccount"
	switch {
	case withPublicId == "":
		return db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "missing public id")
	}
	ac := AllocAccount()
	ac.PublicId = withPublicId

	if err := r.reader.LookupById(ctx, ac); err != nil {
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg("account not found"))
	}

	oplogWrapper, err := r.kms.GetWrapper(ctx, ac.ScopeId, kms.KeyPurposeOplog)
	if err != nil {
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithCode(errors.Encrypt), errors.WithMsg("unable to get oplog wrapper"))
	}
	metadata, err := ac.oplog(ctx, oplog.OpType_OP_TYPE_DELETE)
	if err != nil {
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg("unable to generate oplog metadata"))
	}

	var rowsDeleted int
	_, err = r.writer.DoTx(
		ctx,
		db.StdRetryCnt,
		db.ExpBackoff{},
		func(_ db.Reader, w db.Writer) (err error) {
			dAc := ac.clone()
			rowsDeleted, err = w.Delete(ctx, dAc, db.WithOplog(oplogWrapper, metadata))
			switch {
			case err != nil:
				return errors.Wrap(ctx, err, op, errors.WithMsg("unable to delete ldap account"))
			case rowsDeleted > 1:
				return errors.New(ctx, errors.MultipleRecords, op, "more than 1 resource would have been deleted")
			}
			return nil
		},
	)
	if err != nil {
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg(withPublicId))
	}

	return rowsDeleted, nil
}

// UpdateAccount updates the repository entry for a.PublicId with the
// values in a for the fields listed in fieldMaskPaths. It returns a new
// Account containing the updated values and a count of the number of
// records updated. a is not changed.
//
// a must contain a valid PublicId. Only a.Name and a.Description can be
// updated. If a.Name is set to a non-empty string, it must be unique within
// a.AuthMethodId.
//
// An attribute of a will be set to NULL in the database if the attribute
// in a is the zero value and it is included in fieldMaskPaths.
func (r *Repository) UpdateAccount(ctx context.Context, scopeId string, a *Account, version uint32, fieldMaskPaths []string, opt ...Option) (*Account, int, error) {
	const op = "ldap.(Repository).UpdateAccount"
	switch {
	case a == nil:
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "missing Account")
	case scopeId == "":
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "missing scope id")
	case version == 0:
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "missing version")
	case a.Account == nil:
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "missing embedded Account")
	case a.PublicId == "":
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidPublicId, op, "missing public id")
	}
	for _, f := range fieldMaskPaths {
		switch {
		case strings.EqualFold(NameField, f):
		case strings.EqualFold(DescriptionField, f):
		default:
			return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidFieldMask, op, f)
		}
	}
	var dbMask, nullFields []string
	dbMask, nullFields = dbw.BuildUpdatePaths(
		map[string]any{
			NameField:        a.Name,
			DescriptionField: a.Description,
		},
		fieldMaskPaths,
		nil,
	)
	if len(dbMask) == 0 && len(nullFields) == 0 {
		return nil, db.NoRowsAffected, errors.New(ctx, errors.EmptyFieldMask, op, "missing field mask")
	}

	oplogWrapper, err := r.kms.GetWrapper(ctx, scopeId, kms.KeyPurposeOplog)
	if err != nil {
		return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithCode(errors.Encrypt),
			errors.WithMsg(("unable to get oplog wrapper")))
	}

	metadata, err := a.oplog(ctx, oplog.OpType_OP_TYPE_UPDATE)
	if err != nil {
		return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg("unable to generate oplog metadata"))
	}

	var rowsUpdated int
	var returnedAccount *Account
	_, err = r.writer.DoTx(ctx, db.StdRetryCnt, db.ExpBackoff{},
		func(_ db.Reader, w db.Writer) error {
			returnedAccount = a.clone()
			var err error
			rowsUpdated, err = w.Update(ctx, returnedAccount, dbMask, nullFields, db.WithOplog(oplogWrapper, metadata), db.WithVersion(&version))
			if err != nil {
				return errors.Wrap(ctx, err, op)
			}
			if rowsUpdated > 1 {
				return errors.New(ctx, errors.MultipleRecords, op, "more than 1 resource would have been updated")
			}
			return nil
		},
	)
	if err != nil {
		switch {
		case errors.IsUniqueError(err):
			return nil, db.NoRowsAffected, errors.New(ctx, errors.NotUnique, op,
				fmt.Sprintf("name %s already exists: %s", a.Name, a.PublicId))
		default:
			return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg(a.PublicId))
		}
	}

	return returnedAccount, rowsUpdated, nil
}
