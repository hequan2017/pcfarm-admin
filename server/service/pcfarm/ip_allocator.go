package pcfarm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"

	pcfarmModel "github.com/flipped-aurora/gin-vue-admin/server/model/pcfarm"
)

var (
	ErrIPPoolExhausted = errors.New("ip pool exhausted")
	ErrInvalidIPRange  = errors.New("invalid ip range")
)

func nextAvailableIP(pool pcfarmModel.IPPool, allocated map[string]struct{}) (string, error) {
	start, err := parseIPv4(pool.StartIP)
	if err != nil {
		return "", err
	}
	end, err := parseIPv4(pool.EndIP)
	if err != nil {
		return "", err
	}
	if start > end {
		return "", fmt.Errorf("%w: start ip is after end ip", ErrInvalidIPRange)
	}

	for current := start; current <= end; current++ {
		ip := uint32ToIPv4(current)
		if _, ok := allocated[ip]; !ok {
			return ip, nil
		}
		if current == ^uint32(0) {
			break
		}
	}
	return "", ErrIPPoolExhausted
}

func parseIPv4(raw string) (uint32, error) {
	addr, err := netip.ParseAddr(raw)
	if err != nil || !addr.Is4() {
		return 0, fmt.Errorf("%w: %s", ErrInvalidIPRange, raw)
	}
	bytes := addr.As4()
	return binary.BigEndian.Uint32(bytes[:]), nil
}

func uint32ToIPv4(value uint32) string {
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[:], value)
	return netip.AddrFrom4(bytes).String()
}
