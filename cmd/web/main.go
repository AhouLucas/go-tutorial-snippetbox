package main

import (
    "database/sql"
    "flag"
    "html/template"
    "log"
    "net/http"
    "os"
    "time"
    "snippetbox.gotutorial/internal/models"

    "github.com/alexedwards/scs/sqlite3store"
    "github.com/alexedwards/scs/v2"
    "github.com/go-playground/form/v4"
    _ "github.com/mattn/go-sqlite3"
)



type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
    formDecoder *form.Decoder
    sessionManager *scs.SessionManager
}


func main() {
    // CLI arguments
    addr := flag.String("addr", ":4000", "HTTP network address")

    // Database source name (DSN)
    dsn := flag.String("dsn", "./db/snippetbox.db", "SQLite data source name")

    flag.Parse()

    // Custom loggers
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    // Create DB connection
    db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }

    // When main finishes, close the db connection
    defer db.Close()

    templateCache, err := newTemplateCache()
    if err != nil {
        errorLog.Fatal(err)
    }

    // Initialize a new decoder instance for form parsing
    formDecoder := form.NewDecoder()

    // Initialize a new session manager and configure it to use our sqlite database
    sessionManager := scs.New()
    sessionManager.Store = sqlite3store.New(db)
    sessionManager.Lifetime = 12 * time.Hour
    sessionManager.Cookie.Secure = true // Send cookies over https instead of http

    // Create application struct that contains dependencies for other .go files
    app := &application{
        errorLog: errorLog,
        infoLog: infoLog,
        snippets: &models.SnippetModel{DB: db},
        templateCache: templateCache,
        formDecoder: formDecoder,
        sessionManager: sessionManager,
    }

    // Create a custom http server to use our error logger instead of the standard logger
    // And that also contains our servermux
    srv := &http.Server{
        Addr: *addr,
        ErrorLog: errorLog,
        Handler: app.routes(),
    }

    infoLog.Printf("Starting server on %s", *addr)
    err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
    errorLog.Fatal(err)
}


func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}