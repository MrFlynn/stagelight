FROM node:latest

# Create workspace and volumes.
RUN mkdir /workspace
VOLUME /output

# Install dependencies and build static copy of site.
WORKDIR /workspace
COPY . .
RUN npm install
RUN npm run build

# Copy output to new directory
CMD ["cp", "-a", "/workspace/dist/.", "/output"]
