package api_test

import (
	"fmt"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/liftedkilt/openchore/internal/model"
)

func TestCuratedAvatarLifecycle(t *testing.T) {
	env := setupTest(t)
	env.createAdmin(t)
	id := env.createChild(t, "Avatar kid")
	other := env.createChild(t, "Other kid")
	profile := fmt.Sprintf("/api/users/%d", id)
	path := profile + "/avatar"
	session := childHeaders(id)
	picture := map[string]string{"avatar_asset_id": "fox"}
	legacy := "https://example.com/old.png"

	image := env.expectStatus(t, "GET", "/api/avatar-catalog/fox", nil, nil, http.StatusOK)
	defer image.Body.Close()
	if got := image.Header.Get("Content-Type"); got != "image/png" {
		t.Fatalf("image type = %q, want image/png", got)
	}
	if _, err := png.Decode(image.Body); err != nil {
		t.Fatal(err)
	}

	env.expectStatus(t, "PUT", path, picture, nil, http.StatusUnauthorized).Body.Close()
	env.expectStatus(t, "PUT", path, picture, childHeaders(other), http.StatusForbidden).Body.Close()
	env.expectStatus(t, "PUT", path, map[string]string{"avatar_asset_id": "unknown"}, session, http.StatusBadRequest).Body.Close()
	env.expectStatus(t, "PUT", path, map[string]string{"avatar_url": legacy}, session, http.StatusOK).Body.Close()
	env.expectStatus(t, "PUT", path, picture, session, http.StatusOK).Body.Close()
	env.expectStatus(t, "PUT", profile, map[string]string{"name": "Renamed"}, adminHeaders(), http.StatusOK).Body.Close()
	env.expectStatus(t, "PUT", profile+"/theme", map[string]string{"theme": "tint"}, session, http.StatusOK).Body.Close()

	var user model.User
	decodeBody(t, env.expectStatus(t, "GET", profile, nil, nil, http.StatusOK), &user)
	if user.AvatarAssetID != "fox" || user.AvatarURL != legacy {
		t.Fatalf("catalog selection and legacy fallback after edits = %q, %q", user.AvatarAssetID, user.AvatarURL)
	}

	decodeBody(t, env.expectStatus(t, "PUT", path, map[string]string{"avatar_url": legacy}, session, http.StatusOK), &user)
	if user.AvatarAssetID != "" || user.AvatarURL != legacy {
		t.Fatalf("legacy selection = %q, %q", user.AvatarAssetID, user.AvatarURL)
	}

	env.expectStatus(t, "PUT", path, picture, session, http.StatusOK).Body.Close()
	decodeBody(t, env.expectStatus(t, "PUT", path, map[string]string{"avatar_asset_id": ""}, session, http.StatusOK), &user)
	if user.AvatarAssetID != "" || user.AvatarURL != "" {
		t.Fatalf("initials selection = %q, %q", user.AvatarAssetID, user.AvatarURL)
	}
}

func TestManagedAvatarCatalog(t *testing.T) {
	dir := t.TempDir()
	data, err := os.ReadFile("avatarassets/fox.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "forest-fox.png"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.png"), []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("OPENCHORE_AVATAR_DIR", dir)
	env := setupTest(t)
	var assets []struct{ ID, Label string }
	decodeBody(t, env.expectStatus(t, "GET", "/api/avatar-catalog", nil, nil, http.StatusOK), &assets)
	if len(assets) != 4 || assets[3].ID != "managed-forest-fox" || assets[3].Label != "forest fox" {
		t.Fatalf("catalog = %#v", assets)
	}
	env.expectStatus(t, "GET", "/api/avatar-catalog/managed-forest-fox", nil, nil, http.StatusOK).Body.Close()
	env.expectStatus(t, "GET", "/api/avatar-catalog/managed-bad", nil, nil, http.StatusNotFound).Body.Close()
}
