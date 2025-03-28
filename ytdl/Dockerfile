FROM python:3.13-alpine

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

RUN apk update && \
    apk add --no-cache ffmpeg

COPY requirements.txt /app/requirements.txt
RUN pip install --upgrade pip && \
    pip install -r /app/requirements.txt

COPY . /app
WORKDIR /app

EXPOSE 8000
CMD [ "gunicorn", "--timeout", "600", "-w", "8", "-b", "0.0.0.0:8000", "main:app" ]
