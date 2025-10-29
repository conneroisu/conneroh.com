# FROM directus/directus:11
# USER root
# # RUN apk add ca-certificates fuse3 sqlite bash curl libspatialite sqlite-libs
# # COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs
# # ARG LIFTFS_CONFIG=litefs.yml
# # ADD ./litefs.yml /tmp/litefs.yml
# # ADD ./setup-db.sh /tmp/setup-db.sh
# # RUN cp /tmp/litefs.yml /etc/litefs.yml
# # ENV DB_FILENAME=/litefs/db.sqlite
# # ENV SQLITE_EXTENSIONS="mod_spatialite"
# RUN apk add curl
# RUN curl --create-dirs -o $HOME/.postgresql/root.crt 'https://cockroachlabs.cloud/clusters/83c093f7-89ef-4461-bca3-eeffb41f5f9b/cert'
# # CMD : pm2-runtime start ecosystem.config.cjs 
# base node image
FROM node:16-bullseye-slim as base

# set for base and all layer that inherit from it
ENV NODE_ENV production

RUN apt-get update && apt-get install -y openssl sqlite3

# Install all node_modules, including dev dependencies
FROM base as deps

WORKDIR /myapp

ADD package.json .npmrc ./
RUN npm install --production=false

# Setup production node_modules
FROM base as production-deps

WORKDIR /myapp

COPY --from=deps /myapp/node_modules /myapp/node_modules
ADD package.json .npmrc ./
RUN npm prune --production

# Finally, build the production image with minimal footprint
FROM base

ENV DATABASE_URL=file:/data/database/data.db
ENV PORT="6055"
ENV NODE_ENV="production"

# add shortcut for connecting to database CLI
RUN echo "#!/bin/sh\nset -x\nsqlite3 \$DATABASE_URL" > /usr/local/bin/database-cli && chmod +x /usr/local/bin/database-cli

WORKDIR /myapp

COPY --from=production-deps /myapp/node_modules /myapp/node_modules

ADD . .

CMD ["bash", "start.sh"]
