/*
Package gosshtun is a golang libary that does secure port
forwarding over ssh.

Also `gosshtun` is a command line utility included here that
demonstrates use of the library; and may be useful standalone.

The intent of having a Go library is so that it can be used
to secure (via SSH tunnel) any other traffic that your
Go application would normally have to do over cleartext TCP.

While you could always run a tunnel as a separate process,
by running the tunnel in process with your application, you
know the tunnel is running when the process is running. It's
just simpler to administer; only one thing to start instead of two.

Also this is much simpler, and much faster, than using a
virtual private network (VPN). For a speed comparison,
consider [1] where SSH is seen to be at least 2x faster
than OpenVPN.

[1] http://serverfault.com/questions/653211/ssh-tunneling-is-faster-than-openvpn-could-it-be

In any case, you should realize that this is only an ssh
client, and not an sshd server daemon. It is the equivalent
to using the ssh client and giving `-L` and/or `-R`.

For example,

    gosshtun -listen 127.0.0.1:89  -sshd jumpy:55  -remote 10.0.1.5:80 -user alice -key ~/.ssh/id_rsa_nopw

is equivalent to

    ssh -N -L 89:10.0.1.5:80 alice@jumpy -port 55

with the addendum that `gosshtun` requires the use of passwordless
private `-key` file, and will never prompt you for a password at the keyboard.
This makes it ideal for embedding inside your application to
secure your (e.g. mysql, postgres, other cleartext) traffic. As
many connections as you need will be multiplexed over the
same ssh tunnel.

# theory of operation

We check the sshd server's host key. We prevent MITM attacks
by only allowing new servers if `-new` is given.

You should give `-new` only once at setup time.

Then the lack of `-new` can protect you on subsequent runs,
because the server's host key must match what we were
given the first time.

# options

 $ gosshtun -h
 Usage of gosshtun:
  -cfg string
        path to our config file
  -esshd string
        (optional) start an in-process embedded sshd (server),
        binding this host:port, with both RSA key and 2FA
        checking; useful for securing -revfwd connections.
  -esshd-host-db string
        (only matters if -esshd is also given) path
        to database holding sshd persistent state
        such as our host key, registered 2FA secrets, etc.
        (default "$HOME/.ssh/gosshtun.sshd.db")
  -key string
        private key for sshd login (default "$HOME/.ssh/id_rsa_nopw")
  -known-hosts string
        path to gosshtun's own known-hosts file (default
        "$HOME/.ssh/gosshtun.cli.known.hosts")
  -listen string
        (forward tunnel) We listen on this host:port locally,
        securely tunnel that traffic to sshd, then send it
        cleartext to -remote. The forward tunnel is active
        if and only if -listen is given.
  -new
        allow connecting to a new sshd host key, and store it
        for future reference. Otherwise prevent MITM attacks by
        rejecting unknown hosts.
  -quiet
        if -quiet is given, we don't log to stdout as each
        connection is made. The default is false; we log
        each tunneled connection.
  -remote string
        (forward tunnel) After traversing the secured forward
        tunnel, -listen traffic flows in cleartext from the
        sshd to this host:port. The foward tunnel is active
        only if -listen is given too.
  -revfwd string
        (reverse tunnel) The gosshtun application will receive
        securely tunneled connections from -revlisten on the
        sshd side, and cleartext forward them to this host:port.
        For security, it is recommended that this be 127.0.0.1:22,
        so that the sshd service on your gosshtun host
        authenticates all remotely initiated traffic.
        See also the -esshd option which can be used to
        secure the -revfwd connection as well.
        The reverse tunnel is active only if -revlisten is given
        too. (default "127.0.0.1:22")
  -revlisten string
        (reverse tunnel) The sshd will listen on this host:port,
        securely tunnel those connections to the gosshtun application,
        whence they will cleartext connect to the -revfwd address.
        The reverse tunnel is active if and only if -revlisten is given.
  -sshd string
        The remote sshd host:port that we establish a secure tunnel to;
        our public key must have been already deployed there.
  -user string
        username for sshd login (default is $USER)
  -v    verbose debug mode
  -write-config string
        (optional) write our config to this path before doing
        connections
 $

# example use of the command

    $ gosshtun -listen localhost:8888 -sshd 10.0.1.68:22 -remote 127.0.0.1:80

means the following two network hops will happen, when a local browser connects to localhost:8888

                           `gosshtun`             `sshd`
    local browser ----> localhost:8888 --(a)--> 10.0.1.68:22 --(b)--> 127.0.0.1:80
      `host A`             `host A`               `host B`              `host B`

where (a) takes place inside the previously established ssh tunnel.

Connection (b) takes place over basic, un-adorned, un-encrypted TCP/IP. Of
course you could always run `gosshtun` again on the remote host to
secure the additional hop as well, but typically -remote is aimed at the 127.0.0.1,
which will be internal to the remote host itself and so needs no encryption.

*/
package gosshtun
