package datasources

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataSourceFromSQL struct {
	Name        string
	Description string
	Id          uint64
	Value       uint64
	Interval    uint64
}

func GetAllDatasources() []DataSourceFromSQL {

	gad_db_conn_string := os.Getenv("DLC_DB_CONN_STRING")
	gad_db, err := sqlx.Connect("postgres", gad_db_conn_string)
	if err != nil {
		log.Fatalln(err)
	}
	rows, err := gad_db.Queryx("SELECT Id, Name, Description, cast (Value * 1000000000 as bigint), Interval FROM datasources ORDER BY Id ASC")
	gad_results := []DataSourceFromSQL{}
	for rows.Next() {
		var r DataSourceFromSQL
		err = rows.Scan(&r.Id, &r.Name, &r.Description, &r.Value, &r.Interval)
		if err != nil {
			log.Fatalf("Error querying datasources: %v", err)
		}
		gad_results = append(gad_results, r)
	}
	gad_db.Close()

	return gad_results
}

func GetDatasource(id uint64) DataSourceFromSQL {

	gd_conn_string := os.Getenv("DLC_DB_CONN_STRING")
	gd_db, err := sqlx.Connect("postgres", gd_conn_string)
	if err != nil {
		log.Fatalln(err)
	}
	gd_results := DataSourceFromSQL{}
	err = gd_db.Get(&gd_results, "SELECT Id, Name, Description, cast (Value * 1000000000 as bigint), Interval FROM datasources WHERE id=$1", id)
	if err != nil {
		log.Fatalf("Error getting datasource specified by ID: %v", err)
	}
	gd_db.Close()
	return gd_results
}

func HasDatasource(id uint64) bool {

	hd_conn_string := os.Getenv("DLC_DB_CONN_STRING")
	hd_db, err := sqlx.Connect("postgres", hd_conn_string)
	if err != nil {
		log.Fatalln(err)
	}
	hd_results := DataSourceFromSQL{}
	err = hd_db.Get(&hd_results, "SELECT count(id) as id FROM datasources WHERE id=$1", id)
	defer hd_db.Close()
	if hd_results.Id > 0 {
		return true
	} else {
		return false
	}
}
