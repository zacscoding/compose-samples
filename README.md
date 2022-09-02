# :books: Docker Compose Samples

This project is intended for use in local development environment i.e. do not use in production environments :thinking:

# Table Of Contents

- [Simple Commands](#Simple-Commands)
- [Postgres](#Postgres)
- [MySQL](#MySQL)
- [MySQL Cluster](#MySQL-Cluster)
- [DynamoDB](#DynamoDB)
- [Ubuntu](#Ubuntu)
- [Alpine Java](#Alpine-Java)
- [OpenAPI](#OpenAPI)
- [InfluxDB](#InfluxDB)
- [Grafana](#Grafana)
- [Zookeeper](#Zookeeper)
- [Zookeeper Cluster](#Zookeeper-Cluster)
- [Kafka](#Kafka)
- [Kafka Cluster](#Kafka-Cluster)
- [RabbitMQ](#RabbitMQ)
- [Redis](#Redis)

---  

# Simple Commands

```shell
$ docker-compose -f ./docker-compose.yaml up -d
$ docker logs -f [container_name]
$ docker-compose -f ./docker-compose.yaml down -v
```

---  

# MariaDB

See [./composes/mariadb](./composes/mariadb) for details

> ### docker-compose.yml

```yaml
version: '3.4'
services:
  mariadb:
    image: mariadb:10.2
    container_name: mariadb
    environment:
      MYSQL_ROOT_PASSWORD: pass
      MYSQL_DATABASE: testdb
      MYSQL_USER: tester
      MYSQL_PASSWORD: tester
    ports:
      - "3306:3306"
    logging:
      driver: syslog
      options:
        tag: "{{.DaemonName}}(image={{.ImageName}};name={{.Name}};id={{.ID}})"
    networks:
      - backend
    restart: on-failure
    volumes:
      - ${PWD}/mariadb:/var/lib/mysql
      - ${PWD}/custom.cnf:/etc/mysql/conf.d/custom.cnf
networks:
  backend:
    driver: bridge
```

> ### Start containers and Checks

```shell
$ cd composes/mariadb
$ docker-compose up -d

$ docker exec -it mariadb bash
root@19d6c77a0f0e:/# mysql -u tester -p testdb
Enter password:
Welcome to the MariaDB monitor.  Commands end with ; or \g.
Your MariaDB connection id is 9
Server version: 10.2.23-MariaDB-1:10.2.23+maria~bionic mariadb.org binary distribution

Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

MariaDB [testdb]> show tables;
Empty set (0.01 sec)

MariaDB [testdb]> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| testdb             |
+--------------------+
2 rows in set (0.00 sec)
```

---  

See [./composes/postgres](./composes/postgres) for details

> ### docker-compose.yml

```yaml
version: '3.1'

services:
  postgresdb:
    image: postgres
    container_name: postgresdb
    restart: always
    environment:
      - POSTGRES_PASSWORD=pass
    ports:
      - "5432:5432"
```  

> ### Start containers and Checks

```shell
$ docker exec -it postgresdb bash
root@f979f2462b94:/# psql -d postgres -U postgres
psql (11.2 (Debian 11.2-1.pgdg90+1))
Type "help" for help.

postgres=#
```

---  

# MySQL

See [./composes/mysql](./composes/mysql) for details

> ### docker-compose.yml

```yaml
version: '3.1'
services:
  mysqldb:
    image: mysql:8.0.17
    container_name: mysqldb
    platform: linux/amd64 # for m1
    command: [ '--default-authentication-plugin=mysql_native_password', '--default-storage-engine=innodb' ]
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=my_database
    ports:
      - 3306:3306
```

> ### Start containers and Checks

```shell
$ docker exec -it mysqldb bash
bash-4.4# mysql -u root -p my_database
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.0.30 MySQL Community Server - GPL

Copyright (c) 2000, 2022, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+--------------------+
| Database           |
+--------------------+
...
```

---  

# MySQL Cluster

See [./composes/mysqlcluster](./composes/mysqlcluster) for details

> ### docker-compose.yml

```yaml
version: '3'
services:
  mysql_master:
    platform: linux/amd64
    image: mysql:5.7
    env_file:
      - ./master/mysql_master.env
    container_name: "mysql_master"
    restart: "no"
    ports:
      - 4406:3306
    volumes:
      - ./master/conf/mysql.conf.cnf:/etc/mysql/conf.d/mysql.conf.cnf
      - ./master/__data:/var/lib/mysql
    networks:
      - overlay

  mysql_slave:
    platform: linux/amd64
    image: mysql:5.7
    env_file:
      - ./slave/mysql_slave.env
    container_name: "mysql_slave"
    restart: "no"
    ports:
      - 5506:3306
    depends_on:
      - mysql_master
    volumes:
      - ./slave/conf/mysql.conf.cnf:/etc/mysql/conf.d/mysql.conf.cnf
      - ./slave/__data:/var/lib/mysql
    networks:
      - overlay

networks:
  overlay:
```

> ### Start containers and Checks

```shell
# Run mysql cluster
$ cd composes/mysqlcluster
$ ./composes/mysqlcluster/build.sh
# Run tests. this example try to write(insert/update/delete) to master and read(select) from slave.
$ go run main.go
2022/09/03 03:10:27 //==============================================
2022/09/03 03:10:27 Try to save an user
[Callback - Create] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
2022/09/03 03:10:27 ================================================//
2022/09/03 03:10:27 //==============================================
2022/09/03 03:10:27 Try to update an user
[Callback - Update] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
2022/09/03 03:10:27 ================================================//
2022/09/03 03:10:28 //==============================================
2022/09/03 03:10:28 Try to find an user by calling First()
[Callback - Query] >> current user: Slave(mydb_slave_user@192.168.96.1), err: <nil>
2022/09/03 03:10:28 ================================================//
2022/09/03 03:10:28 //==============================================
2022/09/03 03:10:28 Try to find an user by calling exec()
[Callback - Row] >> current user: Slave(mydb_slave_user@192.168.96.1), err: <nil>
2022/09/03 03:10:28 ================================================//
2022/09/03 03:10:29 //==============================================
2022/09/03 03:10:29 Try to find an user with manual switching
[Callback - Query] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
2022/09/03 03:10:29 ================================================//
2022/09/03 03:10:29 //==============================================
2022/09/03 03:10:29 Try to delete an user
[Callback - Delete] >> current user: Master(mydb_user@192.168.96.1), err: <nil>
2022/09/03 03:10:29 ================================================//
```

---

# DynamoDB

See [./composes/dynamodb](./composes/dynamodb) for details

> ### docker-compose.yml

```yaml
version: '3.1'
services:
  dynamodb:
    image: amazon/dynamodb-local:latest
    container_name: dynamodb
    ports:
      - "8000:8000"
    volumes:
      - $HOME/.aws:/root/.aws
  dynamodb-ui:
    restart: always
    image: aaronshaf/dynamodb-admin
    container_name: dynamodb-ui
    environment:
      - DYNAMO_ENDPOINT=http://dynamodb:8000
    ports:
      - 8001:8001
    volumes:
      - $HOME/.aws:/root/.aws
  dynamodb-init:
    image: amazon/aws-cli
    entrypoint: /bin/sh -c
    container_name: dynamodb-init
    command: "/dynamodb/init.sh"
    environment:
      - ENVIRONMENT=LOCAL
      - HOST=dynamodb:8000
    depends_on:
      - dynamodb
    volumes:
      - ${HOME}/.aws:/root/.aws
      - ./init.sh:/dynamodb/init.sh
      - ./tables:/dynamodb/tables
```

> ### Start containers and Checks

```shell
$ cd composes/dynamodb
$ docker-compose up -d
```

Connect to [http://localhost:8001](http://localhost:8001) in your browser.

---  

# Ubuntu

See [./composes/ubuntu](./composes/ubuntu) for details

> ### Dockerfile

```
FROM ubuntu:16.04

RUN apt-get update && apt-get install -y sudo && apt-get install -y openssh-server
RUN apt-get install net-tools
RUN mkdir /var/run/sshd
RUN echo 'root:rootpw' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

# Install docker
RUN apt-get install --assume-yes  apt-transport-https  ca-certificates  curl  gnupg-agent  software-properties-common && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && apt-get update && apt-get install --assume-yes docker-ce docker-ce-cli containerd.io

# Install docker-compose (TODO : dependencies)
# RUN curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
```

> ### Start containers and Checks

```shell
$ cd composes/ubuntu
$ docker-compose up -d
# $ docker build -t eg_sshd .
$ ssh root@localhost -p 49154
password: rootpw
```

---  

# Alpine Java

See [./composes/alpine/java](./composes/alpine/java) for details

> ### docker-compose.yml

```yaml
version: '3.1'

services:
  alpine-java:
    platform: linux/x86_64
    image: openjdk:8-jre-alpine
    container_name: alpine-java
    command:
      sh -c 'echo before comamnd; cd /tmp; java HelloWorld scan01; echo after command'
    volumes:
      - ${PWD}/HelloWorld.class:/tmp/HelloWorld.class # this class compiled from 1.8
```

> ### Start containers and Checks

```shell
$ cd composes/alpine/java
$ docker-compose up
Starting alpine-java ... done
Attaching to alpine-java
alpine-java    | before comamnd
alpine-java    | Hello World~!
alpine-java    | after command
alpine-java exited with code 0
```

---  

# OpenAPI

See [./composes/openapi](./composes/openapi) for details

> ### docker-compose.yml

```yaml
version: '3.4'
services:
  swagger-ui:
    image: swaggerapi/swagger-ui
    container_name: swagger-ui
    environment:
      - SWAGGER_JSON=/config/sample-api.yaml
      - BASE_URL=/swagger
    ports:
      - "8080:8080"
    volumes:
      - ${PWD}/sample-api.yaml:/config/sample-api.yaml
  httpd:
    image: httpd:latest
    container_name: httpd
    ports:
      - "8081:80"
    volumes:
      - ./sample-api.html:/usr/local/apache2/htdocs/docs.html
```

> ### build open api yaml to html by using redoc-cli

```shell
$ cd composes/openapi 
$ npm i -g redoc-cli
$ redoc-cli build ./sample-api.yaml -o ./sample-api.html
Prerendering docs

🎉 bundled successfully in: ./sample-api.html (1086 KiB) [⏱ 0.082s]
```

> ### Start containers and Checks

```shell
$ cd composes/openapi
$ docker-compose up
```

- swagger-ui: [http://localhost:8080/swagger](http://localhost:8080/swagger)
- generated openapi docs: [http://localhost:8081/docs.html](http://localhost:8081/docs.html)

---  

# InfluxDB

See [./composes/influxdb](./composes/influxdb) for details

> ### docker-compose.yml

```yaml
version: '3.4'
services:
  influxdb:
    image: influxdb:latest
    container_name: influxdb
    ports:
      - "8083:8083"
      - "8086:8086"
      - "8090:8090"
    environment:
      - INFLUXDB_DB=db0
      - INFLUXDB_ADMIN_USER=admins
      - INFLUXDB_ADMIN_PASSWORD=pass
    volumes:
      # Mount data directory for persistent.
      - ./influxdb/__data:/var/lib/influxdb
  chronograf:
    image: chronograf:latest
    ports:
      - '8888:8888'
    volumes:
      - ./__chronograf-storage:/var/lib/chronograf
    depends_on:
      - influxdb
    environment:
      - INFLUXDB_URL=http://influxdb:8086
      - INFLUXDB_USERNAME=admin
      - INFLUXDB_PASSWORD=pass
```  

> ### Start containers and Checks

```shell
$ cd composes/influxdb
$ docker-compose up -d
```

Connect to [http://localhost:8888](http://localhost:8888) in your browser.

---  

# Grafana

See [./composes/grafana](./composes/grafana) for details

> ### docker-compose.yml

```yaml
version: '3'
services:
  grafana:
    container_name: grafana
    image: grafana/grafana:latest
    user: "1000"
    ports:
      - 3000:3000
    volumes:
      - ./__gfdata:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_USER=admins
      - GF_SECURITY_ADMIN_PASSWORD=pass
```  

> ### Start containers and Checks

```shell
$ cd composes/grafana
$ docker-compose up -d
```

Connect to [http://localhost:3000](http://localhost:3000) in your browser.

---  

# Zookeeper

See [./composes/zookeeper](./composes/zookeeper) for details

> ### docker-compose.yml

```shell
version: '3.1'

services:
  zoo1:
    image: zookeeper:latest
    restart: always
    hostname: zoo1
    ports:
      - "2181:2181"
    environment:
      ZOO_MY_ID: 1
      ZOO_SERVERS: server.1=0.0.0.0:2888:3888;2181
    volumes:
      - ./__zookeeper1/data:/data
      - ./__zookeeper1/datalog:/datalog
```

> ### Start containers and Checks

```shell
$ cd composes/zookeeper
$ docker-compose up

# This example try to create and delete znode with watcher.
$ go run main.go
2022/09/03 03:21:50 [EventLoop] EventOccur: {EventSession StateConnecting  <nil> [::1]:2181}
2022/09/03 03:21:50 [ZKClient]connected to [::1]:2181
2022/09/03 03:21:50 [EventLoop] EventOccur: {EventSession StateConnected  <nil> [::1]:2181}
2022/09/03 03:21:50 [ZKClient]authenticated: id=72058349199360000, timeout=4000
2022/09/03 03:21:50 [EventLoop] EventOccur: {EventSession StateHasSession  <nil> [::1]:2181}
2022/09/03 03:21:50 [ZKClient]re-submitting `0` credentials after reconnect
2022/09/03 03:21:50 '/MyFirstZnode' exists: false, stat: &{0 0 0 0 0 0 0 0 0 0 0}
2022/09/03 03:21:50 '/MyFirstZnode' create result: /MyFirstZnode
2022/09/03 03:21:50 Stopping event loop
2022/09/03 03:21:50 [ZKClient]recv loop terminated: EOF
2022/09/03 03:21:50 [ZKClient]send loop terminated: <nil>
```

---

# Zookeeper Cluster

See [./composes/zookeepercluster](./composes/zookeepercluster) for details

> ### docker-compose.yml

```yaml
version: '3.1'

services:
  zookeeper1:
    image: 'bitnami/zookeeper:latest'
    container_name: zookeeper1
    ports:
      - '12181:2181'
      - '12888:2888'
      - '13888:3888'
    volumes:
      - __zookeeper1_data:/bitnami
    environment:
      - ZOO_SERVER_ID=1
      - ALLOW_ANONYMOUS_LOGIN=yes
      - ZOO_SERVERS=0.0.0.0:2888:3888,zookeeper2:2888:3888,zookeeper3:2888:3888
  zookeeper2:
    image: 'bitnami/zookeeper:latest'
    container_name: zookeeper2
    ports:
      - '22181:2181'
      - '22888:2888'
      - '23888:3888'
    volumes:
      - __zookeeper2_data:/bitnami
    environment:
      - ZOO_SERVER_ID=2
      - ALLOW_ANONYMOUS_LOGIN=yes
      - ZOO_SERVERS=zookeeper1:2888:3888,0.0.0.0:2888:3888,zookeeper3:2888:3888
  zookeeper3:
    image: 'bitnami/zookeeper:latest'
    container_name: zookeeper3
    ports:
      - '32181:2181'
      - '32888:2888'
      - '33888:3888'
    volumes:
      - __zookeeper3_data:/bitnami
    environment:
      - ZOO_SERVER_ID=3
      - ALLOW_ANONYMOUS_LOGIN=yes
      - ZOO_SERVERS=zookeeper1:2888:3888,zookeeper2:2888:3888,0.0.0.0:2888:3888

volumes:
  __zookeeper1_data:
    driver: local
  __zookeeper2_data:
    driver: local
  __zookeeper3_data:
    driver: local
```

> ### Start containers and Checks

```shell
$ cd composes/zookeepercluster
$ docker-compose up

# This example try to create and delete znode with watcher.
$ go run main.go
2022/09/03 03:24:25 Conn State: StateDisconnected
2022/09/03 03:24:25 [EventLoop] EventOccur: {EventSession StateConnecting  <nil> [::1]:22181}
2022/09/03 03:24:25 [ZKClient] connected to [::1]:22181
2022/09/03 03:24:25 [EventLoop] EventOccur: {EventSession StateConnected  <nil> [::1]:22181}
2022/09/03 03:24:26 [ZKClient] authenticated: id=144115953666293760, timeout=4000
2022/09/03 03:24:26 [ZKClient] re-submitting `0` credentials after reconnect
2022/09/03 03:24:26 [EventLoop] EventOccur: {EventSession StateHasSession  <nil> [::1]:22181}
2022/09/03 03:24:26 '/MyFirstZnode' exists: true, stat: &{4294967306 4294967306 1662136980002 1662136980002 0 0 0 0 6 0 4294967306}
2022/09/03 03:24:26 '/MyFirstZnode' create result: /MyFirstZnode
2022/09/03 03:24:26 Stopping event loop
2022/09/03 03:24:26 [ZKClient] recv loop terminated: EOF
2022/09/03 03:24:26 [ZKClient] send loop terminated: <nil>
```

---  

# Kafka

See [./composes/kafka](./composes/kafka) for details

> ### docker-compose.yml

```yaml
version: '3.4'

services:
  zoo1:
    image: zookeeper:latest
    restart: always
    container_name: zoo1
    ports:
      - "2181:2181"
    environment:
      - ZOO_MY_ID=1
      - ZOO_SERVERS=server.1=0.0.0.0:2888:3888;2181
  kafka1:
    image: confluentinc/cp-kafka:5.2.1
    hostname: kafka1
    ports:
      - "9092:9092"
    container_name: kafka1
    environment:
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka1:19092,LISTENER_DOCKER_EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181"
      KAFKA_BROKER_ID: 1
      KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    volumes:
      - ./__kafka_data/kafka1/data:/var/lib/kafka/data
    depends_on:
      - zoo1
  kafdrop:
    image: obsidiandynamics/kafdrop
    restart: "no"
    hostname: kafdrop
    container_name: kafdrop
    ports:
      - "9000:9000"
    environment:
      KAFKA_BROKERCONNECT: "kafka1:19092"
      JVM_OPTS: "-Xms16M -Xmx48M -Xss180K -XX:-TieredCompilation -XX:+UseStringDeduplication -noverify"
    depends_on:
      - "kafka1"
```

> ### Start containers and Checks

```shell
$ cd composes/kafka
$ docker-compose up

$ go run main.go
2022/09/03 02:25:02 Skip to create a new topic sample-message because already exists
2022/09/03 02:25:05 [Consumer] Setup Session. memberid:sarama-349f44c9-f0b8-477a-9af2-85a272281ead
2022/09/03 02:25:05 consume message: message-1
2022/09/03 02:25:05 consume message: message-2
2022/09/03 02:25:05 consume message: message-3
2022/09/03 02:25:05 consume message: message-4
2022/09/03 02:25:05 consume message: message-5
```

Check kafka topics from kafdrop([http://localhost:9000](http://localhost:9000) in your browser).

---  

# Kafka Cluster

See [./composes/kafkacluster](./composes/kafkacluster) for details

> ### docker-compose.yml

```yaml
version: '3.4'

services:
  zoo1:
    image: zookeeper:latest
    container_name: zoo1
    hostname: zoo1
    ports:
      - "2181:2181"
    environment:
      - ZOO_MY_ID=1
      - ZOO_SERVERS=server.1=0.0.0.0:2888:3888;2181
  kafka1:
    image: confluentinc/cp-kafka:5.2.1
    hostname: kafka1
    ports:
      - "9092:9092"
    container_name: kafka1
    restart: always
    environment:
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka1:19092,LISTENER_DOCKER_EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181"
      KAFKA_BROKER_ID: 1
      KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    volumes:
      - ./__kafka1_data:/var/lib/kafka/data
    depends_on:
      - zoo1
  kafka2:
    image: confluentinc/cp-kafka:5.2.1
    hostname: kafka2
    ports:
      - "9093:9093"
    container_name: kafka2
    restart: always
    environment:
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka2:19093,LISTENER_DOCKER_EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:9093
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181"
      KAFKA_BROKER_ID: 2
      KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
    volumes:
      - ./__kafka2_data:/var/lib/kafka/data
    depends_on:
      - zoo1
  kafka3:
    image: confluentinc/cp-kafka:5.2.1
    hostname: kafka3
    ports:
      - "9094:9094"
    container_name: kafka3
    restart: always
    environment:
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka3:19094,LISTENER_DOCKER_EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:9094
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_ZOOKEEPER_CONNECT: "zoo1:2181"
      KAFKA_BROKER_ID: 3
      KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
    volumes:
      - ./__kafka3_data:/var/lib/kafka/data
    depends_on:
      - zoo1
  kafdrop:
    image: obsidiandynamics/kafdrop
    restart: "no"
    ports:
      - "9000:9000"
    container_name: kafdrop
    environment:
      KAFKA_BROKERCONNECT: "kafka1:19092"
      JVM_OPTS: "-Xms16M -Xmx48M -Xss180K -XX:-TieredCompilation -XX:+UseStringDeduplication -noverify"
    depends_on:
      - "kafka1"
```

> ### Start containers and Checks

```shell
$ cd composes/kafkacluster
$ docker-compose up

$ go run main.go
2022/09/03 02:41:19 Success to create a new topic: sample-message
2022/09/03 02:41:32 [Consumer] Setup Session. memberid:sarama-338051c9-1d6a-44e7-8926-e2470ce6fcb4
2022/09/03 02:41:33 consume message: message-1
2022/09/03 02:41:33 consume message: message-2
2022/09/03 02:41:33 consume message: message-3
2022/09/03 02:41:33 consume message: message-4
2022/09/03 02:41:33 consume message: message-5
```

Check kafka brokers and topics from kafdrop (Connect to http://localhost:9000 in your browser).

---  

# RabbitMQ

See [./composes/rabbitmq](./composes/rabbitmq) for details

```yaml  
version: '3'

services:
  rabbitmq:
    image: "rabbitmq:3-management"
    hostname: "rabbit"
    ports:
      - "15672:15672"
      - "5672:5672"
    labels:
      NAME: "rabbitmq"
    volumes:
      - ./rabbitmq-isolated.conf:/etc/rabbitmq/rabbitmq.config
```  

> ### rabbitmq-isolated.conf

```
[
 {rabbit,
  [
   %% The default "guest" user is only permitted to access the server
   %% via a loopback interface (e.g. localhost).
   %% {loopback_users, [<<"guest">>]},
   %%
   %% Uncomment the following line if you want to allow access to the
   %% guest user from anywhere on the network.
   {loopback_users, []},
   {default_vhost,       "/"},
   {default_user,        "user"},
   {default_pass,        "secret"},
   {default_permissions, [".*", ".*", ".*"]}
  ]}
].
```

> ### Start containers and Checks

```shell
$ cd cd composes/rabbitmq
$ docker-compose up -d
```

---  

# Redis

See [./composes/redis](./composes/redis)

> ### docker-compose.yml

```yaml
version: '3.1'

services:
  redis:
    image: redis:latest
    ports:
      - "6379:6379"
    restart: always
```

```yaml
$ docker exec -it redis /usr/local/bin/redis-cli INCR mycounter
(integer) 1
```

> ### Start containers and Checks

```shell
$ cd composes/redis 
$ docker-compose up -d

$ docker exec -it redis /usr/local/bin/redis-cli INCR mycounter
(integer) 1
```

---  

# Redis Cluster
; TBD

# Jaeger
; TBD