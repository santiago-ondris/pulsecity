package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/pulsecity/services/gateway/internal/domain"
)

type actor struct {
	kind         string
	guestToken   string
	sessionToken string
	user         domain.User
}

func (d Dependencies) register(w http.ResponseWriter, r *http.Request) {
	var request domain.RegisterRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	displayName := strings.TrimSpace(request.DisplayName)
	password := request.Password

	if email == "" || displayName == "" || len(password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "email, display_name y password de al menos 8 caracteres son obligatorios",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if _, _, found, err := d.Store.GetUserCredentialsByEmail(ctx, email); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to validate email",
		})
		return
	} else if found {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": "email already registered",
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to hash password",
		})
		return
	}

	user, err := d.Store.CreateUser(ctx, email, displayName, string(passwordHash))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
		return
	}

	session, err := d.Store.CreateUserSession(ctx, user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create user session",
		})
		return
	}

	writeJSON(w, http.StatusCreated, session)
}

func (d Dependencies) login(w http.ResponseWriter, r *http.Request) {
	var request domain.LoginRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	user, passwordHash, found, err := d.Store.GetUserCredentialsByEmail(ctx, request.Email)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load user",
		})
		return
	}
	if !found || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.Password)) != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid credentials",
		})
		return
	}

	session, err := d.Store.CreateUserSession(ctx, user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create user session",
		})
		return
	}

	writeJSON(w, http.StatusOK, session)
}

func (d Dependencies) getCurrentSession(w http.ResponseWriter, r *http.Request) {
	sessionToken := sessionTokenFromRequest(r)
	if sessionToken == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing session token",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ok, err := d.Store.TouchUserSession(ctx, sessionToken)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to validate session",
		})
		return
	}
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid session token",
		})
		return
	}

	session, found, err := d.Store.GetUserSession(ctx, sessionToken)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load session",
		})
		return
	}
	if !found {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid session token",
		})
		return
	}

	writeJSON(w, http.StatusOK, session)
}

func (d Dependencies) requireActor(w http.ResponseWriter, r *http.Request) (actor, bool) {
	sessionToken := sessionTokenFromRequest(r)
	guestToken := guestTokenFromRequest(r)
	if sessionToken != "" && guestToken != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "send either session token or guest token, not both",
		})
		return actor{}, false
	}

	if sessionToken != "" {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		ok, err := d.Store.TouchUserSession(ctx, sessionToken)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to validate session token",
			})
			return actor{}, false
		}
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid session token",
			})
			return actor{}, false
		}

		session, found, err := d.Store.GetUserSession(ctx, sessionToken)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to load session",
			})
			return actor{}, false
		}
		if !found {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "invalid session token",
			})
			return actor{}, false
		}

		return actor{
			kind:         "user",
			sessionToken: sessionToken,
			user:         session.User,
		}, true
	}

	if guestToken == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing auth token",
		})
		return actor{}, false
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ok, err := d.Store.TouchGuestSession(ctx, guestToken)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to validate guest token",
		})
		return actor{}, false
	}
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid guest token",
		})
		return actor{}, false
	}

	return actor{
		kind:       "guest",
		guestToken: guestToken,
	}, true
}

func sessionTokenFromRequest(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-Session-Token"))
}

func gameOwnedBy(currentActor actor, game domain.GameSetup) bool {
	if currentActor.kind == "user" {
		return currentActor.user.UserID != "" && game.UserID == currentActor.user.UserID
	}

	return currentActor.guestToken != "" && game.GuestToken == currentActor.guestToken
}
