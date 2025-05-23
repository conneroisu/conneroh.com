# Use Ubuntu as base image for better compatibility with GitHub Codespaces
FROM ubuntu:24.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Install system dependencies needed for Nix and development
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    sudo \
    xz-utils \
    ca-certificates \
    build-essential \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user for development
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Switch to the non-root user
USER $USERNAME
WORKDIR /home/$USERNAME

# Install Nix using the Determinate Nix Installer (more reliable for containers)
RUN curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install linux \
    --extra-conf "sandbox = false" \
    --extra-conf "experimental-features = nix-command flakes" \
    --extra-conf "filter-syscalls = false" \
    --init none \
    --no-confirm

# Source Nix environment
ENV PATH="/home/$USERNAME/.nix-profile/bin:/nix/var/nix/profiles/default/bin:$PATH"
ENV NIX_PROFILES="/nix/var/nix/profiles/default /home/$USERNAME/.nix-profile"
ENV NIX_SSL_CERT_FILE="/etc/ssl/certs/ca-certificates.crt"

# Copy flake files and other necessary project files
COPY --chown=$USERNAME:$USERNAME flake.nix flake.lock /workspace/
COPY --chown=$USERNAME:$USERNAME flake-scripts.nix /workspace/
COPY --chown=$USERNAME:$USERNAME go.mod go.sum /workspace/

# Set workspace as working directory
WORKDIR /workspace

# Pre-build the development environment to speed up container startup
RUN . /home/$USERNAME/.nix-profile/etc/profile.d/nix.sh && \
    nix develop .#devcontainer --command echo "Development environment built successfully"

# Create a script to enter the development shell
RUN echo '#!/bin/bash\n\
. /home/$USERNAME/.nix-profile/etc/profile.d/nix.sh\n\
cd /workspace\n\
exec nix develop .#devcontainer --command "$@"' > /home/$USERNAME/.local/bin/dev-shell && \
    chmod +x /home/$USERNAME/.local/bin/dev-shell

# Create .bashrc that automatically enters the dev shell
RUN echo '# Source Nix environment\n\
if [ -f /home/'$USERNAME'/.nix-profile/etc/profile.d/nix.sh ]; then\n\
    . /home/'$USERNAME'/.nix-profile/etc/profile.d/nix.sh\n\
fi\n\
\n\
# Auto-enter development shell if in workspace\n\
if [ "$PWD" = "/workspace" ] && [ -z "$IN_NIX_SHELL" ]; then\n\
    echo "Entering Nix development shell..."\n\
    exec nix develop .#devcontainer\n\
fi' >> /home/$USERNAME/.bashrc

# Create a script to manually enter dev shell
RUN echo '#!/bin/bash\n\
cd /workspace\n\
exec /home/'$USERNAME'/.local/bin/dev-shell bash' > /usr/local/bin/dev && \
    sudo chmod +x /usr/local/bin/dev

# Set environment variables for the container
ENV REPO_ROOT=/workspace
ENV IN_CONTAINER=true

# Create local bin directory and add to PATH
RUN mkdir -p /home/$USERNAME/.local/bin
ENV PATH="/home/$USERNAME/.local/bin:$PATH"

# Expose common development ports
EXPOSE 8080 3000 5173 8000

# Set the default command
CMD ["/bin/bash"]

# Health check to ensure Nix is working
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD . /home/$USERNAME/.nix-profile/etc/profile.d/nix.sh && nix --version || exit 1
