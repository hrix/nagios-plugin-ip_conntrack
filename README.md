#Nagios ip_conntrack check plugin
<pre><code>
check_ip_conntrack.pl v1.0 - Nagios Plugin

usage:
 check_ip_conntrack.pl -w <warnlevel> -c <critlevel>

options:
 -w PERCENT   Percent used when to warn
 -c PERCENT   Percent used when critical
</code></pre>
* * *
##INSTALLATION
###1. get the file
<pre><code>
 cd /usr/lib64/nagios/plugins/
 wget https://raw.github.com/S1100/nagios-plugin-ip_conntrack/master/check_ip_conntrack.pl
 chmod 755 check_ip_conntrack_max.pl
</code></pre>
###2. check the response
<pre><code>
$ time /usr/lib64/nagios/plugins/check_ip_conntrack.pl -w 80 -c 90
ip_conntrack OK - table usage = 0.1%, count = 75

real    0m0.075s
user    0m0.024s
sys     0m0.030s
</code></pre>
If your server needs over 5 second, should not to use this plugin.

###3. allow nrpe user to count ip_conntrack
<pre><code>
$ ps aux|grep nrpe|grep -v grep
nagios   27083  0.0  0.0   5096   724 ?        Ss   Dec20   0:05 /usr/sbin/nrpe -c /etc/nagios/nrpe.cfg -d

# visudo
--- edit permission
# for check_ip_conntrack
 %nagios       ALL=NOPASSWD:/usr/bin/wc

# comment out this line
#Defaults    requiretty
</code></pre>
###4. add check_ip_conntrack command on your nrpe.cfg
<pre><code>
# command[check_ip_conntrack]=/usr/lib64/nagios/plugins/check_ip_conntrack.pl -w 80 -c 90
# service nrpe restart
</code></pre>
###5. check from server
* FROM Nagios server
<pre><code>
/usr/lib64/nagios/plugins/check_nrpe -H [node IP address] -c check_ip_conntrack
</code></pre>
If NG, you can check the node's /var/log/secure and do over again node's visudo.
###6. add service on your nagios config
<pre><code>
define service{
  use                     huge
  host_name               huge_deliver1
  service_description     check_ip_conntrack
  check_command           check_nrpe!check_ip_conntrack
}
</code></pre>