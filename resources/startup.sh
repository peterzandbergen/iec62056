apt-get update -y && \
apt-get upgrade -y && \
apt-get install curl git -y && \
curl https://dl.google.com/go/go1.14.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
tar -C /usr/local -xzf /tmp/go.tar.gz && \
export PATH=$PATH:/usr/local/go/bin && \
mkdir -p /home/emserver && \
cd /home/emserver && \
git clone https://github.com/peterzandbergen/iec62056.git /home/emserver/iec62056 && \
cd /home/emserver/iec62056 && \
git checkout dev/export-via-rest && \
mkdir -p /home/emserver/emlog-db && \
gsutil -m cp -r gs://emeter-db-backup/emlog-db/* /home/emserver/emlog-db/ && \
echo done
nohup go run cmd/emserver/main.go --local-cache-path /home/emserver/emlog-db --port 80  2>&1 
