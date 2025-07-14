package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/pushed_auth_request/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// require interface implementation at compile time.
var _ repository.PushedAuthRequestCache = (*PushedAuthRequestCache)(nil)

type PushedAuthRequestCache struct {
	client *database.DB

	cache *lru.Cache[string, *oidc.AuthRequest]
}

func Start(dbClient *database.DB, amountOfCachedAuthRequests uint16) *PushedAuthRequestCache {
	cache := &PushedAuthRequestCache{
		client: dbClient,
	}

	requestURICache, err := lru.New[string, *oidc.AuthRequest](int(amountOfCachedAuthRequests))
	logging.OnError(err).Info("pushed auth request cache disabled")
	if err == nil {
		cache.cache = requestURICache
	}

	return cache
}

func (c *PushedAuthRequestCache) Health(ctx context.Context) error {
	return c.client.PingContext(ctx)
}

func (c *PushedAuthRequestCache) GetPushedAuthRequestByRequestURI(ctx context.Context, requestURI string) (*oidc.AuthRequest, error) {
	if authRequest, ok := c.getCachedByID(ctx, requestURI); ok {
		return authRequest, nil
	}

	request, err := c.getAuthRequest(ctx, requestURI, requestURI, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}

	c.CachePushedAuthRequest(ctx, requestURI, request)

	return request, nil
}

func (c *PushedAuthRequestCache) SavePushedAuthRequest(
	ctx context.Context, requestURI string, request *oidc.AuthRequest, expiresAt time.Time) error {
	b, err := json.Marshal(request)
	if err != nil {
		return zerrors.ThrowInternal(err, "CACHE-TfX8Eq", "Errors.Internal")
	}

	_, err = c.client.Exec(
		"INSERT INTO auth.pushed_auth_requests (request_uri, request, instance_id, expires_at) VALUES($1, $2, $3, $4)",
		requestURI, b, authz.GetInstance(ctx).InstanceID(), expiresAt)
	if err != nil {
		return zerrors.ThrowInternal(err, "CACHE-B4ZGaB", "Errors.Internal")
	}

	c.CachePushedAuthRequest(ctx, requestURI, request)

	return nil
}

func (c *PushedAuthRequestCache) DeletePushedAuthRequest(ctx context.Context, requestURI string) error {
	_, err := c.client.Exec("DELETE FROM auth.pushed_auth_requests WHERE instance_id = $1 and request_uri = $2", authz.GetInstance(ctx).InstanceID(), requestURI)
	if err != nil {
		return zerrors.ThrowInternal(err, "CACHE-J2mAJj", "unable to delete auth request")
	}

	c.deleteFromCache(ctx, requestURI)

	return nil
}

func (c *PushedAuthRequestCache) getAuthRequest(ctx context.Context, requestURI, value, instanceID string) (*oidc.AuthRequest, error) {
	var (
		requestData []byte
		expiresAt   time.Time
	)

	err := c.client.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			return row.Scan(&requestData, &expiresAt)
		},
		"SELECT request, expires_at FROM auth.pushed_auth_requests WHERE instance_id = $1 and request_uri = $2",
		instanceID, requestURI,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "CACHE-aMz1qe", "Errors.AuthRequest.NotFound")
		}

		return nil, zerrors.ThrowInternal(err, "CACHE-V8TPX2", "Errors.Internal")
	}

	if expiresAt.Before(time.Now()) {
		return nil, zerrors.ThrowNotFound(err, "CACHE-cA1OPE", "Errors.AuthRequest.NotFound")
	}

	var request oidc.AuthRequest
	err = json.Unmarshal(requestData, &request)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "CACHE-mvAP3X", "Errors.Internal")
	}

	return &request, nil
}

func (c *PushedAuthRequestCache) getCachedByID(ctx context.Context, requestID string) (*oidc.AuthRequest, bool) {
	if c.cache == nil {
		return nil, false
	}

	authRequest, ok := c.cache.Get(cacheKey(ctx, requestID))

	logging.WithFields("hit", ok).Info("get from pushed auth request cache")

	return authRequest, ok
}

func (c *PushedAuthRequestCache) CachePushedAuthRequest(ctx context.Context, requestURI string, request *oidc.AuthRequest) {
	if c.cache == nil {
		return
	}

	c.cache.Add(cacheKey(ctx, requestURI), request)
}

func cacheKey(ctx context.Context, value string) string {
	return fmt.Sprintf("%s-%s", authz.GetInstance(ctx).InstanceID(), value)
}

func (c *PushedAuthRequestCache) deleteFromCache(ctx context.Context, requestURI string) {
	if c.cache == nil {
		return
	}

	key := cacheKey(ctx, requestURI)

	_, ok := c.cache.Get(key)
	if !ok {
		return
	}

	c.cache.Remove(key)
}
