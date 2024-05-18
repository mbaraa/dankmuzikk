FROM python:3.11-slim-bookworm AS build

WORKDIR /app

ADD requirements.txt /tmp
RUN apt-get update -y && \
    apt-get install -y gcc libmariadb-dev mariadb-client && \
    pip install --upgrade pip && \
    pip install --no-cache-dir -r /tmp/requirements.txt

COPY . /app


FROM python:3.11-alpine

FROM python:3.11-slim-bookworm

RUN apt-get update -y && \
    apt-get install -y mariadb-client ffmpeg

WORKDIR /app

COPY --from=build /usr/local/lib/python3.11/site-packages  /usr/local/lib/python3.11/site-packages
COPY --from=build /usr/local/bin/ /usr/local/bin/
COPY --from=build /app .

EXPOSE 8000
CMD [ "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000" ]

