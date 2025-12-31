package captcha

import (
	"context"
	"strconv"
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

func (mng *Manager) VerifySolution(ctx context.Context, payload altcha.Payload) (bool, error) {
	// Check the solution and early return if it's wrong
	ok, err := altcha.VerifySolutionSafe(payload, mng.cfg.HmacKey, true)

	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	// If it's valid, avoid reusing the challenge
	used, err := mng.data.IsUsed(ctx, payload.Challenge)

	if err != nil {
		return false, err
	}

	if used {
		return false, nil
	}

	// Get challenge expiration time
	params := altcha.ExtractParams(payload)
	expires := params.Get("expires")

	if expires == "" {
		return false, nil
	}

	expireTime, err := strconv.ParseInt(expires, 10, 64) // it's not in the past because we already checked it

	if err != nil {
		return false, err
	}

	// The challenge is OK and it's never been used
	ttl := time.Until(time.Unix(expireTime, 0)) + time.Minute // adding a safety margin

	err = mng.data.MarkUsed(ctx, payload.Challenge, ttl)

	if err != nil {
		return false, err
	}

	return true, nil
}
