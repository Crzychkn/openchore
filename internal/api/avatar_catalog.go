package api

import (
	"bytes"
	"embed"
	"errors"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed avatarassets/*.png
var bundledAvatars embed.FS

const maxManagedAvatarBytes = 256 << 10

var errInvalidAvatar = errors.New("expected a PNG up to 256 KiB and 512 by 512 pixels")

var managedAvatarSlug = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type avatarAsset struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	image []byte
}

type avatarCatalog []avatarAsset

func loadAvatarCatalog(directory string) avatarCatalog {
	var catalog avatarCatalog
	for _, asset := range []avatarAsset{
		{ID: "fox", Label: "Fox"},
		{ID: "saturn", Label: "Saturn"},
		{ID: "sprout", Label: "Sprout"},
	} {
		asset.image, _ = bundledAvatars.ReadFile("avatarassets/" + asset.ID + ".png")
		catalog = append(catalog, asset)
	}
	if directory == "" {
		return catalog
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Printf("avatar catalog: %v", err)
		return catalog
	}
	for _, entry := range files {
		name := entry.Name()
		slug, isPNG := strings.CutSuffix(name, ".png")
		if !entry.Type().IsRegular() || !isPNG || len(slug) > 48 || !managedAvatarSlug.MatchString(slug) {
			continue
		}
		if len(catalog) >= 35 {
			break
		}
		image, err := readManagedAvatar(filepath.Join(directory, name))
		if err != nil {
			log.Printf("avatar catalog: cannot load %s: %v", name, err)
			continue
		}
		catalog = append(catalog, avatarAsset{
			ID: "managed-" + slug, Label: strings.ReplaceAll(slug, "-", " "), image: image,
		})
	}
	return catalog
}

func readManagedAvatar(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxManagedAvatarBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxManagedAvatarBytes {
		return nil, errInvalidAvatar
	}
	config, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width < 1 || config.Height < 1 || config.Width > 512 || config.Height > 512 {
		return nil, errInvalidAvatar
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var sanitized bytes.Buffer
	if err := png.Encode(&sanitized, img); err != nil {
		return nil, err
	}
	return sanitized.Bytes(), nil
}

func (catalog avatarCatalog) find(id string) *avatarAsset {
	for i := range catalog {
		if catalog[i].ID == id {
			return &catalog[i]
		}
	}
	return nil
}

func (catalog avatarCatalog) list(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, catalog)
}

func (catalog avatarCatalog) serveImage(w http.ResponseWriter, r *http.Request) {
	asset := catalog.find(chi.URLParam(r, "id"))
	if asset == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=300")
	_, _ = w.Write(asset.image)
}
