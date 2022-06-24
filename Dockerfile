FROM scratch

COPY ./bootpeer /bootpeer

ENTRYPOINT [ "/bootpeer" ] 
