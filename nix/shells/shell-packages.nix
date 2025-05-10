{
  self',
  pkgs,
  scripts,
  scriptPackages,
}: let
  internal_port = 8080;
  force_https = true;
  processes = ["app"];
  version = self'.shortRev or "dirty";
  src = ./../../.;

  databaseFiles = pkgs.runCommand "database-files" {} ''
    mkdir -p $out/root
    cp ${./../../master.db} $out/root/master.db
  '';

  preBuild = ''
    ${pkgs.templ}/bin/templ generate
    ${pkgs.tailwindcss}/bin/tailwindcss \
        --minify \
        -i ./input.css \
        -o ./cmd/conneroh/_static/dist/style.css \
        --cwd .
  '';

  settingsFormat = pkgs.formats.toml {};

  flyDevConfig = {
    app = "conneroh-com-dev";
    primary_region = "ord";
    build = {};
    http_service = {
      inherit internal_port force_https processes;
      auto_stop_machines = "stop";
      auto_start_machines = true;
      min_machines_running = 0;
    };
    vm = [
      {
        memory = "512M";
        cpu_kind = "shared";
        cpus = 1;
      }
    ];
  };

  flyProdConfig = {
    app = "conneroh-com";
    primary_region = "ord";
    build = {};
    http_service = {
      inherit internal_port force_https processes;
      auto_stop_machines = "stop";
      auto_start_machines = true;
      min_machines_running = 0;
    };
    vm = [
      {
        memory = "1gb";
        cpu_kind = "shared";
        cpus = 2;
      }
    ];
  };

  flyDevToml = settingsFormat.generate "fly.dev.toml" flyDevConfig;
  flyProdToml = settingsFormat.generate "fly.toml" flyProdConfig;
in
  {
    conneroh = pkgs.buildGoModule {
      inherit src version preBuild;
      vendorHash = "sha256-kOGauV5dMTcHvSR7uWvY1dcKR4WqlWccDfnXtycsRVI=";
      name = "conneroh.com";
      goSum = ./../../go.sum;
      subPackages = ["."];
    };
    update = pkgs.buildGoModule {
      inherit src version preBuild;
      vendorHash = "sha256-kOGauV5dMTcHvSR7uWvY1dcKR4WqlWccDfnXtycsRVI=";
      goSum = ./../../go.sum;
      name = "update";
      subPackages = ["./cmd/update"];
      doCheck = false;
    };
    C-conneroh = pkgs.dockerTools.buildImage {
      created = "now";
      tag = "latest";
      name = "conneroh";
      config = {
        WorkingDir = "/root";
        Cmd = ["/bin/conneroh.com"];
        ExposedPorts = {"8080/tcp" = {};};
        Env = [
          "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
          "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
        ];
      };
      copyToRoot = [
        self'.packages.conneroh
        databaseFiles
      ];
    };
    deployPackage = pkgs.writeShellApplication {
      name = "deployPackage";
      runtimeInputs = with pkgs; [
        doppler
        skopeo
        flyctl
        cacert
      ];
      bashOptions = ["errexit" "pipefail"];
      text = ''
        set -e
        arg=$1
        TOKEN=""
        FLY_NAME=""
        CONFIG_FILE=""

        [ -z "$arg" ] && { echo "No argument provided. Please provide a target environment. (dev or prod)"; exit 1; }

        if [ "$arg" = "dev" ]; then
          [ -z "$FLY_DEV_AUTH_TOKEN" ] && FLY_DEV_AUTH_TOKEN="$(doppler secrets get --plain FLY_DEV_AUTH_TOKEN)"
          TOKEN="$FLY_DEV_AUTH_TOKEN"
          export FLY_NAME="conneroh-com-dev"
          export CONFIG_FILE=${flyDevToml}
        else
          [ -z "$FLY_AUTH_TOKEN" ] && FLY_AUTH_TOKEN="$(doppler secrets get --plain FLY_AUTH_TOKEN)"
          TOKEN="$FLY_AUTH_TOKEN"
          export FLY_NAME="conneroh-com"
          export CONFIG_FILE=${flyProdToml}
        fi

        REGISTRY="registry.fly.io/$FLY_NAME"
        echo "Copying image to Fly.io... to $REGISTRY"

        skopeo copy \
          --insecure-policy \
          docker-archive:"${self'.packages.C-conneroh}" \
          "docker://$REGISTRY:latest" \
          --dest-creds x:"$TOKEN" \
          --format v2s2

        echo "Deploying to Fly.io..."
        fly deploy \
          --remote-only \
          -c "$CONFIG_FILE" \
          -i "$REGISTRY:latest" \
          -t "$TOKEN"
      '';
    };
  }
  // pkgs.lib.genAttrs (builtins.attrNames scripts) (name: scriptPackages.${name})
