package data

import "context"

// incrementUsage increments usage column at record's with the given id.
func (s *service) incrementUsage(ctx context.Context, id string) error {

	// increment
	result, err := s.DB.ExecContext(ctx, `
UPDATE
  shortcuts
SET
  usage = usage + 1
WHERE
  shortcut_id = $1
  AND deleted_at IS NULL;
	`, id)

	// check changes
	if i, _ := result.RowsAffected(); i != 1 {
		return ErrIDNotFound
	}

	return err
}
