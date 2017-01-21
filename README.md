# sshego, a usable ssh library for Go

### executive summary

Google's "golang.org/x/crypto/ssh" library offers a
fantastic full implementation of the ssh
client and server protocols.
However this library is minimalistic by
design, cumbersome to figure out how to
use with RSA keys, and needs additional code to
support support tunneling
and receiving connections as an sshd.

`sshego` bridges this usability gap,
providing a drop-in Go library
to secure your tcp connections. In 
places `sshego` can be used in preference to
a virtual-private-network (VPN), for both
convenience and speed. Moreover the SSH
protocol's man-in-the-middle attack protection is
better than a VPN in almost all cases.

#### usable three-factor auth in an embeddable sshd

For strong security, our embedded sshd
offers three-factor auth (3FA). The three
security factors are: a passphrase ("what you know");
a 6-digit Google authenticator code (TOTP/RFC 6238;
"what you have": your phone); and
the use of PKI in the form of 4096-bit RSA keys.

To promote strong passphrases, we follow the
inspiration of https://xkcd.com/936/, and offer a user-
friendly 3-word starting prompt (the user
completes the sentence) to spark the
user's imagination in creating
a strong and memorizable passphrase. Passphrases of up
to 100 characters are supported.

Although not for the super-security conscious,
if desired and configured, passphrases can automatically
be backed up to email (via the Mailgun email service).

On new account creation with `gosshtun -adduser yourlogin`,
we will attempt to pop-up the QR-code on your
local desktop for quick Google Authenticator setup
on your phone.


### introducing sshego: a gopher's do-it-yourself ssh tunneling library

`sshego` is a golang (Go) library for ssh
tunneling (secure port forwarding). It also offers an
embeddable 3-factor authentication sshd server,
which can be useful for securing reverse forwards.

This means you can easily create an ssh-based vpn
with 3-factor authentication requrirements:
the embedded sshd requires passphrase, RSA keys,
and a TOTP Google Authenticator one-time password.

In addition to the libary, `gosshtun` is also
a command line utility (see the cmd/ subdir) that
demonstrates use of the library and may prove
useful on its own.

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

In its principal use, sshego is the equivalent
to using the ssh client and giving `-L` and/or `-R`.
It acts like an ssh client without a remote shell; it simply
tunnels other TCP connections securely.

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

# granuality of access

If you don't trust the other users on the host where your
process is running, you can also use sshego to (a) secure a direct
TCP connection (see DialConfig.Dial() and the example in cli_test.go; https://github.com/glycerine/sshego/blob/master/cli_test.go#L72);
or (b) forward via a file-system secured unix-domain sockets.
The first option (a) would disallow any other process (even
under the same user) from using
your connection, and the second (b) would disallow any other user from
accessing your tunnel, so long as you use the file-system permissions
to make the unix-domain socket path inaccessible to others.

# theory of operation

`gosshtun` and `sshego` will check the sshd server's host key. We prevent MITM attacks
by only allowing new servers if `-new` (a.k.a. SshegoConfig.AddIfNotKnown == true) is given.

When running the standalone `gosshtun` to test
a foward, you should give `-new` only once at setup time.

Then the lack of `-new` protects you on subsequent runs,
because the server's host key must match what we were
given the very first time.

# flags accepted, see `gosshtun -h` for complete list

~~~
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
        (default "$HOME/.ssh/.sshego.sshd.db")        
  -key string
        private key for sshd login (default "$HOME/.ssh/id_rsa_nopw")
  -known-hosts string
        path to gosshtun's own known-hosts file (default
        "$HOME/.ssh/.sshego.cli.known.hosts")
  -listen string
        (forward tunnel) We listen on this host:port locally,
        securely tunnel that traffic to sshd, then send it
        cleartext to -remote. The forward tunnel is active
        if and only if -listen is given.  If host starts with
        a '/' then we treat it as the path to a unix-domain
        socket to listen on, and the port can be omitted.
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
        only if -listen is given too.  If host starts with a
        '/' then we treat it as the path to a unix-domain
        socket to forward to, and the port can be omitted.
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
~~~

# installation

  go get github.com/glycerine/sshego/...

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


# specifying username to login to sshd host with

The `-user` flag should be used if your local $USER is different from that on the sshd host.

# source code for gosshtun command

See `github.com/glycerine/sshego/cmd/gosshtun/main.go` for the source code. This
also serves as an example of how to use the library.

# host key storage location (default)

`~/.ssh/.sshego.known.hosts.json.snappy`

# prep before running

a) install your passwordless ssh-private key in `~/.ssh/id_rsa_nopw` or use -key to say where it is.

b) add the corresponding public key to the user's .ssh/authorized_keys file on the sshd host.

# config file format

a) see demo.env for an example

b) run `gosshtun -write-config -` to generate a sample config file to stdout

c) comments are allowed; lines must start with `#`, comments continue until end-of-line

d) fields recognized (see `gosshtun -write-config -` for a full list)

~~~
#
# config file for sshego:
#
SSHD_ADDR="1.2.3.4:22"
FWD_LISTEN_ADDR="127.0.0.1:8888"
FWD_REMOTE_ADDR="127.0.0.1:22"
REV_LISTEN_ADDR=""
REV_REMOTE_ADDR=""
SSHD_LOGIN_USERNAME="$USER"
SSH_PRIVATE_KEY_PATH="$HOME/.ssh/id_rsa_nopw"
SSH_KNOWN_HOSTS_PATH="$HOME/.ssh/.sshego.known.hosts"
#
# optional in-process sshd
#
EMBEDDED_SSHD_HOST_DB_PATH="$HOME/.ssh/.sshego.sshd.db"
EMBEDDED_SSHD_LISTEN_ADDR="127.0.0.1:2022"
~~~

d) special environment reads

* The SSHD_LOGIN_USERNAME will subsitute $USER from the environment, if present.

* The *PATH keys will substitute $HOME from the environment, if present.

# MIT license

See the LICENSE file.

# Author

Jason E. Aten, Ph.D.

