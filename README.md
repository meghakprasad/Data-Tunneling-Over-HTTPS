# Data Tunneling Over HTTPS

Data tunneling over HTTPS refers to the process of encapsulating one type of data within the HTTPS protocol for secure transmission over a network. HTTPS, or Hypertext Transfer Protocol Secure, is a protocol that ensures secure communication between a client (such as a web browser) and a server over the internet. It uses encryption to protect data integrity and confidentiality.

Data tunneling involves embedding data from one protocol or application within the HTTPS protocol to leverage its security features. This approach is commonly used when the underlying network or infrastructure does not directly support the desired protocol, or when additional security measures are required.

For example, let's say you have a non-secure protocol (Protocol A) that you want to transmit over a network that only allows HTTPS traffic. You can establish an HTTPS connection between the client and server, and then encapsulate the Protocol A data within the HTTPS packets. This way, the data is protected by the security mechanisms provided by HTTPS, such as encryption and server authentication.

Data tunneling over HTTPS is often used in scenarios like:

1. VPN (Virtual Private Network): VPN services often use HTTPS as a means of tunneling data securely over the internet, ensuring privacy and encryption for the transmitted data.
2. Secure proxying: In some cases, data is tunneled over HTTPS to bypass firewalls or network restrictions by making the traffic appear as regular HTTPS traffic.
3. Legacy application integration: If you have a legacy application that communicates using an insecure protocol, you can tunnel the communication through HTTPS to add an extra layer of security.

It's important to note that data tunneling over HTTPS does not modify the HTTPS protocol itself but leverages its existing infrastructure for secure data transmission.

# To Run The Project

Three Ubuntu Machines are required to run the project.
Set network settings to use Bridged Adapter with Promiscuous mode set to Allow All on all 3 machines.

consider,
IP Address of Ubuntu1:192.168.29.224
IP Address of Ubuntu2:192.168.29.200
IP Address of Ubuntu3:192.168.29.137

1) In the first ubuntu (i.e., Ubuntu1) install Golang and RabbitMq.Create a user with the username as admin1 and password as admin123 in the RabbitMq and grant admin privileges.

2) Add app1 to ubuntu1.

3) In the second ubuntu (i.e., Ubuntu2) install Golang and add app2 to Ubuntu2.

4) In the third ubuntu (i.e., Ubuntu3) install mongodb. Create a user in the admin database with username: mongoAdmin and password: admin123 and set the role to userAdminAnyDatabase

5) Install Ghostunnel on all 3 machines.

6) Create certs on Ubuntu1 and Ubuntu3 and copy the certs to ubuntu2.

7) Start Ghostunnel between ubuntu1 and ubuntu2 considering ubuntu1 as the server and ubuntu2 as the client.

Example: 

(IN UBUNTU1)

	nc -l localhost 8080

	sudo ghostunnel server     --listen 0.0.0.0:443     --target localhost:8080    --key certi/server.key --cert certi/server.crt     --cacert certi/ca.pem     --allow-all

(IN UBUNTU2)

	ghostunnel client     --listen localhost:8080     --target 192.168.29.224:443     --key certi/server.key     --cert certi/server.crt     --cacert certi/ca.pem

	nc -v localhost 8080


8) Start Ghostunnel between ubuntu3 and ubuntu2 on port8070 considering ubuntu3 as the server and ubuntu2 as the client.

Example:

(IN UBUNTU3)

	nc -l localhost 8070

	sudo ghostunnel server     --listen 0.0.0.0:443     --target localhost:8070    --key cert/server.key --cert cert/server.crt     --cacert cert/ca.pem     --allow-all


(IN UBUNTU2)

	ghostunnel client     --listen localhost:8070     --target 192.168.29.137:443     --key cert/server.key     --cert cert/server.crt     --cacert cert/ca.pem

	nc -v localhost 8070


9) once the connection is estblished and the machines are able to communicate Start app1 in ubuntu1 and start app2 in ubuntu2.

10) Use curl command to send a message from app1 to app2.

Example: 

	curl --header "Content-Type: application/json"      --request POST      --data '{"name":"Jane","email":"jane@gmail.com"}'       http://localhost:8090/people

11) The message recieved will be inserted to the database in Ubuntu3.
