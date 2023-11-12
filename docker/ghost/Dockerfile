ARG VERSION
FROM node:18 as s3-builder
WORKDIR /root
RUN npm install ghost-storage-adapter-s3@2.8.0

FROM ghost:$VERSION
COPY --from=s3-builder /root/node_modules/ghost-storage-adapter-s3/ /var/lib/ghost/content/adapters/storage/s3/
RUN npm install ghost-storage-base@1.0.0 && npm install aws-sdk@2.1474.0
ENV storage__active s3
