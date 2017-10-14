
## Relay Server
Application that helps you bypass nasty firewalls. In order to use the relay
server, please read this guide carefully as there are some considerations to
keep in mind when writing the code for your server.

## Architecture
Understanding the architecture of the server, will help us understand how
to write code for it.

As the relayServer cannot initiate any connections, it is setup to only receive
them. To this end, we spawn 2 sockets to listen for connections from the client
and the server each time a new RelayRequest is received.

When a client connects to the relayServer, we notify the server by sending a
string, delimited by the '\n' character, containing the host and port of the
relayServer that is handling this connection from the client. We then expect
the server to establish a connection to the relay server.

Once the server establishes this connection, we are left with 2 connections
that we then keep in synchrony. That is, everything that is written to the
clientConnection is also be also written to the serverConnection and vice versa.

## Clients
As mentioned earlier, servers that want to use our relayServer must first connect
to the relayServer and then listen for notifications of connections (a string sent
through a TCP stream). When the server learns there is a new client, it must then
create a new connection to the host received in the notification.
The final step is to create a connection to the server itself and synchronize this
connection with the one receiving the data from the relayServer. You can see two
examples of servers that use this method:

(echoServer.go)[https://gist.github.com/anonymous/a31dadee64238118229eb3ff13f1a340]
(webServer.go)[https://gist.github.com/anonymous/511f9943eab6d3bcdd75a3629e51486e]

