// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package target

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/boundary/internal/db"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/kms"
	"github.com/hashicorp/boundary/internal/oplog"
)

// CreateAlias inserts a into the repository and returns a new
// Alias containing the alias's PublicId. a is not changed. a must
// contain a valid ScopeId. a must not contain a PublicId. The PublicId is
// generated and assigned by this method. opt is ignored.
//
// Description, DestinationId, and HostId are optional.
//
// Value must be globally unique.
func (r *Repository) CreateAlias(ctx context.Context, a *Alias, opt ...Option) (*Alias, error) {
	const op = "target.(Repository).CreateAlias"
	switch {
	case a == nil:
		return nil, errors.New(ctx, errors.InvalidParameter, op, "nil Alias")
	case a.Alias == nil:
		return nil, errors.New(ctx, errors.InvalidParameter, op, "nil embedded Alias")
	case a.Value == "":
		return nil, errors.New(ctx, errors.InvalidParameter, op, "no value")
	case a.ScopeId == "":
		return nil, errors.New(ctx, errors.InvalidParameter, op, "no scope id")
	case a.PublicId != "":
		return nil, errors.New(ctx, errors.InvalidParameter, op, "public id not empty")
	}
	a = a.clone()

	id, err := newAliasId(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op)
	}
	a.PublicId = id

	oplogWrapper, err := r.kms.GetWrapper(ctx, a.ScopeId, kms.KeyPurposeOplog)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg("unable to get oplog wrapper"))
	}

	metadata := newAliasMetadata(a, oplog.OpType_OP_TYPE_CREATE)

	var newAlias *Alias
	_, err = r.writer.DoTx(
		ctx,
		db.StdRetryCnt,
		db.ExpBackoff{},
		func(_ db.Reader, w db.Writer) error {
			newAlias = a.clone()
			err := w.Create(
				ctx,
				newAlias,
				db.WithOplog(oplogWrapper, metadata),
			)
			if err != nil {
				return errors.Wrap(ctx, err, op)
			}
			return nil
		},
	)

	if err != nil {
		if errors.IsUniqueError(err) {
			switch {
			case strings.Contains(err.Error(), `"alias_value_uq"`):
				return nil, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("alias value %q is already in use", a.Value)))
			case strings.Contains(err.Error(), `"alias_target_scope_id_name_uq"`):
				return nil, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("in scope %q, the name %q is already in use", a.ScopeId, a.Name)))
			}
		}
		if strings.Contains(err.Error(), `violates foreign key constraint "target_fkey"`) {
			return nil, errors.Wrap(ctx, err, op, errors.WithCode(errors.NotFound), errors.WithMsg("target with specified destination id %q was not found", a.GetDestinationId()))
		}
		return nil, errors.Wrap(ctx, err, op)
	}
	return newAlias, nil
}

// UpdateAlias updates the repository entry for a.PublicId with the
// values in a for the fields listed in fieldMask. It returns a new
// Alias containing the updated values and a count of the number of
// records updated. a is not changed.
func (r *Repository) UpdateAlias(ctx context.Context, a *Alias, version uint32, fieldMask []string, opt ...Option) (*Alias, int, error) {
	const op = "target.(Repository).UpdateAlias"
	switch {
	case a == nil:
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "nil Alias")
	case a.Alias == nil:
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "nil embedded Alias")
	case a.PublicId == "":
		return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "no public id")
	case len(fieldMask) == 0:
		return nil, db.NoRowsAffected, errors.New(ctx, errors.EmptyFieldMask, op, "empty field mask")
	}

	var dbMask, nullFields []string
	for _, f := range fieldMask {
		switch {
		case strings.EqualFold("value", f) && a.Value == "":
			return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "value cannot be empty")
		case strings.EqualFold("value", f) && a.Value != "":
			dbMask = append(dbMask, "value")
		case strings.EqualFold("name", f) && a.Name == "":
			nullFields = append(nullFields, "name")
		case strings.EqualFold("name", f) && a.Name != "":
			dbMask = append(dbMask, "name")
		case strings.EqualFold("description", f) && a.Description == "":
			nullFields = append(nullFields, "description")
		case strings.EqualFold("description", f) && a.Description != "":
			dbMask = append(dbMask, "description")
		case strings.EqualFold("DestinationId", f) && a.DestinationId == "":
			nullFields = append(nullFields, "DestinationId")
		case strings.EqualFold("DestinationId", f) && a.DestinationId != "":
			dbMask = append(dbMask, "DestinationId")
		case strings.EqualFold("HostId", f) && a.HostId == "":
			nullFields = append(nullFields, "HostId")
		case strings.EqualFold("HostId", f) && a.HostId != "":
			dbMask = append(dbMask, "HostId")
		default:
			return nil, db.NoRowsAffected, errors.New(ctx, errors.InvalidFieldMask, op, fmt.Sprintf("invalid field mask: %s", f))
		}
	}

	oplogWrapper, err := r.kms.GetWrapper(ctx, a.ScopeId, kms.KeyPurposeOplog)
	if err != nil {
		return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg("unable to get oplog wrapper"))
	}

	a = a.clone()

	metadata := newAliasMetadata(a, oplog.OpType_OP_TYPE_UPDATE)

	var rowsUpdated int
	var returnedAlias *Alias
	_, err = r.writer.DoTx(
		ctx,
		db.StdRetryCnt,
		db.ExpBackoff{},
		func(_ db.Reader, w db.Writer) error {
			returnedAlias = a.clone()
			var err error
			rowsUpdated, err = w.Update(
				ctx,
				returnedAlias,
				dbMask,
				nullFields,
				db.WithOplog(oplogWrapper, metadata),
				db.WithVersion(&version),
			)
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
		if errors.IsUniqueError(err) {
			switch {
			case strings.Contains(err.Error(), `"alias_value_uq"`):
				return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("for alias %s: alias value %q is already in use", a.PublicId, a.Value)))
			case strings.Contains(err.Error(), `"alias_target_scope_id_name_uq"`):
				return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("in scope %s, the name %q is already in use", a.ScopeId, a.Name)))
			}
		}
		if strings.Contains(err.Error(), `violates foreign key constraint "target_fkey"`) {
			return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithCode(errors.NotFound), errors.WithMsg("target with specified destination id %q was not found", a.GetDestinationId()))
		}
		return nil, db.NoRowsAffected, errors.Wrap(ctx, err, op)
	}

	return returnedAlias, rowsUpdated, nil
}

// LookupAlias returns the Alias for id. Returns nil, nil if no
// Alias is found for id.
func (r *Repository) LookupAlias(ctx context.Context, id string, opt ...Option) (*Alias, error) {
	const op = "target.(Repository).LookupAlias"
	if id == "" {
		return nil, errors.New(ctx, errors.InvalidParameter, op, "no public id")
	}
	a := allocAlias()
	a.PublicId = id
	if err := r.reader.LookupByPublicId(ctx, a); err != nil {
		if errors.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("failed for: %s", id)))
	}
	return a, nil
}

// DeleteAlias deletes id from the repository returning a count of the
// number of records deleted.
func (r *Repository) DeleteAlias(ctx context.Context, id string, opt ...Option) (int, error) {
	const op = "target.(Repository).DeleteAlias"
	if id == "" {
		return db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "no public id")
	}

	a := allocAlias()
	a.PublicId = id
	if err := r.reader.LookupByPublicId(ctx, a); err != nil {
		if errors.IsNotFoundError(err) {
			return db.NoRowsAffected, nil
		}
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("failed for %s", id)))
	}
	if a.ScopeId == "" {
		return db.NoRowsAffected, errors.New(ctx, errors.InvalidParameter, op, "no project id")
	}
	oplogWrapper, err := r.kms.GetWrapper(ctx, a.ScopeId, kms.KeyPurposeOplog)
	if err != nil {
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg("unable to get oplog wrapper"))
	}

	metadata := newAliasMetadata(a, oplog.OpType_OP_TYPE_DELETE)

	var rowsDeleted int
	var deleteAlias *Alias
	_, err = r.writer.DoTx(
		ctx,
		db.StdRetryCnt,
		db.ExpBackoff{},
		func(_ db.Reader, w db.Writer) error {
			deleteAlias = a.clone()
			var err error
			rowsDeleted, err = w.Delete(
				ctx,
				deleteAlias,
				db.WithOplog(oplogWrapper, metadata),
			)
			if err != nil {
				return errors.Wrap(ctx, err, op)
			}
			if rowsDeleted > 1 {
				return errors.New(ctx, errors.MultipleRecords, op, "more than 1 resource would have been deleted")
			}
			return nil
		},
	)

	if err != nil {
		return db.NoRowsAffected, errors.Wrap(ctx, err, op, errors.WithMsg(fmt.Sprintf("delete failed for %s", a.PublicId)))
	}

	return rowsDeleted, nil
}
