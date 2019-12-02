FROM node:latest

# Create workspace and volumes.
VOLUME /output
RUN mkdir /workspace

# Install dependencies and build static copy of site.
WORKDIR /workspace
COPY . .
RUN npm install
RUN npm run build

# Copy the contents of the dist folder to the output folder.
COPY /workspace/dist /output