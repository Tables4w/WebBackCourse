package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type form struct {
	fio      string
	tel      string
	email    string
	date     string
	gender   string
	favlangs []int
	bio      string
}

func process(w http.ResponseWriter, r *http.Request) {
	var formerrors []int
	r.ParseForm()
	var f form
	err := Validate(&f, r.Form, formerrors)
	if err != nil {
		fmt.Print(err)
	} else {
		err := WriteForm(&f)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func Validate(f *form, form url.Values, formerrors []int) (err error) {
	var check bool = false
	for key, value := range form {
		if key == "Fio" {
			var v string = value[0]
			r, err := regexp.Compile(`^[A-Za-zА-Яа-яЁё\s]{1,150}$`)
			if err != nil {
				fmt.Print(err)
			}
			if !r.MatchString(v) {
				formerrors = append(formerrors, 1)
			} else {
				f.fio = v
			}
		}

		if key == "Tel" {
			var v string = value[0]
			r, err := regexp.Compile(`^\+[0-9]{1,29}$`)
			if err != nil {
				fmt.Print(err)
			}
			if !r.MatchString(v) {
				formerrors = append(formerrors, 2)
			} else {
				f.tel = v
			}
		}

		if key == "Email" {
			var v string = value[0]
			r, err := regexp.Compile(`^[A-Za-z0-9._%+-]{1,30}@[A-Za-z0-9.-]{1,20}\.[A-Za-z]{1,10}$`)
			if err != nil {
				fmt.Print(err)
			}
			if !r.MatchString(v) {
				formerrors = append(formerrors, 3)
			} else {
				f.email = v
			}
		}

		if key == "Birth_date" {
			var v string = value[0]
			r, err := regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
			if err != nil {
				fmt.Print(err)
			}
			if !r.MatchString(v) {
				formerrors = append(formerrors, 4)
			} else {
				f.date = v
			}
		}

		if key == "Gender" {
			var v string = value[0]
			if v != "Male" && v != "Female" {
				formerrors = append(formerrors, 5)
			} else {
				f.gender = v
			}
		}

		if key == "Bio" {
			var v string = value[0]
			f.bio = v
		}

		if key == "Familiar" {
			var v string = value[0]
			if v == "on" {
				check = true
			}
		}

		if key == "Favlangs" {
			for _, p := range value {
				np, err := strconv.Atoi(p)
				if err != nil {
					fmt.Print(err)
					formerrors = append(formerrors, 6)
					break
				} else {
					if np < 1 || np > 11 {
						formerrors = append(formerrors, 6)
						break
					} else {
						f.favlangs = append(f.favlangs, np)
					}
				}
			}
		}
	}
	if !check {
		formerrors = append(formerrors, 8)
	}
	if len(formerrors) == 0 {
		return nil
	}
	return errors.New("validation failed")
}

func WriteForm(f *form) (err error) {
	/*
		postgresHost := os.Getenv("POSTGRES_HOST")
		postgresUser := os.Getenv("POSTGRES_USER")
		postgresPassword := os.Getenv("POSTGRES_PASSWORD")
		postgresDB := os.Getenv("POSTGRES_DB")
		connectStr := "host=" + postgresHost + " user=" + postgresUser +
		" password=" + postgresPassword +
		" dbname=" + postgresDB + " sslmode=disable"
	*/
	postgresUser := "postgres"
	postgresPassword := "123"
	postgresDB := "back3"
	connectStr := "user=" + postgresUser +
		" password=" + postgresPassword +
		" dbname=" + postgresDB + " sslmode=disable"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return err
	}
	defer db.Close()
	var insertsql = []string{
		"INSERT INTO forms",
		"(fio, tel, email, birth_date, gender, bio)",
		"VALUES ($1, $2, $3, $4, $5, $6) returning form_id",
	}
	var form_id int
	err = db.QueryRow(strings.Join(insertsql, ""), f.fio, f.tel,
		f.email, f.date, f.gender, f.bio).Scan(&form_id)
	if err != nil {
		fmt.Print("SEEEEEEEEEEEX")
		return err
	}

	for _, v := range f.favlangs {
		_, err = db.Exec("INSERT INTO favlangs VALUES ($1, $2)", form_id, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/process", process)
	server.ListenAndServe()
}
