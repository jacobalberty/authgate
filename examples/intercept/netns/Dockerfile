FROM debian:12-slim
RUN apt-get update && apt-get install -y iptables && rm -rf /var/lib/apt/lists/*
COPY intercept-rules.sh /intercept-rules.sh
ENTRYPOINT ["/bin/sh","/intercept-rules.sh"]