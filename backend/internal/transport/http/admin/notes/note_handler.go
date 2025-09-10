package notes

import (
	"strconv"

	"github.com/ArtisanCloud/PowerXPlugin/internal/contracts"
	srvnotes "github.com/ArtisanCloud/PowerXPlugin/internal/services/admin/notes"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

type NoteHandler struct{ NoteService *srvnotes.NoteService }

func NewNoteHandler(deps *app.Deps) *NoteHandler {
	return &NoteHandler{NoteService: srvnotes.NewNoteService(deps.DB)}
}

// ====== Handlers ======

func (h *NoteHandler) GetNotes(c *gin.Context) {
	var q NoteListRequest
	if err := c.ShouldBindQuery(&q); err != nil {
		contracts.ResponseBadRequest(c, "invalid query: "+err.Error())
		return
	}

	page, pageSize := q.Page, q.PageSize
	res, err := h.NoteService.List(c.Request.Context(), q.Q, q.TeamID, q.MemberID, page, pageSize)
	if err != nil {
		contracts.ResponseInternalError(c, err)
		return
	}
	contracts.ResponseSuccess(c, res)
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		contracts.ResponseBadRequest(c, "invalid id")
		return
	}
	n, err := h.NoteService.GetByID(c.Request.Context(), id)
	if err != nil {
		contracts.ResponseNotFound(c, "not found: "+err.Error())
		return
	}
	contracts.ResponseSuccess(c, n)
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		contracts.ResponseBadRequest(c, "invalid body: "+err.Error())
		return
	}
	n, err := h.NoteService.Create(c.Request.Context(), req.Title, req.Content, req.Author, req.TeamID, req.MemberID)
	if err != nil {
		contracts.ResponseInternalError(c, err)
		return
	}
	contracts.ResponseSuccess(c, n)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		contracts.ResponseBadRequest(c, "invalid id")
		return
	}
	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		contracts.ResponseBadRequest(c, "invalid body: "+err.Error())
		return
	}
	n, err := h.NoteService.Update(c.Request.Context(), id, req.Title, req.Content, req.Author, req.TeamID, req.MemberID)
	if err != nil {
		contracts.ResponseInternalError(c, err)
		return
	}
	contracts.ResponseSuccess(c, n)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		contracts.ResponseBadRequest(c, "invalid id")
		return
	}
	if err := h.NoteService.Delete(c.Request.Context(), id); err != nil {
		contracts.ResponseInternalError(c, err)
		return
	}
	contracts.ResponseSuccess(c, gin.H{"ok": true})
}

func parseUint64(s string) (uint64, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	return uint64(u), err
}
