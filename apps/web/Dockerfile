FROM node:22-alpine

WORKDIR /src/apps/web
COPY apps/web/package*.json ./
RUN npm ci
COPY apps/web/ ./
COPY contracts/ /src/contracts/
COPY data/ /src/data/
ENV HOSTNAME=0.0.0.0
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--hostname", "0.0.0.0"]
