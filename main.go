package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"fmt"
	"time"
	"os"
	"database/sql"
	"regexp"
	"strconv"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"github.com/go-sql-driver/mysql"

)

var db *sql.DB

type Photograph struct {
	ID int
	Name string
	Title sql.NullString
	Description sql.NullString
}

type Collection struct {
	ID int
	Title string
	Description string
	Photographs []Photograph
}

type GalleryPipeline struct {
	Collections []Collection
}

type CollectionPipeline struct {
	Collections []Collection
	Viewing Collection
}

var projectID = "exalted-kayak-372004"

var templates = template.Must(template.ParseFiles("gallery.gohtml", "collection.gohtml"))

// listFiles lists objects within specified bucket.
func listFiles(bucket string) ([]string, error) {
	var output []string
	// bucket := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return output, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return output, fmt.Errorf("Bucket(%q).Objects: %v", bucket, err)
		}
		output = append(output, attrs.Name)
	}
	return output, nil
}

func getCollections() ([]Collection, error) {
	var collections []Collection

	rows, err := db.Query("SELECT * FROM collections")
	if err != nil {
		return nil, fmt.Errorf("collections: %v", err)
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Title, &c.Description); err != nil {
			return nil, fmt.Errorf("collections: %v", err)
		}
		collections = append(collections, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("articles: %v", err)
	}
	return collections, nil
}

func getCollection(collectionID int) (Collection, error) {
	var c Collection
	row := db.QueryRow("SELECT title, description FROM collections WHERE id = ?", collectionID)
	err := row.Scan(&c.Title, &c.Description);
	if err != nil {
		return Collection{}, fmt.Errorf("collection %d: %v", collectionID, err)
	}

	rows, err := db.Query("SELECT id, name, title, description FROM photographs INNER JOIN entries ON photographs.id = entries.photoID WHERE collectionID = ?", collectionID)
	if err != nil {
		return Collection{}, fmt.Errorf("collection %d: %v", collectionID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Photograph
		if err := rows.Scan(&p.ID, &p.Name, &p.Title, &p.Description); err != nil {
			return Collection{}, fmt.Errorf("collection: %v", err)
		}
		c.Photographs = append(c.Photographs, p)
	}
	return c, nil
}

func getPhotographs() ([]Photograph, error) {
	// An articles slice to hold data from returned rows.
	var photographs []Photograph

	rows, err := db.Query("SELECT * FROM photographs")
	if err != nil {
		return nil, fmt.Errorf("articles: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var p Photograph
		if err := rows.Scan(&p.ID, &p.Name, &p.Title, &p.Description); err != nil {
			return nil, fmt.Errorf("articles: %v", err)
		}
		photographs = append(photographs, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("articles: %v", err)
	}
	return photographs, nil
}

func dbSetup() {

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Args[1],
		Passwd: os.Args[2],
		Net:    "tcp",
		Addr:   "database.diegogutierrez.org",
		DBName: "photography",
		ParseTime: true,
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Print("Connected to database!\n")
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	cs, _ := getCollections()
	gp := GalleryPipeline{Collections: cs}
	err := templates.ExecuteTemplate(w, "gallery.gohtml", gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/collection/([0-9]+)$")

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
	}
	id, _ := strconv.Atoi(m[1])
	c, _ := getCollection(id)
	cs, _ := getCollections()
	gp := CollectionPipeline{Collections: cs, Viewing: c}
	err := templates.ExecuteTemplate(w, "collection.gohtml", gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	dbSetup()

	http.HandleFunc("/", galleryHandler)
	http.HandleFunc("/collection/", collectionHandler)

	// CSS files don't require special serving
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets"))))

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}