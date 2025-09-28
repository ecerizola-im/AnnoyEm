package memes

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/ecerizola-im/AnnoyEm/internal/app/utils"
	"github.com/ecerizola-im/AnnoyEm/web/views/layouts"
	memesui "github.com/ecerizola-im/AnnoyEm/web/views/memes"
	"github.com/ecerizola-im/AnnoyEm/web/views/memes/components"
)

type MemeHandler struct {
	svc                   *MemeService
	MaxUploadBytes        int64
	MaxInMemoryBytes      int64
	ValidMemeContentTypes []string
}

func NewHandler(svc *MemeService) *MemeHandler {
	return &MemeHandler{
		svc:              svc,
		MaxUploadBytes:   10 << 20, // 10MB default
		MaxInMemoryBytes: 32 << 20, // 32MB default for multipart form parsing
		ValidMemeContentTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"image/bmp",
		},
	}
}

func (handler *MemeHandler) Register(mux *http.ServeMux) {
	// POST /receipts -> upload a new receipt
	mux.HandleFunc("GET /memes", handler.getMemes)
	mux.HandleFunc("POST /memes", handler.uploadMeme)
	mux.HandleFunc("GET /memes/{id}/download", handler.downloadMeme) // download memes
	mux.HandleFunc("/", handler.handleBase)
}

func (handler *MemeHandler) handleBase(rw http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}

	http.Redirect(rw, r, "/memes", http.StatusSeeOther)
}

func (handler *MemeHandler) downloadMeme(rw http.ResponseWriter, r *http.Request) {

	memeID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		utils.WriteError(rw, http.StatusBadRequest, fmt.Errorf("invalid meme id: %v", err))
		return
	}

	meme, err := handler.svc.GetMeme(r.Context(), memeID)
	if err != nil {
		utils.WriteError(rw, http.StatusNotFound, fmt.Errorf("meme not found: %v", err))
		return
	}

	file, err := handler.svc.GetMemeFile(r.Context(), *meme.UUID)

	if err != nil {
		utils.WriteError(rw, http.StatusInternalServerError, fmt.Errorf("failed to retrieve meme file"))
		return
	}

	defer file.Close()

	rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", meme.OriginalFileName))
	rw.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(rw, file); err != nil {
		utils.WriteError(rw, http.StatusInternalServerError, fmt.Errorf("failed to send meme file: %v", err))
		return
	}

}

func (handler *MemeHandler) uploadMeme(rw http.ResponseWriter, r *http.Request) {

	defer io.Copy(io.Discard, r.Body)
	r.Body = http.MaxBytesReader(rw, r.Body, handler.MaxUploadBytes)

	if r.ContentLength > 0 && r.ContentLength > handler.MaxUploadBytes {
		utils.WriteError(rw, http.StatusRequestEntityTooLarge, fmt.Errorf("file too big: max %d MB", handler.MaxUploadBytes/(1<<20)))
		return
	}

	if err := r.ParseMultipartForm(handler.MaxInMemoryBytes); err != nil {
		utils.WriteError(rw, http.StatusBadRequest, fmt.Errorf("invalid multipart form: %v", err))
		return
	}

	file, header, err := r.FormFile("meme")

	if err != nil {
		utils.WriteError(rw, http.StatusBadRequest, fmt.Errorf("file field 'meme' is required"))
		return
	}

	defer file.Close()

	if err := handler.validateFileMetadata(header, file); err != nil {
		utils.WriteError(rw, http.StatusBadRequest, err)
		return
	}

	origName := strings.TrimSpace(header.Filename)

	memeId, svcErr := handler.svc.AddMeme(r.Context(), origName, file)
	if svcErr != nil {
		http.Error(rw, svcErr.Error(), http.StatusInternalServerError)
		return
	}

	addedMeme, err := handler.svc.GetMeme(r.Context(), memeId)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to retrieve added meme: %v", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	var memeComponent = components.MemeRow(*addedMeme)

	if err := memeComponent.Render(r.Context(), rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (handler *MemeHandler) getMemes(rw http.ResponseWriter, r *http.Request) {

	memes, err := handler.svc.GetMemes(r.Context())

	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get memes: %v", err), http.StatusInternalServerError)
		return
	}

	memeListComponent := components.MemeList(memes)
	bodyComponent := memesui.Memes("Upload Meme", memeListComponent)

	layoutComponent := layouts.Layout("EICS Tracker", bodyComponent)

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := layoutComponent.Render(r.Context(), rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *MemeHandler) validateFileMetadata(fileHeader *multipart.FileHeader, file multipart.File) error {

	var errors []error
	if err := handler.validateMemeFileHasName(fileHeader); err != nil {
		errors = append(errors, err)
	}

	if err := handler.validateMemeFileType(file); err != nil {
		errors = append(errors, err)
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("invalid receipt file: %v", errors)
}

func (handler *MemeHandler) validateMemeFileHasName(fileHeader *multipart.FileHeader) error {
	if strings.TrimSpace(fileHeader.Filename) == "" {
		return fmt.Errorf("file name is empty")
	}
	return nil
}

func (handler *MemeHandler) validateMemeFileType(file multipart.File) error {

	sniffBuf := make([]byte, 512)
	_, err := file.Read(sniffBuf)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(sniffBuf)
	if !(slices.Contains(handler.ValidMemeContentTypes, contentType)) {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}

	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to reset file reader: %w", err)
		}
	}

	return nil
}
