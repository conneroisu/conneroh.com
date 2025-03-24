# Scripts for development shell
{lib}: {
  dx = {
    exec = ''$EDITOR $REPO_ROOT/flake.nix'';
    description = "Edit the flake.nix";
  };
  clean = {
    exec = ''
      git clean -fdx
    '';
    description = "Clean Project";
  };
  tests = {
    exec = ''go test -v -short ./...'';
    description = "Run go tests with short flag";
  };
  unit-tests = {
    exec = ''go test -v ./...'';
    description = "Run all go tests";
  };
  lint = {
    exec = ''golangci-lint run'';
    description = "Run golangci-lint";
  };
  build = {
    exec = ''nix build .#packages.x86_64-linux.conneroh'';
    description = "Build the package";
  };
  update = {
    exec = ''go run $REPO_ROOT/cmd/update'';
    description = "Run update command with doppler";
  };
  restart = {
    exec = ''rm -f $REPO_ROOT/master.db && go run $REPO_ROOT/cmd/update'';
    description = "Run restart command with doppler";
  };
  generate-all-profile = {
    exec = ''
      TIMEFORMAT="%R seconds"

      echo "Starting build process..."

      # Profile go generate
      echo "Running go generate..."
      START_TIME=$(date +%s.%N)
      go generate $REPO_ROOT/... &
      GO_PID=$!

      # Profile templ generate
      echo "Running templ generate..."
      TEMPL_START_TIME=$(date +%s.%N)
      templ generate $REPO_ROOT &
      TEMPL_PID=$!

      # Profile bun build
      echo "Running bun build..."
      BUN_START_TIME=$(date +%s.%N)
      bun build \
          $REPO_ROOT/index.js \
          --minify \
          --minify-syntax \
          --minify-whitespace \
          --minify-identifiers \
          --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &
      BUN_PID=$!

      # Profile tailwindcss
      echo "Running tailwindcss..."
      TAILWIND_START_TIME=$(date +%s.%N)
      tailwindcss \
          --minify \
          -i $REPO_ROOT/input.css \
          -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css \
          --cwd $REPO_ROOT &
      TAILWIND_PID=$!

      # Wait for go generate to complete
      wait $GO_PID
      GO_END_TIME=$(date +%s.%N)
      GO_DURATION=$(echo "$GO_END_TIME - $START_TIME" | bc)
      echo "go generate completed in $GO_DURATION seconds"

      # Wait for templ generate to complete
      wait $TEMPL_PID
      TEMPL_END_TIME=$(date +%s.%N)
      TEMPL_DURATION=$(echo "$TEMPL_END_TIME - $TEMPL_START_TIME" | bc)
      echo "templ generate completed in $TEMPL_DURATION seconds"

      # Wait for bun build to complete
      wait $BUN_PID
      BUN_END_TIME=$(date +%s.%N)
      BUN_DURATION=$(echo "$BUN_END_TIME - $BUN_START_TIME" | bc)
      echo "bun build completed in $BUN_DURATION seconds"

      # Wait for tailwindcss to complete
      wait $TAILWIND_PID
      TAILWIND_END_TIME=$(date +%s.%N)
      TAILWIND_DURATION=$(echo "$TAILWIND_END_TIME - $TAILWIND_START_TIME" | bc)
      echo "tailwindcss completed in $TAILWIND_DURATION seconds"

      # Calculate total duration
      TOTAL_DURATION=$(echo "$TAILWIND_END_TIME - $START_TIME" | bc)
      echo "Total build process completed in $TOTAL_DURATION seconds"

      echo "Build process finished!"
    '';
    description = "Generate all with profiling information";
  };
  "generate-reload" = {
    exec = ''
      templ generate $REPO_ROOT &
      wait
    '';
    description = "Generate templ files and wait for completion";
  };
  "generate-all" = {
    exec = ''
      go generate $REPO_ROOT/... &

      templ generate $REPO_ROOT &

      bun build \
          $REPO_ROOT/index.js \
          --minify \
          --minify-syntax \
          --minify-whitespace  \
          --minify-identifiers \
          --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &

      tailwindcss \
          --minify \
          -i $REPO_ROOT/input.css \
          -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css \
          --cwd $REPO_ROOT &

      wait
    '';
    description = "Generate all files in parallel";
  };
  format = {
    exec = ''
      export REPO_ROOT=$(git rev-parse --show-toplevel) # needed

      go fmt $REPO_ROOT/...

      git ls-files \
          --others \
          --exclude-standard \
          --cached \
          -- '*.cc' '*.h' '*.proto' \
          | xargs clang-format -i

      git ls-files \
        --others \
        --exclude-standard \
        --cached \
        -- '*.js' '*.ts' '*.css' '*.md' '*.json' \
        | xargs prettier --write

      golines \
        -l \
        -w \
        --max-len=80 \
        --shorten-comments \
        --ignored-dirs=.devenv .
    '';
    description = "Format code files";
  };
  run = {
    exec = ''air'';
    description = "Run the application with air for hot reloading";
  };
}
