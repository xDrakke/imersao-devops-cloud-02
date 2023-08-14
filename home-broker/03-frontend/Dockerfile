FROM node:20-slim
ARG next_url=localhost:3000
ARG next_url_api=localhost:3001
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
ENV NEXT_PUBLIC_NEST_URL=${next_url}
ENV NEXT_PUBLIC_NEST_URL_API=${next_url_api}
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]