FROM bufbuild/buf

WORKDIR /repo

RUN apk add --no-cache shadow python3 py3-pip nodejs npm
RUN python -m pip install --ignore-installed betterproto[compiler]==2.0.0b6
RUN npm install -g @bufbuild/protoc-gen-es@^1

# From: https://jtreminio.com/blog/running-docker-containers-as-current-host-user/
ARG USER_ID=1000
ARG GROUP_ID=1000
RUN if [ ${USER_ID:-0} -ne 0 ] && [ ${GROUP_ID:-0} -ne 0 ]; then \
    if getent group www-data ; then groupdel www-data; fi &&\
    groupadd -o -g ${GROUP_ID} -r www-data &&\
    useradd -l -u ${USER_ID} -g www-data www-data &&\
    install -d -m 0755 -o www-data -g www-data /home/www-data &&\
    chown --changes --silent --no-dereference --recursive \
          ${USER_ID}:${GROUP_ID} \
        /home/www-data \
        /repo \
;fi

USER www-data

ENTRYPOINT [ "/repo/protocol/generate.sh" ]
