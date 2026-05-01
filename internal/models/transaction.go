package models

func (m *SnippetModel) Transaction() error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	// This ensures there is always a rollback;
	// if the function commits to DB it rollbacks to the final state
	// after the successful transaction
	defer tx.Rollback()

	err = tx.Commit()
	return err
}
