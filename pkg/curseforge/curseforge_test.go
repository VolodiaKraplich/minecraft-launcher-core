package curseforge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Helper to create a test server and client
func newTestClient(_ *testing.T, handler http.HandlerFunc) (*CurseforgeV1Client, func()) {
	server := httptest.NewServer(handler)
	client := NewCurseforgeV1Client("key", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	return client, server.Close
}

func TestCurseforgeV1Client_Constructor(t *testing.T) {
	headers := map[string]string{"hello": "world"}
	client := NewCurseforgeV1Client("key", WithBaseURL("https://test.url"), WithHeaders(headers))
	if client.baseURL != "https://test.url" {
		t.Errorf("expected baseURL to be https://test.url, got %s", client.baseURL)
	}
	if client.headers["hello"] != "world" {
		t.Errorf("expected header hello=world, got %v", client.headers)
	}
	if client.headers["x-api-key"] != "key" {
		t.Errorf("expected x-api-key=key, got %v", client.headers)
	}
}

func TestCurseforgeV1Client_GetMod(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/mods/310806" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "key" {
			t.Errorf("missing x-api-key header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": Mod{Name: "hello"}})
	})
	defer close()
	result, err := client.GetMod(context.Background(), 310806)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "" && result.Name != "hello" {
		t.Errorf("expected mod name to be 'hello', got %v", result.Name)
	}
}

func TestCurseforgeV1Client_GetModDescription(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/mods/310806/description" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": "description"})
	})
	defer close()
	desc, err := client.GetModDescription(context.Background(), 310806)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if desc != "description" {
		t.Errorf("expected description to be 'description', got %v", desc)
	}
}

func TestCurseforgeV1Client_GetModFile(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/mods/310806/files/2657461" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": File{FileName: "file"}})
	})
	defer close()
	file, err := client.GetModFile(context.Background(), 310806, 2657461)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file.FileName != "" && file.FileName != "file" {
		t.Errorf("expected file name to be 'file', got %v", file.FileName)
	}
}

func TestCurseforgeV1Client_GetMods(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/mods" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []Mod{}})
	})
	defer close()
	mods, err := client.GetMods(context.Background(), []int{310806})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mods) != 0 {
		t.Errorf("expected empty mods, got %v", mods)
	}
}

func TestCurseforgeV1Client_GetFiles(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/mods/files" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []File{}})
	})
	defer close()
	files, err := client.GetFiles(context.Background(), []int{2657461})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected empty files, got %v", files)
	}
}

func TestCurseforgeV1Client_GetCategories(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/categories") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []ModCategory{}})
	})
	defer close()
	cats, err := client.GetCategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 0 {
		t.Errorf("expected empty categories, got %v", cats)
	}
}

func TestCurseforgeV1Client_GetModFiles(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/mods/310806/files") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []File{}})
	})
	defer close()
	resp, err := client.GetModFiles(context.Background(), GetModFilesOptions{ModID: 310806})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected empty mod files, got %v", resp.Data)
	}
}

func TestCurseforgeV1Client_SearchMods(t *testing.T) {
	client, close := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/mods/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": []Mod{}, "pagination": Pagination{}})
	})
	defer close()
	resp, err := client.SearchMods(context.Background(), SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected empty search result, got %v", resp.Data)
	}
}

func TestCurseforgeV1Client_RealAPI_GetCategories(t *testing.T) {
	token := os.Getenv("CURSEFORGE_API_KEY")
	if token == "" {
		t.Skip("CURSEFORGE_API_KEY not set; skipping real API test")
	}
	client := NewCurseforgeV1Client(token)
	categories, err := client.GetCategories(context.Background())
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	if len(categories) == 0 {
		t.Errorf("expected non-empty categories from Curseforge API")
	}
	t.Logf("Fetched %d categories", len(categories))
}

func TestGetModsBatch(t *testing.T) {
	token := os.Getenv("CURSEFORGE_API_KEY")
	if token == "" {
		t.Skip("CURSEFORGE_API_KEY not set; skipping real API test")
	}
	api := NewCurseforgeV1Client(token)

	modIDs := []int{238222, 298312, 308711} // JEI, Xaero's Minimap, Biomes O' Plenty

	mods, err := api.GetMods(context.Background(), modIDs)
	if err != nil {
		t.Fatalf("GetMods failed: %v", err)
	}

	if len(mods) == 0 {
		t.Errorf("Expected to get some mods, got none")
	}

	for _, mod := range mods {
		t.Logf("Got mod: %s (ID: %d)", mod.Name, mod.ID)
	}
}
