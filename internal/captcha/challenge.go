package captcha

import (
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

func (mng *Manager) CreateChallenge() (altcha.Challenge, error) {
	expires := time.Now().Add(mng.cfg.Expiration)

	return altcha.CreateChallenge(altcha.ChallengeOptions{
		HMACKey: mng.cfg.HmacKey,
		Expires: &expires,
	})
}
