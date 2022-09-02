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
# RUN apt-get install --assume-yes  apt-transport-https  ca-certificates  curl  gnupg-agent  software-properties-common && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && apt-get update && apt-get install --assume-yes docker-ce docker-ce-cli containerd.io

# Install docker-compose (TODO : dependencies)
# RUN curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]