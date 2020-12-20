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
	if err != nil {
		return err
	}

	// check changes
	if i, _ := result.RowsAffected(); i != 1 {
		return ErrIDNotFound
	}

	return nil
}

// logUsage logs the record with the given id into a usages table.
func (s *service) logUsage(ctx context.Context, id string) error {

	// store
	//TODO insert foreign key
	result, err := s.DB.ExecContext(ctx, `
INSERT INTO
  usages (shortcut_id)
VALUES
  ($1);
	`, id)
	if err != nil {
		return err
	}

	// check changes
	if i, _ := result.RowsAffected(); i != 1 {
		return ErrIDNotFound
	}

	return nil
}
