package vault

import (
	"fmt"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
	common "swagtask/internal/vault/common"
)

func HandlerGetVaults(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vaults, err := common.GetVaultsWithCollaborators(queries, utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	page := common.NewVaultsPage(vaults, true, user.PathToPfp, user.Username)
	templates.Render(w, "vaults-page", page)
}

func HandlerCreateVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vaults, err := common.CreateVault(queries, r.FormValue("vault_name"), r.FormValue("vault_description"), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "vault-card", vaults)
}

func HandlerDeleteVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	errDelete := common.DeleteVault(queries, utils.PgUUID(r.PathValue("vaultId")), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, errDelete) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}

func HandlerUpdateVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	var locked bool
	if r.FormValue("vault_locked") == "" {
		locked = false
	} else {
		locked = true
	}

	vault, errUpdate := common.UpdateVault(queries, utils.PgUUID(r.PathValue("vaultId")), utils.PgUUID(user.ID), r.FormValue("vault_name"), r.FormValue("vault_description"), locked, r.Context())
	if utils.CheckError(w, r, errUpdate) {
		fmt.Println("yea error here")
		return
	}

	templates.Render(w, "vault-card", vault)
}

func HandlerAddCollaboratorToVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vault, errUpdate := common.AddCollaboratorToVault(queries,
		r.FormValue("collaborator_username"),
		utils.PgUUID(user.ID),
		utils.PgUUID(r.PathValue("vaultId")),
		r.FormValue("collaborator_role"),
		r.Context())
	if utils.CheckError(w, r, errUpdate) {
		fmt.Println("yea error here")
		return
	}

	templates.Render(w, "vault-card", vault)
}

func HandlerDeleteCollaboratorToVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	collaboratorUsername := r.FormValue("collaborator_username")
	if collaboratorUsername == user.Username {
		http.Error(w, "You cant remove the owner as a collaborator, consider deleting the vault if needed.", http.StatusBadRequest)
		return
	}
	vault, errUpdate := common.RemoveCollaboratorFromVault(queries,
		utils.PgUUID(user.ID),
		utils.PgUUID(r.PathValue("vaultId")),
		collaboratorUsername,
		r.Context())
	if utils.CheckError(w, r, errUpdate) {
		fmt.Println("yea error here")
		return
	}

	templates.Render(w, "vault-card", vault)
}
