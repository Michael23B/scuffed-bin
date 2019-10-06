package cmd

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "github.com/davecgh/go-spew/spew"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "os"

  _ "github.com/lib/pq"
  "github.com/spf13/cobra"
)

type Post struct {
  ID   string `json:"post-id"`
  Body string `json:"post-body"`
  Ts   string `json:"time-stamp"`
}

var (
  host     = os.Getenv("DBHOST")
  port     = os.Getenv("DBPORT")
  user     = os.Getenv("DBUSER")
  password = os.Getenv("DBPASSWORD")
  dbname   = os.Getenv("DBNAME")
)

var globalDB *sql.DB

func getPost(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "*")
  w.Header().Set("Content-Type", "application/json")
  var post Post
  postID := mux.Vars(r)["id"]
  log.Println("[getPost] - Fetching ID: " + postID)

  sqlStatement := fmt.Sprintf(`
	select uri, posts, date 
	from posts 
	where uri = '%s'`, postID)

  err := globalDB.QueryRow(sqlStatement).Scan(&post.ID, &post.Body, &post.Ts)
  if err != nil {
    log.Printf("[getPost] - err: %v\n", err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  fmt.Printf("[getPost] - Found record {ID: %s, TimeStamep: %s }\n", post.ID, post.Ts)
  if err = json.NewEncoder(w).Encode(post); err != nil {
    log.Printf("[getPost] - err: %v\n", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
}

func uploadPost(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "*")
  w.Header().Set("Content-Type", "application/json")

  var post Post
  _ = json.NewDecoder(r.Body).Decode(&post)
  spew.Dump(post)

  sqlStatement := fmt.Sprintf(`
	INSERT INTO posts (posts)
	VALUES ('%s')
	RETURNING uri, date`, post.Body)

  err := globalDB.QueryRow(sqlStatement).Scan(&post.ID, &post.Ts)
  if err != nil {
    log.Printf("[uploadPost] - err: %v\n", err)
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  fmt.Printf("[uploadPost] - New record {ID: %s, TimeStamep: %s }\n", post.ID, post.Ts)
  if err = json.NewEncoder(w).Encode(post); err != nil {
    log.Printf("[uploadPost] - err: %v\n", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
}

var rootCmd = &cobra.Command{
  Use:   "server",
  Short: "Webserver for scuffed-bin project",
  Long:  `Webserver for scuffed-bin project`,
  Run: func(cmd *cobra.Command, args []string) {
    var err error
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
      "password=%s dbname=%s sslmode=disable",
      host, port, user, password, dbname)

    globalDB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
      panic(err)
    }
    defer globalDB.Close()

    err = globalDB.Ping()
    if err != nil {
      log.Printf("Error on pinging database: %v\n", err)
      os.Exit(1)
    }

    fmt.Println("Successfully connected!")

    router := mux.NewRouter()
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, "../src/index.html")
    })

    router.HandleFunc("/post/{id}", getPost).Methods("GET")
    router.HandleFunc("/post", uploadPost).Methods("POST", "OPTIONS")
    log.Fatal(http.ListenAndServe(":8070", router))
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
}
