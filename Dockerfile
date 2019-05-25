FROM scratch

COPY ./dist /bin

CMD [ "/bin/mona" ]