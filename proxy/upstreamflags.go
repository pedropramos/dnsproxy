package proxy

import (
	"fmt"
	"strings"
)

// UpstreamFlags contains parsed flags for an upstream server.
type UpstreamFlags struct {
	// NoCache indicates whether responses from this upstream should be cached.
	NoCache bool
}

// Default UpstreamFlags (when no flags are passed)
func NewUpstreamFlags() UpstreamFlags {
    return UpstreamFlags{
        NoCache: false,
    }
}

func (o *UpstreamFlags) Merge(other UpstreamFlags) {
    // Logical OR: if either is true, the result is true
    o.NoCache = o.NoCache || other.NoCache
}

// upstreamFlagsSuffix is the delimiter that separates the upstream address
// from its flags.
const upstreamFlagsSuffix = "?"

// upstreamFlagsDelimiter is the delimiter between multiple flags.
const upstreamFlagsDelimiter = ","

// parseUpstreamFlags parses the flags suffix from an upstream address.
// It returns the clean address (without flags) and the parsed flags.
//
// Examples:
//   - "1.1.1.1?nocache" → address="1.1.1.1", flags={NoCache: true}
//   - "1.1.1.1" → address="1.1.1.1", flags={NoCache: false}
func parseUpstreamFlags(upstreamAddr string) (address string, flags UpstreamFlags, err error) {
	flags = NewUpstreamFlags()

	// Check if address has flags suffix
	if !strings.Contains(upstreamAddr, upstreamFlagsSuffix) {
		return upstreamAddr, flags, nil
	}

	// Split address and flags
	parts := strings.SplitN(upstreamAddr, upstreamFlagsSuffix, 2)
	if len(parts) != 2 {
		return "", flags, fmt.Errorf("invalid upstream format: %q", upstreamAddr)
	}

	address = parts[0]
	flagsStr := parts[1]

	if flagsStr == "" {
		return address, flags, nil
	}

	// Parse individual flags
	flagsList := strings.Split(flagsStr, upstreamFlagsDelimiter)
	for _, flagStr := range flagsList {
		flagStr = strings.TrimSpace(flagStr)
		if flagStr == "" {
			continue
		}

		// Parse the flag
		if err = flags.parseFlag(flagStr); err != nil {
			return "", flags, fmt.Errorf(
				"parsing flag %q: %w",
				flagStr,
				err,
			)
		}
	}

	return address, flags, nil
}

// parseFlag parses a single flag string and updates the UpstreamFlags.
func (o *UpstreamFlags) parseFlag(flagStr string) (err error) {
	switch strings.ToLower(flagStr) {
	case "nocache":
		o.NoCache = true
		return nil
	default:
		// Return error for unknown flags to prevent silent failures
		return fmt.Errorf("unknown upstream flag: %q", flagStr)
	}
}
