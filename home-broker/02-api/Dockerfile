FROM node:18 
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY prisma ./prisma
RUN npx prisma generate
#RUN npx prisma migrate dev  
COPY . .
RUN npm run build
EXPOSE 3000
RUN chmod +x ./entrypoint.sh 
CMD [ "./entrypoint.sh" ]
# CMD [ "node", "dist/main.js" ]