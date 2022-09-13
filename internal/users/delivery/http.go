package delivery

import (
	"net/http"
	"strconv"

	"github.com/Ollub/user_service/internal/session"
	"github.com/Ollub/user_service/internal/users"
	"github.com/Ollub/user_service/internal/users/usecase"
	"github.com/Ollub/user_service/pkg/log"
	"github.com/Ollub/user_service/pkg/utils/http_utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	sessions *session.SessionsJWTVer
	users    *usecase.Manager
}

func NewHandler(sessionManager *session.SessionsJWTVer, userManager *usecase.Manager) *Handler {
	return &Handler{sessionManager, userManager}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	loginReq, err := http_utils.FromBody[LoginReq](r)
	if err != nil {
		log.Clog(ctx).Error("Marshaling error", log.Fields{"err": err})
		http_utils.HttpError(w, "Provided payload can not be marshalled", http.StatusBadRequest)
		return
	}

	user, err := h.users.CheckPassByEmail(ctx, loginReq.Email, loginReq.Password)
	switch err {
	case nil:
		// all is ok
	case usecase.UserNotFoundError:
		http_utils.HttpError(w, "User not found", http.StatusNotFound)
	case usecase.BadPasswordError:
		http_utils.HttpError(w, "Wrong password provided", http.StatusBadRequest)
	default:
		log.Clog(ctx).Error("Error during checking user password", log.Fields{"err": err})
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
	}
	if err != nil {
		return
	}
	if err == usecase.UserNotFoundError {
		http_utils.HttpError(w, "User not found", http.StatusNotFound)
		return
	}
	if err == usecase.BadPasswordError {
		http_utils.HttpError(w, "Wrong password provided", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Clog(ctx).Error("Error during checking user password", log.Fields{"err": err})
		http_utils.HttpError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	token, err := h.sessions.Create(r.Context(), user)
	if err != nil {
		log.Clog(ctx).Error("Cant issue token", log.Fields{"err": err.Error()})
		http_utils.HttpError(w, "Internal error during token creation", http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, &LoginResp{token, user.ID}, http.StatusOK)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIn, err := http_utils.FromBody[users.UserIn](r)
	if err != nil {
		log.Clog(ctx).Error("Marshaling error", log.Fields{"err": err})
		http_utils.HttpError(w, "Provided payload can not be marshalled", http.StatusBadRequest)
		return
	}

	if err := validateUser(userIn); err != nil {
		http_utils.HttpError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	user, err := h.users.Create(ctx, userIn)
	if err == usecase.UserExistsError {
		http_utils.HttpError(w, "User exist", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Clog(ctx).Error("Unexpected error during user creation", log.Fields{"err": err.Error()})
		http_utils.HttpError(w, "Internal error during user creation", http.StatusInternalServerError)
		return
	}

	token, err := h.sessions.Create(r.Context(), user)
	if err != nil {
		log.Clog(ctx).Error("Cant issue token", log.Fields{"err": err.Error()})
		http_utils.HttpError(w, "Internal error during token creation", http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, &LoginResp{token, user.ID}, http.StatusCreated)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.users.ListUsers(r.Context())
	if err != nil {
		http_utils.HttpError(w, "Internal error while listing users", http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, ListUsersResp{items}, http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http_utils.HttpError(w, "Provided userId can not be converted to integer", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	sess := session.FromContext(ctx)
	if sess.UserID != uint32(userId) {
		http_utils.HttpError(w, "User can update only his profile", http.StatusForbidden)
		return
	}

	payload, err := http_utils.FromBody[users.UserUpdate](r)
	if err != nil {
		log.Clog(ctx).Error("Marshaling error", log.Fields{"err": err})
		http_utils.HttpError(w, "Provided payload can not be marshalled", http.StatusBadRequest)
		return
	}

	user, err := h.users.PartialUpdate(ctx, uint32(userId), payload)
	if err == usecase.UserNotFoundError {
		http_utils.HttpError(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Clog(ctx).Error("User update error", log.Fields{"userId": userId, "err": err.Error()})
		http_utils.HttpError(w, "Internal during user update", http.StatusInternalServerError)
		return
	}
	http_utils.JsonResp(w, user, http.StatusOK)
}
