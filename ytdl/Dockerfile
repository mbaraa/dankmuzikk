FROM python:3.11-alpine AS build

WORKDIR /app

ADD requirements.txt /tmp
RUN apk install mariadb && \
    pip install --upgrade pip && \
    pip install -r /tmp/requirements.txt

COPY . /app


FROM python:3.11-alpine

RUN apk add ffmpeg mariadb

WORKDIR /app

COPY --from=build /usr/local/lib/python3.11/site-packages  /usr/local/lib/python3.11/site-packages
COPY --from=build /usr/local/bin/ /usr/local/bin/
COPY --from=build /app .

EXPOSE 8000
CMD [ "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000" ]

