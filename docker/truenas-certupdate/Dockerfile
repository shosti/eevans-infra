FROM python:3.9-alpine3.12

COPY entrypoint.sh /entrypoint.sh

RUN wget -O /deploy-freenas.py \
            https://raw.githubusercontent.com/danb35/deploy-freenas/de822753c47fb801075e407f11146ac68cbd7ad7/deploy_freenas.py \
         && chmod +x /deploy-freenas.py \
         && chmod +x /entrypoint.sh \
         && pip install requests \
         && addgroup -S certupdate \
         && adduser -S certupdate -G certupdate

USER certupdate
ENTRYPOINT /entrypoint.sh
