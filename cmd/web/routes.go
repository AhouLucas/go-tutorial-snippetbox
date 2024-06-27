package main

import (
	"net/http"
    
    "github.com/julienschmidt/httprouter"
    "github.com/justinas/alice" // Package for middleware chaining
)

func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.notFound(w)
    })

    fileServer := http.FileServer(http.Dir("./ui/static"))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

    dynamicMiddleWare := alice.New(app.sessionManager.LoadAndSave)

    // Unprotected routes (doesn't require authentication)
    router.Handler(http.MethodGet, "/", dynamicMiddleWare.ThenFunc(app.home))
    router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMiddleWare.ThenFunc(app.snippetView))
    router.Handler(http.MethodGet, "/user/signup", dynamicMiddleWare.ThenFunc(app.userSignup))
    router.Handler(http.MethodPost, "/user/signup", dynamicMiddleWare.ThenFunc(app.userSignupPost))
    router.Handler(http.MethodGet, "/user/login", dynamicMiddleWare.ThenFunc(app.userLogin))
    router.Handler(http.MethodPost, "/user/login", dynamicMiddleWare.ThenFunc(app.userLoginPost))

    // Protected routes (require authentication)
    protectedMiddleWare := dynamicMiddleWare.Append(app.requireAuthentication)

    router.Handler(http.MethodGet, "/snippet/create", protectedMiddleWare.ThenFunc(app.snippetCreate))
    router.Handler(http.MethodPost, "/snippet/create", protectedMiddleWare.ThenFunc(app.snippetCreatePost))
    router.Handler(http.MethodPost, "/user/logout", protectedMiddleWare.ThenFunc(app.userLogoutPost))

    standardMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

    return standardMiddleWare.Then(router)
}