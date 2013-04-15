#!/usr/bin/perl -w
# $Id: check_ip_conntrack.pl
use strict;
use Getopt::Std;

use vars qw($opt_c $opt_w
            $count_ip_conntrack $max_ip_conntrack
            $crit_level $warn_level
            %exit_codes @memlist
            $percent $fmt_pct 
            $verb_err $command_line);

# Predefined exit codes for Nagios
%exit_codes   = ('UNKNOWN' ,-1,
                 'OK'      , 0,
                 'WARNING' , 1,
                 'CRITICAL', 2,);

# Turn this to 1 to see reason for parameter errors (if any)
$verb_err     = 0;

my @conntrack_count_files = qw(
        /proc/sys/net/netfilter/nf_conntrack_count
        /proc/sys/net/ipv4/netfilter/ip_conntrack_count
);

my $count;
foreach (@conntrack_count_files) {
    if ( -r $_) {
        chomp($count = `cat $_`);
        last;
    }
}


my @conntrack_max_files = qw(
        /proc/sys/net/nf_conntrack_max
        /proc/sys/net/netfilter/nf_conntrack_max
        /proc/sys/net/ipv4/ip_conntrack_max
        /proc/sys/net/ipv4/netfilter/ip_conntrack_max
);

my $max;
foreach (@conntrack_max_files) {
    if ( -r $_) {
        chomp($max = `cat $_`);
        last;
    }
}

if (!$max) {
    print "ip_conntrack_max is NG\n";
    exit;
}

# Define the calculating scalars
$count_ip_conntrack  = $count;
$max_ip_conntrack = $max;

# Get the options
if ($#ARGV le 0)
{
  &usage;
}
else
{
  getopts('c:w:');
}

# Shortcircuit the switches
if (!$opt_w or $opt_w == 0 or !$opt_c or $opt_c == 0)
{
  print "*** You must define WARN and CRITICAL levels!" if ($verb_err);
  &usage;
}


$warn_level   = $opt_w;
$crit_level   = $opt_c;

$percent    = $count_ip_conntrack / $max_ip_conntrack * 100;
$fmt_pct    = sprintf "%.1f", $percent;
if ($percent >= $crit_level)
{
print "ip_conntrack CRITICAL - $fmt_pct% ($count_ip_conntrack) used\n";
exit $exit_codes{'CRITICAL'};
}
elsif ($percent >= $warn_level)
{
print "ip_conntrack WARNING - $fmt_pct% ($count_ip_conntrack) used\n";
exit $exit_codes{'WARNING'};
}
else
{
print "ip_conntrack OK - table usage = $fmt_pct%, count = $count_ip_conntrack\n";
exit $exit_codes{'OK'};
}

# Show usage
sub usage()
{
  print "\ncheck_ip_conntrack.pl v1.0 - Nagios Plugin\n\n";
  print "usage:\n";
  print " check_ip_conntrack.pl -w <warnlevel> -c <critlevel>\n\n";
  print "options:\n";
  print " -w PERCENT   Percent used when to warn\n";
  print " -c PERCENT   Percent used when critical\n";
  exit $exit_codes{'UNKNOWN'}; 
}