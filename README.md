# Data-Tunneling-through-HTTPS

Three Ubuntu Machines are required to run the project.Set network settings to use Bridged Adapter with Promiscuous mode set to Allow All on all 3 machines.

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
