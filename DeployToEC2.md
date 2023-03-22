# start by crate repo and passe it to it (PUPlic then converter to private)

Go to >> https://us-west-2.console.aws.amazon.com/

EC2 >> Launch an instance 

&& Reourses [
    https://medium.com/hackernoon/deploying-a-go-application-on-aws-ec2-76390c09c2c5
    https://aws.plainenglish.io/how-to-install-mongodb-in-aws-ec2-99958976abda
    https://www.youtube.com/watch?v=7vf210p2tJg
    https://github.com/aaronwht/EC2    
]
# 80 8080 5001 8090
 ssh -i "ec2acesskey.pem" ec2-user@ec2-44-214-1-200.compute-1.amazonaws.com

create and download key pair (login) | named GoServerAccess 
|| Get the File on the App folder

|| Edit network settings and add security group rule
|| type custom TCP With Port range AppPort || 5000
|| source type > Anywhere

[launching-instance]

>> Go EC2>>instances

Select our instance (instance name)

Click Connect
Example:
>> ssh -i "GoServerAccess.pem" ec2-user@ec2-34-209-204-151.us-west-2.compute.amazonaws.com

[copy-our-instace-link] || Public DNS (IPv4)
>> ec2-34-209-204-151.us-west-2.compute.amazonaws.com

Connect to our server
>> ssh -i "GoServerAccess.pem" ec2-user@ec2-34-209-204-151.us-west-2.compute.amazonaws.com

!!AFTER CONNECTING!!

>> sudo yum update -y
>> sudo yum install -y golang
>> export GOROOT=/usr/lib/golang
>> export PATH=$PATH:$GOROOT/bin
>> go version
 


>> ls
>> pwd
[/home/ec2-user]

# clone the repo
 git  clone https://github.com/ahmedkhalaf1996/GoAws.git
 git  clone https://github.com/Danc-app/GolangServer.git
>> git clone https://github.com/Danc-app/GolangServer.git
>> git clone 
>> clear
>> ls
GolangServer

>> cd GolangServer
>>  go build main.go

>>   ./main
Go to puplic dns ipv4 and add the port number 5000
http://ec2-34-209-204-151.us-west-2.compute.amazonaws.com:5000/

>> crtl+C
>> Start Setup mongodb
>> cd ..
>> cd ..
>> cd ..

>> sudo vim /etc/yum.repos.d/mongodb-org-6.0.repo

#### hit I to enter INSERT mode, copy and paste the following

[mongodb-org-6.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/amazon/2/mongodb-org/6.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-6.0.asc

#### hit esc key to exit INSERT mode, and run :wq to save the file.
####then hit enter

>> sudo yum install -y mongodb-org

>> mongod --version
>> sudo systemctl start mongod
>> sudo systemctl status mongod


# run server using pm2
Install Node Version Manager:
>> curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash

>> exit
[terminateAccessSSH-thenGoBack]
>>  ssh -i "GoServerAccess.pem" ec2-user@ec2-34-209-204-151.us-west-2.compute.amazonaws.com

>>  nvm -v
>>  nvm install 16

>> exit
[terminateAccessSSH-thenGoBack]
>>  ssh -i "GoServerAccess.pem" ec2-user@ec2-34-209-204-151.us-west-2.compute.amazonaws.com

>> node -v

>> npm install pm2 -g
>> pm2 -version

>> cd GoAws
>>  pm2 start ./main

>> pm2 stop ./main

>> pm2 restart ./main
>> pm2 reload ./main


>> mongosh

>> use MY_DATABASE
OR
>> use test

>>  db.dropDatabase() 
## show databases names to list existing databases
>> show dbs

<!-- 
db.product.insert({
  name: "Apple",
  description: "This is an Apple",
  price: 2.50
});
... -->

<!-- db.product.find({}) -->

db.users.find({})


# remove 
rm -rf


## Update Git Repo

>> git pull origin


----------------------
https://www.youtube.com/watch?v=q-XEGbipOVw

https://app.zerossl.com/certificate/install/aeaf62a8a3180fa2608f8bcda8cd0154

create enctance with allow http and https also add our server ports

after connecting

>> sudo yum update -y 
>> sudo amazon-linux-extras install -y nginx1 
>> sudo service nginx start 
>> sudo yum install -y golang 
>> export PATH=$PATH:$GOROOT/bin 

>> go version

>> git clone https://github.com/Danc-app/GolangServer.git

>> sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-ports 8080 
>> sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-ports 5000

>> sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-ports 8080 
>> sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-ports 5000 

>> sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-ports 5001 
>> sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-ports 5001 

