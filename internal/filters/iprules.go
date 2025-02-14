package filters

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/ma111e/melody/internal/logging"
)

// IPRanges abstracts an array of IPRange
type IPRanges []IPRange

// IPRules groups the whitelisted and blacklisted ip rules
type IPRules struct {
	WhitelistedIPs IPRanges
	BlacklistedIPs IPRanges
}

// IPRange is a range of IP represented by a lower and an upper bound
type IPRange struct {
	Lower net.IP
	Upper net.IP
}

// NewIPRange created a new ip range from a lower and an upper bound
func NewIPRange(lower net.IP, upper net.IP) IPRange {
	return IPRange{
		Lower: lower,
		Upper: upper,
	}
}

//
//func (iprl *IPRules) ParseRules(rules []string) {
//	for _, rawRule := range rules {
//		rule := strings.Replace(rawRule, " ", "", -1)
//
//		if strings.HasPrefix(rawRule, "not") {
//			rule = strings.TrimPrefix(rule, "not")
//
//			if strings.Contains(rule, "-") {
//				err := iprl.BlacklistRange(rule)
//				if err != nil {
//					log.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
//					log.Println(err)
//					os.Exit(1)
//				}
//				continue
//			} else if strings.Contains(rule, "/") {
//				err := iprl.BlacklistCIDR(rule)
//				if err != nil {
//					log.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
//					log.Println(err)
//					os.Exit(1)
//				}
//				continue
//			}
//
//			err := iprl.Blacklist(rule)
//			if err != nil {
//				log.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
//				//log.Println(err)
//				os.Exit(1)
//			}
//			continue
//		}
//
//		if strings.Contains(rule, "-") {
//			err := iprl.WhitelistRange(rule)
//			if err != nil {
//				log.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
//				log.Println(err)
//				os.Exit(1)
//			}
//		} else if strings.Contains(rule, "/") {
//			err := iprl.WhitelistCIDR(rule)
//			if err != nil {
//				log.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
//				//log.Println(err)
//				os.Exit(1)
//			}
//			continue
//		} else {
//			err := iprl.Whitelist(rule)
//			if err != nil {
//				log.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
//				log.Println(err)
//				os.Exit(1)
//			}
//		}
//	}
//
//	iprl.BlacklistedIPs.MergeOverlapping()
//	iprl.WhitelistedIPs.MergeOverlapping()
//}

// ParseRules loads a whitelist and a blacklist into a set of IPRules
func (iprl *IPRules) ParseRules(whitelist []string, blacklist []string) {
	for _, rawRule := range whitelist {
		rule := strings.Replace(rawRule, " ", "", -1)

		if strings.Contains(rule, "-") {
			err := iprl.WhitelistRange(rule)
			if err != nil {
				logging.Errors.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
				logging.Errors.Println(err)
				os.Exit(1)
			}
		} else if strings.Contains(rule, "/") {
			err := iprl.WhitelistCIDR(rule)
			if err != nil {
				logging.Errors.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
				logging.Errors.Println(err)
				os.Exit(1)
			}
			continue
		} else {
			err := iprl.Whitelist(rule)
			if err != nil {
				logging.Errors.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
				logging.Errors.Println(err)
				os.Exit(1)
			}
		}
	}

	for _, rawRule := range blacklist {
		rule := strings.Replace(rawRule, " ", "", -1)

		if strings.Contains(rule, "-") {
			err := iprl.BlacklistRange(rule)
			if err != nil {
				logging.Errors.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
				logging.Errors.Println(err)
				os.Exit(1)
			}
			continue
		} else if strings.Contains(rule, "/") {
			err := iprl.BlacklistCIDR(rule)
			if err != nil {
				logging.Errors.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
				logging.Errors.Println(err)
				os.Exit(1)
			}
			continue
		}

		err := iprl.Blacklist(rule)
		if err != nil {
			logging.Errors.Println(fmt.Sprintf("Failed to parse the IP rule [%s]:", rule))
			logging.Errors.Println(err)
			os.Exit(1)
		}
		continue
	}

	iprl.BlacklistedIPs.MergeOverlapping()
	iprl.WhitelistedIPs.MergeOverlapping()
}

//
// IPRanges methods
//

// MergeOverlapping optimize the parsed IPRange by keeping only non-overlapping ranges
func (iprgs *IPRanges) MergeOverlapping() {
	workSlice := make(IPRanges, len(*iprgs))
	copy(workSlice, *iprgs)

	for i := 0; i < len(workSlice); i++ {
		for idx, candidate := range workSlice {
			if candidate.Equals(workSlice[i]) {
				// Skip
				continue
			}

			if candidate.ContainsIPRange(workSlice[i]) {
				workSlice.RemoveAt(i)
				i = 0 // Restart upper loop
				break
			}

			if workSlice[i].ContainsIPRange(candidate) {
				workSlice.RemoveAt(idx)
				i = 0 // Restart upper loop
				break
			}

			if !candidate.IsUpperOrLowerBoundary(workSlice[i].Lower) && candidate.IsUpperOrLowerBoundary(workSlice[i].Upper) {
				// Replace the candidate's upper with the current's upper
				workSlice[idx].Upper = workSlice[i].Upper
				workSlice.RemoveAt(i)
				i = 0 // Restart upper loop
				break
			}

			if !candidate.IsUpperOrLowerBoundary(workSlice[i].Upper) && candidate.IsUpperOrLowerBoundary(workSlice[i].Lower) {
				// Replace the candidate's lower with the current's lower
				workSlice[idx].Lower = workSlice[i].Lower
				workSlice.RemoveAt(i)
				i = 0 // Restart upper loop
				break
			}
		}
	}

	*iprgs = workSlice
}

// RemoveAt is an helper that removes a range at the the given index
func (iprgs *IPRanges) RemoveAt(index int) {
	workSlice := make(IPRanges, len(*iprgs))
	copy(workSlice, *iprgs)

	workSlice = append(workSlice[:index], workSlice[index+1:]...)
	*iprgs = workSlice
}

// Add is an helper that adds a range made of a single IP
func (iprgs *IPRanges) Add(ip net.IP) {
	ipr := NewIPRange(ip, ip)
	*iprgs = append(*iprgs, ipr)
}

// AddString is an helper that parses and adds a range of IP from a string
func (iprgs *IPRanges) AddString(ipstr string) error {
	var ip net.IP

	if val := net.ParseIP(ipstr); val != nil {
		ip = val
	} else {
		return fmt.Errorf("invalid IP [%s]", ipstr)
	}

	iprgs.Add(ip.To4())

	return nil
}

// AddRange is an helper that adds a range of IP made of two IPs
func (iprgs *IPRanges) AddRange(lower net.IP, upper net.IP) {
	ipr := NewIPRange(lower, upper)
	*iprgs = append(*iprgs, ipr)
}

//
// IPRange methods
//

// ContainsIPString is an helper that checks if a range contains the given IP string
func (iprg IPRange) ContainsIPString(ipstr string) bool {
	var ip net.IP
	if val := net.ParseIP(ipstr); val != nil {
		ip = val
	} else {
		return false
	}

	return iprg.ContainsIP(ip)
}

// ContainsIP is an helper that checks if a range contains the given IP
func (iprg IPRange) ContainsIP(ip net.IP) bool {
	if bytes.Compare(ip.To4(), iprg.Lower) >= 0 && bytes.Compare(ip.To4(), iprg.Upper) <= 0 {
		return true
	}

	return false
}

// ContainsIPRange is an helper that checks if a range contains the given IP range
func (iprg IPRange) ContainsIPRange(iprange IPRange) bool {
	if iprg.ContainsIP(iprange.Lower.To4()) && iprange.ContainsIP(iprange.Upper.To4()) {
		return true
	}

	return false
}

// IsUpperOrLowerBoundary is an helper that checks if the given IP is either the lower of the upper bound of a range
func (iprg IPRange) IsUpperOrLowerBoundary(ip net.IP) bool {
	if !net.IP.Equal(ip.To4(), iprg.Lower) && !net.IP.Equal(ip.To4(), iprg.Upper) {
		return false
	}

	return true
}

// Equals is an helper that checks if an IPRange is equal to another
func (iprg *IPRange) Equals(iprange IPRange) bool {
	return net.IP.Equal(iprg.Upper, iprange.Upper.To4()) && net.IP.Equal(iprg.Lower, iprange.Lower.To4())
}

//
// CIDR
//

// WhitelistCIDR parses and adds a CIDR string to the IPRules' whitelist
func (iprl *IPRules) WhitelistCIDR(rawIPCIDR string) error {
	_, ipnet, err := net.ParseCIDR(rawIPCIDR)
	if err != nil {
		return err
	}

	ipFrom, ipTo, err := addressRange(ipnet)
	if err != nil {
		return err
	}

	iprl.WhitelistedIPs.AddRange(ipFrom, ipTo)

	return nil
}

// BlacklistCIDR parses and adds a CIDR string to the IPRules' blacklist
func (iprl *IPRules) BlacklistCIDR(rawIPCIDR string) error {
	_, ipnet, err := net.ParseCIDR(rawIPCIDR)
	if err != nil {
		return err
	}

	ipFrom, ipTo, err := addressRange(ipnet)
	if err != nil {
		return err
	}

	iprl.BlacklistedIPs.AddRange(ipFrom, ipTo)

	return nil
}

//
// Ranges
//

// WhitelistRange parses and adds an IP range string to the IPRules' whitelist
func (iprl *IPRules) WhitelistRange(rawIPRange string) error {
	var ipFrom net.IP
	var ipTo net.IP

	hostRange := strings.Split(rawIPRange, "-")

	lower, higher := hostRange[0], hostRange[1]

	if val := net.ParseIP(lower); val != nil {
		ipFrom = val
	} else {
		return fmt.Errorf("[%s-%s] is not a valid IP address range", lower, higher)
	}

	if val := net.ParseIP(higher); val != nil {
		ipTo = val
	} else {
		return fmt.Errorf("[%s-%s] is not a valid IP address range", lower, higher)
	}

	checkValidIPRange(ipFrom, ipTo)

	iprl.WhitelistedIPs.AddRange(ipFrom, ipTo)

	return nil
}

// BlacklistRange parses and adds an IP range string to the IPRules' blacklist
func (iprl *IPRules) BlacklistRange(rawIPRange string) error {
	var ipFrom net.IP
	var ipTo net.IP

	hostRange := strings.Split(rawIPRange, "-")

	lower, higher := hostRange[0], hostRange[1]

	if val := net.ParseIP(lower); val != nil {
		ipFrom = val
	} else {
		return fmt.Errorf("[%s-%s] is not a valid IP address range", lower, higher)
	}

	if val := net.ParseIP(higher); val != nil {
		ipTo = val
	} else {
		return fmt.Errorf("[%s-%s] is not a valid IP address range", lower, higher)
	}

	checkValidIPRange(ipFrom, ipTo)

	iprl.BlacklistedIPs.AddRange(ipFrom, ipTo)

	return nil
}

//
// Single IPs
//

// Whitelist checks the validity of an IP string and adds it to the IPRules' whitelist
func (iprl *IPRules) Whitelist(ip string) error {
	checkValidIP(ip)

	if err := iprl.WhitelistedIPs.AddString(ip); err != nil {
		return err
	}

	return nil
}

// Blacklist checks the validity of an IP string and adds it to the IPRules' blacklist
func (iprl *IPRules) Blacklist(ip string) error {
	checkValidIP(ip)

	if err := iprl.BlacklistedIPs.AddString(ip); err != nil {
		return err
	}

	return nil
}

//
// Checks
//

func checkValidIP(ipstr string) {
	if !isValidIPString(ipstr) {
		log.Println(fmt.Sprintf("[%s] is not a valid IP address", ipstr))
		os.Exit(1)
	}
}

func checkValidIPRange(lower net.IP, upper net.IP) {
	if !isValidIPRange(lower, upper) {
		log.Println(fmt.Sprintf("[%s-%s] is not a valid host range", lower.String(), upper.String()))
		os.Exit(1)
	}
}

func isValidIPString(ipstr string) bool {
	if val := net.ParseIP(ipstr); val == nil {
		return false
	}

	return true
}

func isValidIPRange(lower net.IP, upper net.IP) bool {
	return bytes.Compare(lower, upper) <= 0
}
