FROM python:3.9-slim-buster AS build
WORKDIR /tests
COPY ./e2e/requirements.txt .
RUN pip3 install --compile -r requirements.txt && rm -rf /root/.cache
