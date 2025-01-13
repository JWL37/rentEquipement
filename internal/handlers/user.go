package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"rentEquipement/internal/repository/user"
	"rentEquipement/internal/session"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo     user.UserRepo
	Tmpl     *template.Template
	Sessions *session.SessionsManager
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.Repo.Authorize(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Error fetching user:", err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Password mismatch:", err)
		return
	}

	_, err = h.Sessions.Create(w, user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error creating session:", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *UserHandler) FormLogin(w http.ResponseWriter, r *http.Request) {

	err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	email := r.FormValue("email")
	phoneNumber := r.FormValue("phoneNumber")
	address := r.FormValue("address")
	log.Println("Registering user:", username, email, phoneNumber, address)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		log.Println("Password hash error:", err)
		return
	}
	user, err := h.Repo.CreateUser(username, hashedPassword, email, phoneNumber, address)

	log.Println(user.Role)

	h.Sessions.Create(w, user.ID, user.Username, user.Role)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *UserHandler) FormRegister(w http.ResponseWriter, r *http.Request) {

	err := h.Tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) FrontInfo(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if sess.Role == "customer" {
		err = h.Tmpl.ExecuteTemplate(w, "customer.html", nil)
		if err != nil {
			http.Error(w, `Template errror`, http.StatusInternalServerError)
			return
		}
		return
	}
	if sess.Role == "admin" {
		err = h.Tmpl.ExecuteTemplate(w, "admin.html", nil)
		if err != nil {
			http.Error(w, `Template errror`, http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "not authorize", http.StatusForbidden)

}

func (h *UserHandler) CustomerInfo(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	log.Println("CustomerInfo: ", username)

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if sess.Role != "customer" {
		http.Error(w, "not customer", http.StatusForbidden)
		return
	}

	usernameFront, err := GetCookieValue(r, "username")
	if usernameFront != username {
		http.Error(w, "no your username", http.StatusForbidden)
		return
	}
	info, err := h.Repo.GetCustomerByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(*info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCookieValue извлекает значение куки по её имени
func GetCookieValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("cookie %s not found", cookieName)
		}
		return "", fmt.Errorf("error retrieving cookie %s: %w", cookieName, err)
	}
	return cookie.Value, nil
}

func (h *UserHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if sess.Role != "admin" {
		http.Error(w, "not admin", http.StatusForbidden)
		return
	}

	logs, err := h.Repo.GetLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if sess.Role != "customer" {
		http.Error(w, "not customer", http.StatusForbidden)
		return
	}

	ordersForCustomer, err := h.Repo.GetListOrdersForCustomer(sess.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ordersForCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	expiredCookie := &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, expiredCookie)
	expiredUsernameCookie := &http.Cookie{
		Name:    "username",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, expiredUsernameCookie)
	// err := h.Sessions.DestroyCurrent(w, r)
	// if err != nil {
	// 	log.Printf("Logout. Error: %s. \n", err)
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	http.Redirect(w, r, "/", http.StatusFound)
}
