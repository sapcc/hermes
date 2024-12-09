/*******************************************************************************
*
* Copyright 2022 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package identity

import (
	"container/list"
	"sync"
	"time"

	"github.com/gophercloud/gophercloud/v2"

	"github.com/sapcc/go-bits/logg"
)

var providerClient *gophercloud.ProviderClient

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
	tokenCache = &keystoneTokenCache{
		tMap:  make(map[string]*keystoneToken),
		eMap:  make(map[time.Time][]string),
		eList: list.New(),
	}
}

func addTokenToCache(cache *keystoneTokenCache, id string, token *keystoneToken) {
	expiryTime, err := time.Parse("2006-01-02T15:04:05.999999Z", token.ExpiresAt)
	if err != nil {
		logg.Error("Not adding token to cache because time '%s' could not be parsed", token.ExpiresAt)
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
	logg.Debug("Added token to cache. Current cache size: %d", cacheSize)
}

func getCachedToken(cache *keystoneTokenCache, id string) *keystoneToken {
	// First, remove expired tokens from cache
	now := time.Now()
	elemsToRemove := []*list.Element{}
	cache.RLock()
	for e := cache.eList.Front(); e != nil; e = e.Next() {
		expiryTime := (e.Value).(time.Time) //nolint:errcheck
		if now.Before(expiryTime) {
			break // list is sorted, so we can stop once we get to an unexpired token
		}
		// We can't remove from the list during the for loop, so remember which ones to delete
		elemsToRemove = append(elemsToRemove, e)
	}
	cache.RUnlock()
	cache.Lock()
	for _, elem := range elemsToRemove {
		cache.eList.Remove(elem)                 // Remove the cached expiry time from the sorted list
		timeToRemove := (elem.Value).(time.Time) //nolint:errcheck
		tokenIDs := cache.eMap[timeToRemove]
		delete(cache.eMap, timeToRemove) // Remove the cached expiry time from the time:tokenIDs map
		for _, tokenID := range tokenIDs {
			delete(cache.tMap, tokenID) // Remove all the cached tokens
		}
	}
	cacheSize := len(cache.tMap)
	cache.Unlock()
	if len(elemsToRemove) > 0 {
		logg.Debug("Removed expired token(s) from cache. Current cache size: %d", cacheSize)
	}
	// Now look for the token in question
	cache.RLock()
	token := cache.tMap[id]
	cache.RUnlock()
	if token != nil {
		logg.Debug("Got token from cache. Current cache size: %d", cacheSize)
	}

	return token
}
