package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"rentEquipement/internal/config"
	"rentEquipement/internal/handlers"
	m "rentEquipement/internal/middleware"
	"rentEquipement/internal/repository/equipment"
	"rentEquipement/internal/repository/maintenance"
	"rentEquipement/internal/repository/review"
	"rentEquipement/internal/repository/user"
	"rentEquipement/internal/session"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := new(config.Config)
	config.ParseConfig(cfg)
	portApp := fmt.Sprintf(":%s", cfg.Port)
	fmt.Println(cfg.DSN)
	db, _ := sql.Open("pgx", cfg.DSN)
	db.SetMaxOpenConns(10)
	defer db.Close()
	err := db.Ping()
	// conn, err := pgx.Connect(context.Background(), cfg.DSN)
	// defer conn.Close(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	templates := template.Must(template.ParseGlob("templates/*"))
	sm := session.NewSessionsManager()
	userRep := user.NewRep(db)
	maintenanceRep := maintenance.NewRep(db)
	equipmentRep := equipment.NewRep(db)
	reviewRep := review.NewRep(db)

	userHandler := handlers.UserHandler{
		Repo:     userRep,
		Tmpl:     templates,
		Sessions: sm,
	}
	maintenanceHandler := handlers.MaintenanceHandler{
		Repo:     maintenanceRep,
		Tmpl:     templates,
		Sessions: sm,
	}
	equipmentHandler := handlers.EquipmentHandler{
		Repo:     equipmentRep,
		Tmpl:     templates,
		Sessions: sm,
	}
	reviewHandler := handlers.ReviewHandler{
		Repo:     reviewRep,
		Tmpl:     templates,
		Sessions: sm,
	}
	auth := m.AuthMiddleware{
		SM: sm,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", equipmentHandler.FrontListEquipments)
	r.Get("/login", userHandler.FormLogin)
	r.Get("/logout", userHandler.Logout)
	r.Get("/register", userHandler.FormRegister)
	r.Get("/equipment/{equipment_ID}", equipmentHandler.FrontEquipmentInfo)

	r.Get("/api/v1/equipments", equipmentHandler.ListEquipments)
	r.Get("/api/v1/equipment/{equipment_ID}", equipmentHandler.EquipmentInfo)
	r.Get("/api/v1/equipment/{equipment_ID}/maintenance", maintenanceHandler.List)
	r.Get("/api/v1/equipment/{equipment_ID}/reviews", reviewHandler.List)
	r.Post("/api/v1/login", userHandler.Login)
	r.Post("/api/v1/register", userHandler.Register)

	r.Group(func(r chi.Router) {
		r.Use(auth.Auth)
		r.Get("/order", equipmentHandler.FormCreateOrder)
		r.Get("/info", userHandler.FrontInfo)
		r.Get("/api/v1/orders", userHandler.ListOrders)
		r.Get("/api/v1/notequipments", equipmentHandler.ListNotEquipments)
		r.Get("/api/v1/customer/{username}", userHandler.CustomerInfo)
		r.Get("/api/v1/admin/equipment/{equipment_ID}/do_available", equipmentHandler.DoAvailable)
		r.Get("/api/v1/admin/getlogs", userHandler.GetLogs)
		r.Post("/api/v1/order", equipmentHandler.CreateOrder)
	})

	log.Println("Server is running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(portApp, r))

}
