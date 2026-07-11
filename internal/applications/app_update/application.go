package app_update

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"
)

// ─── Application ─────────────────────────────────────────────────────────────

type AppUpdateApplication struct {
	db            *gorm.DB
	errorHandler  http_errors.HttpErrorHandler
	webhookSecret string
	githubToken   string
	s3Client      *s3.Client
	r2Bucket      string
	r2PublicURL   string
}

func NewAppUpdateApplication(
	db *gorm.DB,
	errorHandler http_errors.HttpErrorHandler,
	webhookSecret string,
	githubToken string,
	r2AccountId string,
	r2AccessKeyId string,
	r2SecretAccessKey string,
	r2Bucket string,
	r2PublicURL string,
) *AppUpdateApplication {
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2AccountId)
	options := s3.Options{
		Region:       "auto",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(r2AccessKeyId, r2SecretAccessKey, "")),
		BaseEndpoint: aws.String(r2Endpoint),
	}
	s3Client := s3.New(options)
	return &AppUpdateApplication{
		db:            db,
		errorHandler:  errorHandler,
		webhookSecret: webhookSecret,
		githubToken:   githubToken,
		s3Client:      s3Client,
		r2Bucket:      r2Bucket,
		r2PublicURL:   r2PublicURL,
	}
}

func (a *AppUpdateApplication) PublicRoutes(r chi.Router) {
	r.Get("/app/latest", endpoint.JSON(a.errorHandler, a.latest))
	r.Post("/webhooks/github", a.webhook)
}

// ─── DB model ────────────────────────────────────────────────────────────────

type appRelease struct {
	Id          int64     `gorm:"primaryKey;autoIncrement"`
	Platform    string    `gorm:"not null"`
	Version     string    `gorm:"not null"`
	VersionCode int64     `gorm:"column:version_code;not null"`
	ApkURL      string    `gorm:"column:apk_url;not null"`
	Notes       string    `gorm:"not null;default:''"`
	Mandatory   bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

func (appRelease) TableName() string { return "app_releases" }

// ─── GET /app/latest ─────────────────────────────────────────────────────────

type latestResponse struct {
	Version   string `json:"version"`
	ApkURL    string `json:"apk_url"`
	Notes     string `json:"notes,omitempty"`
	Mandatory bool   `json:"mandatory,omitempty"`
}

func (a *AppUpdateApplication) latest(r *http.Request) (*latestResponse, int, error) {
	platform := r.URL.Query().Get("platform")
	if platform == "" {
		platform = "android"
	}
	var rel appRelease
	err := a.db.Where("platform = ?", platform).Order("version_code DESC").First(&rel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusNoContent, nil
	}
	if err != nil {
		return nil, 0, cerr.NewInternalError("fetching latest release", err)
	}
	response := &latestResponse{
		Version:   rel.Version,
		ApkURL:    rel.ApkURL,
		Notes:     rel.Notes,
		Mandatory: rel.Mandatory,
	}
	return response, http.StatusOK, nil
}

// ─── POST /webhooks/github ────────────────────────────────────────────────────

type githubReleasePayload struct {
	Action  string `json:"action"`
	Release struct {
		TagName    string `json:"tag_name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		Assets     []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Size               int64  `json:"size"`
		} `json:"assets"`
	} `json:"release"`
}

func (a *AppUpdateApplication) webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20))
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	if !validHMAC(body, r.Header.Get("X-Hub-Signature-256"), a.webhookSecret) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("X-GitHub-Event") != "release" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var p githubReleasePayload
	err = json.Unmarshal(body, &p)
	if err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if p.Action != "published" || p.Release.Draft || p.Release.Prerelease {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if !strings.HasPrefix(p.Release.TagName, "build-v") {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	version := strings.TrimPrefix(p.Release.TagName, "build-v")

	var downloadURL string
	for _, asset := range p.Release.Assets {
		if strings.HasSuffix(strings.ToLower(asset.Name), ".apk") {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		http.Error(w, "no apk asset", http.StatusUnprocessableEntity)
		return
	}

	versionCode, err := parseVersionCode(version)
	if err != nil {
		http.Error(w, "invalid version: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	apkURL, err := a.downloadAndUpload(r.Context(), downloadURL, version)
	if err != nil {
		slog.Error("app_update: download/upload failed", "version", version, "error", err)
		http.Error(w, "failed to process APK", http.StatusInternalServerError)
		return
	}

	err = a.upsert("android", version, versionCode, apkURL, firstLine(p.Release.Body))
	if err != nil {
		slog.Error("app_update: upsert failed", "version", version, "error", err)
		http.Error(w, "failed to save release", http.StatusInternalServerError)
		return
	}

	slog.Info("app_update: release published", "version", version, "apk_url", apkURL)
	w.WriteHeader(http.StatusAccepted)
}

// ─── R2 upload ───────────────────────────────────────────────────────────────

func (a *AppUpdateApplication) downloadAndUpload(ctx context.Context, downloadURL, version string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	if a.githubToken != "" {
		request.Header.Set("Authorization", "token "+a.githubToken)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("downloading APK: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub returned %d", response.StatusCode)
	}

	key := fmt.Sprintf("releases/android/vitae-%s.apk", version)
	contentType := "application/vnd.android.package-archive"

	input := &s3.PutObjectInput{
		Bucket:      aws.String(a.r2Bucket),
		Key:         aws.String(key),
		Body:        response.Body,
		ContentType: aws.String(contentType),
	}
	if response.ContentLength > 0 {
		input.ContentLength = aws.Int64(response.ContentLength)
	}

	_, err = a.s3Client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("uploading to R2: %w", err)
	}

	return strings.TrimRight(a.r2PublicURL, "/") + "/" + key, nil
}

// ─── DB upsert ───────────────────────────────────────────────────────────────

func (a *AppUpdateApplication) upsert(platform, version string, versionCode int64, apkURL, notes string) error {
	rel := appRelease{
		Platform:    platform,
		Version:     version,
		VersionCode: versionCode,
		ApkURL:      apkURL,
		Notes:       notes,
	}
	onConflict := clause.OnConflict{
		Columns:   []clause.Column{{Name: "platform"}, {Name: "version"}},
		DoUpdates: clause.AssignmentColumns([]string{"apk_url", "notes", "version_code"}),
	}
	err := a.db.Clauses(onConflict).Create(&rel).Error
	if err != nil {
		return err
	}
	return nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func validHMAC(body []byte, header, secret string) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(header, prefix) || secret == "" {
		return false
	}
	want, err := hex.DecodeString(strings.TrimPrefix(header, prefix))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hmac.Equal(want, mac.Sum(nil))
}

func parseVersionCode(version string) (int64, error) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) != 3 {
		return 0, fmt.Errorf("expected X.Y.Z, got %q", version)
	}
	major, err1 := strconv.ParseInt(parts[0], 10, 64)
	minor, err2 := strconv.ParseInt(parts[1], 10, 64)
	patch, err3 := strconv.ParseInt(parts[2], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, fmt.Errorf("non-numeric segment in %q", version)
	}
	return major*1_000_000 + minor*1_000 + patch, nil
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	i := strings.IndexByte(s, '\n')
	if i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}
