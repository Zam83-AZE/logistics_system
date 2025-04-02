package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zam83-AZE/logistics_system/internal/domain/auth"
	"github.com/Zam83-AZE/logistics_system/internal/domain/dashboard"
	"github.com/Zam83-AZE/logistics_system/internal/middleware"
	"github.com/Zam83-AZE/logistics_system/pkg/db"
	"github.com/Zam83-AZE/logistics_system/pkg/logger"
	"github.com/Zam83-AZE/logistics_system/pkg/session"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func main() {
	// Loq inisializasiyası
	log := logger.NewLogger()
	log.Info("Logistics System işə salınır")

	// Verilənlər bazasına qoşulma
	database, err := db.Connect()
	if err != nil {
		log.WithError(err).Fatal("Verilənlər bazasına qoşulma xətası")
	}
	defer database.Close()

	// Router inisializasiyası
	router := mux.NewRouter()

	// Sessiya mağazası yaratma
	store := sessions.NewCookieStore([]byte("logistics-system-secret-key"))
	sessionManager := session.NewManager(store)

	// // Mütləq yol istifadə edin
	// wd, err := os.Getwd()
	// fmt.Println("Cari işçi qovluq:", wd)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// staticDir := filepath.Join(wd, "web", "static")
	// fmt.Println("staticDir qovluq:", staticDir)

	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Statik faylların təqdim edilməsi
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// tmplFiles, _ := filepath.Glob("web/templates/**/*.html")
	// if len(tmplFiles) == 0 {
	// 	log.Fatal("Heç bir şablon faylı tapılmadı!")
	// }
	// fmt.Println("Tapılan şablonlar:", tmplFiles)

	// Şablonların emalı
	tmpl, err := template.ParseGlob("web/templates/**/*.html")
	//tmpl := template.New("templates")
	// Hər bir şablonu açıq şəkildə yükləyin
	// _, err = tmpl.ParseFiles(
	// 	filepath.Join(wd, "web", "templates", "layout.html"),
	// 	filepath.Join(wd, "web", "templates", "auth", "login.html"),
	// 	filepath.Join(wd, "web", "templates", "dashboard", "index.html"),
	// )
	if err != nil {
		log.WithError(err).Fatal("Şablonların emalı zamanı xəta")
	}

	// Middleware tətbiqi
	router.Use(middleware.Logging(log))
	router.Use(sessionManager.Middleware)

	// Marşrutların qeydiyyatı
	auth.RegisterRoutes(router, database, tmpl, sessionManager)

	// Autentifikasiya tələb edən marşrutlar üçün alt-router
	secureRouter := router.PathPrefix("/").Subrouter()
	secureRouter.Use(middleware.RequireAuth(sessionManager))

	// Dashboard marşrutlarının qeydiyyatı
	dashboard.RegisterRoutes(secureRouter, database, tmpl)

	// Server tərifləri
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Serverin başladılması
	go func() {
		log.Info("Server :8080 portunda başladılır")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server başlatma xətası")
		}
	}()

	// Təmiz bağlanma
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Info("Server bağlanır")
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server məcburi bağlandı")
	}

	log.Info("Server bağlandı")
}
