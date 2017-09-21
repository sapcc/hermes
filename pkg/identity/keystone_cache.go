package identity

import (
	"container/list"
	"github.com/gophercloud/gophercloud"
	"github.com/sapcc/hermes/pkg/util"
	"sync"
	"time"
)

// Cache type used for the name caches
type cache struct {
	// "Inherit from" sync.RWMutex, so the cache can be locked during access/update
	sync.RWMutex
	// The actual cache - a simple map with no expiry
	//  (the total number of items will only be in the 10000s, ~100 bytes per item, so ~1Mb per cache)
	m map[string]string
}

var providerClient *gophercloud.ProviderClient
var domainNameCache *cache
var projectNameCache *cache
var userNameCache *cache
var userIdCache *cache
var roleNameCache *cache
var groupNameCache *cache

// Token cache
type keystoneTokenCache struct {
	// "Inherit from" sync.RWMutex, so the cache can be locked during access/update
	sync.RWMutex
	//  tMap:  Cached tokens (keystoneToken struct) accessible by the token ID from the request
	tMap map[string]*keystoneToken // map tokenID to token struct
	//  eList: A sorted list of token expiry times, so we don't have to scan the whole list
	//         every time we check to see what's expired
	eList *list.List // sorted list of expiration times
	//  eMap:  If we know a token is expired at time T, we use this map to look up the tokenID
	//         so we can then remove the token from tMap.
	eMap map[time.Time][]string // map expiration time to list of tokenIDs
}

var tokenCache *keystoneTokenCache

func init() {
	domainNameCache = &cache{m: make(map[string]string)}
	projectNameCache = &cache{m: make(map[string]string)}
	userNameCache = &cache{m: make(map[string]string)}
	userIdCache = &cache{m: make(map[string]string)}
	roleNameCache = &cache{m: make(map[string]string)}
	groupNameCache = &cache{m: make(map[string]string)}
	tokenCache = &keystoneTokenCache{
		tMap:  make(map[string]*keystoneToken),
		eMap:  make(map[time.Time][]string),
		eList: list.New(),
	}
}

func updateCache(cache *cache, key string, value string) {
	cache.Lock()
	cache.m[key] = value
	cache.Unlock()
}

func getFromCache(cache *cache, key string) (string, bool) {
	cache.RLock()
	value, exists := cache.m[key]
	cache.RUnlock()
	return value, exists
}

func addTokenToCache(cache *keystoneTokenCache, id string, token *keystoneToken) {
	expiryTime, err := time.Parse("2006-01-02T15:04:05.999999Z", token.ExpiresAt)
	if err != nil {
		util.LogWarning("Not adding token to cache because time '%s' could not be parsed", token.ExpiresAt)
		return
	}
	cache.Lock()
	cache.eList.PushBack(expiryTime)
	// If the expiryTime is earlier than the last item in the list,
	// move it to the correct place so that we keep the list sorted
	lastItem := cache.eList.Back()
	if cache.eList.Back() != nil && expiryTime.Before(lastItem.Value.(time.Time)) {
		for e := cache.eList.Back(); e != nil; e = e.Prev() {
			if expiryTime.After((e.Value).(time.Time)) {
				cache.eList.MoveAfter(cache.eList.Back(), e)
			}
		}
	}
	if cache.eMap[expiryTime] == nil {
		cache.eMap[expiryTime] = []string{id}
	} else {
		cache.eMap[expiryTime] = append(cache.eMap[expiryTime], id)
	}
	cache.tMap[id] = token
	cacheSize := len(cache.tMap)
	cache.Unlock()
	util.LogDebug("Added token to cache. Current cache size: %d", cacheSize)
}

func getCachedToken(cache *keystoneTokenCache, id string) *keystoneToken {
	// First, remove expired tokens from cache
	now := time.Now()
	elemsToRemove := []*list.Element{}
	cache.RLock()
	for e := cache.eList.Front(); e != nil; e = e.Next() {
		expiryTime := (e.Value).(time.Time)
		if now.Before(expiryTime) {
			break // list is sorted, so we can stop once we get to an unexpired token
		}
		// We can't remove from the list during the for loop, so remember which ones to delete
		elemsToRemove = append(elemsToRemove, e)
	}
	cache.RUnlock()
	cache.Lock()
	for _, elem := range elemsToRemove {
		cache.eList.Remove(elem) // Remove the cached expiry time from the sorted list
		time := (elem.Value).(time.Time)
		tokenIds := cache.eMap[time]
		delete(cache.eMap, time) // Remove the cached expiry time from the time:tokenIDs map
		for _, tokenId := range tokenIds {
			delete(cache.tMap, tokenId) // Remove all the cached tokens
		}
	}
	cacheSize := len(cache.tMap)
	cache.Unlock()
	if len(elemsToRemove) > 0 {
		util.LogDebug("Removed expired token(s) from cache. Current cache size: %d", cacheSize)
	}
	// Now look for the token in question
	cache.RLock()
	token := cache.tMap[id]
	cache.RUnlock()
	if token != nil {
		util.LogDebug("Got token from cache. Current cache size: %d", cacheSize)
	}

	return token
}
