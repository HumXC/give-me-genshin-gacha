package database

import (
	"encoding/json"
	"fmt"
	"give-me-genshin-gacha/gacha"
	"io"
	"os"
	"path"
	"testing"
)

var testDir = "../test"
var testDB = path.Join(testDir, "test.db")

func TestMain(m *testing.M) {
	clear()
	m.Run()
}
func clear() {
	os.Remove(testDB)
}
func fakeDB(t *testing.T) GachaDB {
	clear()
	fd, err := fakeData(t)
	if err != nil {
		t.Fatal(err)
	}
	db, err := NewDB(testDB)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Add(fd["200"])
	if err != nil {
		t.Fatal(err)
	}
	err = db.Add(fd["301"])
	if err != nil {
		t.Fatal(err)
	}
	err = db.Add(fd["302"])
	if err != nil {
		t.Fatal(err)
	}
	err = db.Add(fd["100"])
	if err != nil {
		t.Fatal(err)
	}
	return db
}
func fakeData(t *testing.T) (map[string][]gacha.RespDataListItem, error) {
	fakeData := make(map[string][]gacha.RespDataListItem)
	f, err := os.Open(path.Join(testDir, "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, &fakeData)
	if err != nil {
		t.Fatal(err)
	}
	return fakeData, nil
}

func TestAdd(t *testing.T) {
	fd, err := fakeData(t)
	if err != nil {
		t.Fatal(err)
	}
	db, err := NewDB(testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = db.Add(fd["200"])
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUids(t *testing.T) {
	db := fakeDB(t)
	defer db.Close()
	got, err := db.GetUids()
	if err != nil {
		t.Fatal(err)
	}
	want := "[111111111 222222222]"
	if fmt.Sprint(got) != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

func TestGetByRank(t *testing.T) {
	db := fakeDB(t)
	defer db.Close()
	got, err := db.GetGachaRank("111111111", "200", "3")
	if err != nil {
		t.Fatal(err)
	}
	want := "1667390760012370286"
	if got[0].ID != want {
		t.Errorf("got: %v, want: %s", got, want)
	}
}

func TestCountIn(t *testing.T) {
	db := fakeDB(t)
	defer db.Close()
	id1, id2 := "1665741960038012486", "1665741960038216486"
	got, err := db.CountIn("111111111", "200", id1, id2)
	if err != nil {
		t.Fatal(err)
	}
	want := 2
	if got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}

func TestCountGacha(t *testing.T) {
	db := fakeDB(t)
	defer db.Close()
	got, err := db.CountGacha("111111111", "200")
	if err != nil {
		t.Fatal(err)
	}
	want := 4
	if got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
	got, err = db.CountGacha("222222222", "301")
	if err != nil {
		t.Fatal(err)
	}
	want = 1
	if got != want {
		t.Errorf("got: %d, want: %d", got, want)
	}
}
