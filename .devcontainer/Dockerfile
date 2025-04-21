
  # Reuse a base image made for devcontainers.
FROM mcr.microsoft.com/vscode/devcontainers/base:buster

# These dependencies are required by Nix.
RUN apt update -y
RUN apt -y install --no-install-recommends curl xz-utils
RUN apt update -y
RUN apt install curl -y
RUN curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install linux \
  --extra-conf "sandbox = false" \
  --init none \
  --no-confirm
ENV PATH="${PATH}:/nix/var/nix/profiles/default/bin"
# Configuration for Nix from the repository shared amongst developers.
RUN mkdir -p /etc/nix
COPY nix.conf /etc/nix/nix.conf

# This loads the development environment when the container is started.
COPY profile.sh /etc/profile.d/devcontainer.sh

ENV USER=vscode

# Rebuilding the container should not throw the Nix store away so we make /nix
# a volume. See `devcontainer.json` later.
VOLUME /nix