package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Predefined exit codes for Nagios
var exit_codes = map[string]int{
	"UNKNOWN":  -1,
	"OK":       0,
	"WARNING":  1,
	"CRITICAL": 2,
}

// Show usage
func usage() {
	fmt.Println("\ncheck_ip_conntrack.pl v1.0 - Nagios Plugin\n")
	fmt.Println("usage:")
	fmt.Println(" check_ip_conntrack.pl -w <warnlevel> -c <critlevel>\n")
	fmt.Println("options:")
	fmt.Println(" -w PERCENT   Percent used when to warn")
	fmt.Println(" -c PERCENT   Percent used when critical")
	os.Exit(exit_codes["UNKNOWN"])
}

func isReadable(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		//fmt.Printf("no such file or directory: %s", filename)
		return false
	}

	// check readable
	_, oerr := os.Open(filename)
	if oerr != nil {
		//fmt.Printf("cant open such file or directory: %s", filename)
		return false
	}
	return true
}

func get_count_value(filename string) (int, error) {
	fbyte, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	//count, nerr := strconv.ParseUint(string(fbyte), 10, 0)
	count, nerr := strconv.Atoi(strings.TrimRight(string(fbyte), "\n\r"))
	if nerr != nil {
		return 0, nerr
	}
	return count, nil
}

func collect_values() (count_ip_conntrack, max_ip_conntrack int) {
	var err error
	conntrack_count_files := []string{
		"/proc/sys/net/netfilter/nf_conntrack_count",
		"/proc/sys/net/ipv4/netfilter/ip_conntrack_count",
	}
	for _, f := range conntrack_count_files {
		if isReadable(f) {
			// fmt.Println("Exists && readable") // debug
			count_ip_conntrack, err = get_count_value(f)
			if err != nil {
				continue
			}
			break
		}
	}

	conntrack_max_files := []string{
		"/proc/sys/net/nf_conntrack_max",
		"/proc/sys/net/netfilter/nf_conntrack_max",
		"/proc/sys/net/ipv4/ip_conntrack_max",
		"/proc/sys/net/ipv4/netfilter/ip_conntrack_max",
	}
	for _, mf := range conntrack_max_files {
		if isReadable(mf) {
			max_ip_conntrack, err = get_count_value(mf)
			if err != nil {
				fmt.Println(err.Error()) // debug
				continue
			}
			break
		}
		fmt.Println(mf, " is not exists")
	}

	if max_ip_conntrack <= 0 {
		fmt.Println("ip_conntrack_max is NG")
		os.Exit(exit_codes["UNKNOWN"])
	}
	// Define the calculating scalars
	return
}

func check_limit(count_ip_conntrack, max_ip_conntrack, warn_level, crit_level int) {
	var percent float64
	percent = float64(count_ip_conntrack) / float64(max_ip_conntrack) * 100.0
	fmt_pct := fmt.Sprintf("%.1f", percent)
	//fmt.Println(fmt_pct) // debug
	switch {
	case percent >= float64(crit_level):
		fmt.Printf("ip_conntrack CRITICAL - %s%% (%d) used\n", fmt_pct, count_ip_conntrack)
		os.Exit(exit_codes["CRITICAL"])
	case percent >= float64(warn_level):
		fmt.Printf("ip_conntrack WARNING - %s%% (%d) used\n", fmt_pct, count_ip_conntrack)
		os.Exit(exit_codes["WARNING"])
	default:
		fmt.Printf("ip_conntrack OK - table usage = %s%%, count = %d \n", fmt_pct, count_ip_conntrack)
		os.Exit(exit_codes["OK"])
	}
}

func main() {
	opt_warn := flag.Int("w", 0, "warning")
	opt_crit := flag.Int("c", 0, "critical")
	flag.Parse()
	warn_level := *opt_warn
	crit_level := *opt_crit
	if warn_level == 0 || crit_level == 0 {
		usage()
	}

	count_ip_conntrack, max_ip_conntrack := collect_values()

	check_limit(count_ip_conntrack, max_ip_conntrack, warn_level, crit_level) // debug
}
