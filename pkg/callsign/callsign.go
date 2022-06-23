package callsign

import (
	"errors"
	"github.com/k0swe/wsjtx-go/v4/pkg/city"
	"regexp"
	"strings"
)

type CallSign struct {
	Number string
	DXCC   *city.DXCCEntry
}

var cqDirectionRegExp = regexp.MustCompile("^[A-Z]{2}$")

func ExtractCallSignFromMessage(msg string, findDXCC bool) (de, dx *CallSign, err error) {

	parts := strings.Split(strings.TrimSpace(msg), " ")
	if len(parts) < 2 {
		return nil, nil, errors.New("format error")
	} else {
		// CQ
		if parts[0] == "CQ" {
			// CQ DX
			if cqDirectionRegExp.MatchString(parts[1]) {
				return buildCallSign(parts[2], true), nil, nil
			} else {
				return buildCallSign(parts[1], true), nil, nil
			}
		} else {
			return buildCallSign(parts[1], true), buildCallSign(parts[0], true), nil
		}
	}
}
func buildCallSign(call string, findDXCC bool) *CallSign {
	c := new(CallSign)
	c.Number = call
	if findDXCC {
		c.DXCC = city.FindDXCC(call)
	}
	return c
}
