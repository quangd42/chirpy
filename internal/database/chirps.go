package database

import "errors"

type ChirpTable struct {
	Chirps    map[int]Chirp `json:"chirps"`
	NextIndex int           `json:"max_index"`
}

type Chirp struct {
	Body     string `json:"body"`
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	// create Chirp, give it ID
	chirp := Chirp{
		ID:       dbs.ChirpTable.NextIndex,
		Body:     body,
		AuthorID: authorID,
	}
	dbs.ChirpTable.Chirps[dbs.ChirpTable.NextIndex] = chirp
	dbs.ChirpTable.NextIndex++
	// writeDB
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}
	for _, chirp := range dbs.ChirpTable.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirpsByAuthor(authorID int) ([]Chirp, error) {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}
	for _, chirp := range dbs.ChirpTable.Chirps {
		if authorID == -1 {
			chirps = append(chirps, chirp)
		} else if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbs.ChirpTable.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp doesn't exist")
	}
	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	// loadDB
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbs.ChirpTable.Chirps, id)

	err = db.writeDB(dbs)
	if err != nil {
		return err
	}
	return nil
}
