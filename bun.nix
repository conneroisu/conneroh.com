# This file was autogenerated by `bun2nix`, editing it is not recommended.
# Consume it with `callPackage` in your actual derivation -> https://nixos-and-flakes.thiscute.world/nixpkgs/callpackage
{
  lib,
  fetchurl,
  runCommand,
  libarchive,
  bun,
  makeWrapper,
  ...
}: let
  # Set of Bun packages to install
  packages = {
    "@alpinejs/anchor" = {
      out_path = "@alpinejs/anchor";
      binaries = {
      };
      pkg = fetchurl {
        name = "@alpinejs/anchor@3.14.9";
        url = "https://registry.npmjs.org/@alpinejs/anchor/-/anchor-3.14.9.tgz";
        hash = "sha256-kMQLV3ItTtpqwPvoXEn8CTmU7EaiS+mKchvLTOnByEc=";
      };
    };
    "@alpinejs/intersect" = {
      out_path = "@alpinejs/intersect";
      binaries = {
      };
      pkg = fetchurl {
        name = "@alpinejs/intersect@3.14.9";
        url = "https://registry.npmjs.org/@alpinejs/intersect/-/intersect-3.14.9.tgz";
        hash = "sha256-wFtr78dKhPhSrwTlcxqBQ+6nsZcPrW4OZo1Ge6mhVOc=";
      };
    };
    "@parcel/watcher" = {
      out_path = "@parcel/watcher";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher/-/watcher-2.5.1.tgz";
        hash = "sha256-v3prVXeihxU8mmv3vpU1VrmY8qV3But6U3HsfrAIjUE=";
      };
    };
    "@parcel/watcher-android-arm64" = {
      out_path = "@parcel/watcher-android-arm64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-android-arm64@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-android-arm64/-/watcher-android-arm64-2.5.1.tgz";
        hash = "sha256-MCubTXv4PRJUN3styKSwjYVKFxp5YOFZVdwN2flNTMg=";
      };
    };
    "@parcel/watcher-darwin-arm64" = {
      out_path = "@parcel/watcher-darwin-arm64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-darwin-arm64@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-darwin-arm64/-/watcher-darwin-arm64-2.5.1.tgz";
        hash = "sha256-hLCFsjrLwDaQSyrlRXRI3KwA2cW5GQSljSFPALc97Vg=";
      };
    };
    "@parcel/watcher-darwin-x64" = {
      out_path = "@parcel/watcher-darwin-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-darwin-x64@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-darwin-x64/-/watcher-darwin-x64-2.5.1.tgz";
        hash = "sha256-X0lP9VLyyaanlSSJArfXQ4K+zsENJdRMHz1ppp2ptsc=";
      };
    };
    "@parcel/watcher-freebsd-x64" = {
      out_path = "@parcel/watcher-freebsd-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-freebsd-x64@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-freebsd-x64/-/watcher-freebsd-x64-2.5.1.tgz";
        hash = "sha256-7WsOuklZyfJ1W7srPVjoLmHE0ECfwUlWtLerC90TldY=";
      };
    };
    "@parcel/watcher-linux-arm-glibc" = {
      out_path = "@parcel/watcher-linux-arm-glibc";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-linux-arm-glibc@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-linux-arm-glibc/-/watcher-linux-arm-glibc-2.5.1.tgz";
        hash = "sha256-pi5qYiT5LQHtlPijtZEMAaxNwHUWoL4mqyAh0uw7QSc=";
      };
    };
    "@parcel/watcher-linux-arm-musl" = {
      out_path = "@parcel/watcher-linux-arm-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-linux-arm-musl@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-linux-arm-musl/-/watcher-linux-arm-musl-2.5.1.tgz";
        hash = "sha256-fofj1XGOmUXderYRdpV8/78ekV2OSe/FM3bAyGPJbvQ=";
      };
    };
    "@parcel/watcher-linux-arm64-glibc" = {
      out_path = "@parcel/watcher-linux-arm64-glibc";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-linux-arm64-glibc@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-linux-arm64-glibc/-/watcher-linux-arm64-glibc-2.5.1.tgz";
        hash = "sha256-FE5nlSjrk7lop/1MsoMLVnGXXXE+SZgHfGbLf5HpoHU=";
      };
    };
    "@parcel/watcher-linux-arm64-musl" = {
      out_path = "@parcel/watcher-linux-arm64-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-linux-arm64-musl@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-linux-arm64-musl/-/watcher-linux-arm64-musl-2.5.1.tgz";
        hash = "sha256-P3P2wUf4k0wqXB/7XrYdGE3O46Wmh2elDii4UPqy1eI=";
      };
    };
    "@parcel/watcher-linux-x64-glibc" = {
      out_path = "@parcel/watcher-linux-x64-glibc";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-linux-x64-glibc@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-linux-x64-glibc/-/watcher-linux-x64-glibc-2.5.1.tgz";
        hash = "sha256-2dmtfQHjsXbm9gP5uSRROhYYTsiw3ih4L45lnBE9UKA=";
      };
    };
    "@parcel/watcher-linux-x64-musl" = {
      out_path = "@parcel/watcher-linux-x64-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-linux-x64-musl@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-linux-x64-musl/-/watcher-linux-x64-musl-2.5.1.tgz";
        hash = "sha256-9/bt7x5rILicRz2GTDfIS8QfuMx0LT1/WDW9X23okDU=";
      };
    };
    "@parcel/watcher-win32-arm64" = {
      out_path = "@parcel/watcher-win32-arm64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-win32-arm64@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-win32-arm64/-/watcher-win32-arm64-2.5.1.tgz";
        hash = "sha256-gxoA8CiYPr2gX7O77pxoSYnuHomDcQeo8e6OvAkxFpY=";
      };
    };
    "@parcel/watcher-win32-ia32" = {
      out_path = "@parcel/watcher-win32-ia32";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-win32-ia32@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-win32-ia32/-/watcher-win32-ia32-2.5.1.tgz";
        hash = "sha256-B+M/7X4SPMJGg7uVdmjk2puLXONLzWbFUDEEkchWOAw=";
      };
    };
    "@parcel/watcher-win32-x64" = {
      out_path = "@parcel/watcher-win32-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@parcel/watcher-win32-x64@2.5.1";
        url = "https://registry.npmjs.org/@parcel/watcher-win32-x64/-/watcher-win32-x64-2.5.1.tgz";
        hash = "sha256-MFZk/ZVJGnPlGIBruXzesw3WIW2sLs6mvmEncFZ722c=";
      };
    };
    "@popperjs/core" = {
      out_path = "@popperjs/core";
      binaries = {
      };
      pkg = fetchurl {
        name = "@popperjs/core@2.11.8";
        url = "https://registry.npmjs.org/@popperjs/core/-/core-2.11.8.tgz";
        hash = "sha256-jgm9+pEgNWaOYs6mEyG84ny9ARuFZyBV2yXScb1jr0k=";
      };
    };
    "@tailwindcss/cli" = {
      out_path = "@tailwindcss/cli";
      binaries = {
        "tailwindcss" = "../@tailwindcss/cli/dist/index.mjs";
      };
      pkg = fetchurl {
        name = "@tailwindcss/cli@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/cli/-/cli-4.1.3.tgz";
        hash = "sha256-If9p+QJQP0NKQqTYKLcJ4XRgbAq7567zgDVyFjjTq8c=";
      };
    };
    "@tailwindcss/node" = {
      out_path = "@tailwindcss/node";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/node@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/node/-/node-4.1.3.tgz";
        hash = "sha256-HcbCe7PeeJiEl3pMDQimif6r8NRpoPCTl5onsPzBTfw=";
      };
    };
    "@tailwindcss/oxide" = {
      out_path = "@tailwindcss/oxide";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide/-/oxide-4.1.3.tgz";
        hash = "sha256-WXG2denX6p/edEiPaulalq4AxkfVz7STZ5hhN/qqQIQ=";
      };
    };
    "@tailwindcss/oxide-android-arm64" = {
      out_path = "@tailwindcss/oxide-android-arm64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-android-arm64@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-android-arm64/-/oxide-android-arm64-4.1.3.tgz";
        hash = "sha256-PPaI0jaDXb4K6pKlsXzU1ENBARCY2Kz6HvZSCqWLkPY=";
      };
    };
    "@tailwindcss/oxide-darwin-arm64" = {
      out_path = "@tailwindcss/oxide-darwin-arm64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-darwin-arm64@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-darwin-arm64/-/oxide-darwin-arm64-4.1.3.tgz";
        hash = "sha256-5dHqNUK/r5pj1LFByhXa1p6lqp3OiZOjDVdc1VklLRA=";
      };
    };
    "@tailwindcss/oxide-darwin-x64" = {
      out_path = "@tailwindcss/oxide-darwin-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-darwin-x64@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-darwin-x64/-/oxide-darwin-x64-4.1.3.tgz";
        hash = "sha256-8YMTri8GhfbmVtlG2GwsKcSzPRXTChb7oEuipc9+rjE=";
      };
    };
    "@tailwindcss/oxide-freebsd-x64" = {
      out_path = "@tailwindcss/oxide-freebsd-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-freebsd-x64@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-freebsd-x64/-/oxide-freebsd-x64-4.1.3.tgz";
        hash = "sha256-OAluVFvo6iK1IAYFm6fMsWs175+xMe9Wt9hnd2M57n4=";
      };
    };
    "@tailwindcss/oxide-linux-arm-gnueabihf" = {
      out_path = "@tailwindcss/oxide-linux-arm-gnueabihf";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-linux-arm-gnueabihf@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-linux-arm-gnueabihf/-/oxide-linux-arm-gnueabihf-4.1.3.tgz";
        hash = "sha256-CntN3CwK0Gx/RLl4ELVneVCUKY1QpJFygUNJ+uYX62I=";
      };
    };
    "@tailwindcss/oxide-linux-arm64-gnu" = {
      out_path = "@tailwindcss/oxide-linux-arm64-gnu";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-linux-arm64-gnu@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-linux-arm64-gnu/-/oxide-linux-arm64-gnu-4.1.3.tgz";
        hash = "sha256-pfnI/OPORRd/CRqOLN5WMYlN4EKC64l6XL2yyOtozcY=";
      };
    };
    "@tailwindcss/oxide-linux-arm64-musl" = {
      out_path = "@tailwindcss/oxide-linux-arm64-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-linux-arm64-musl@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-linux-arm64-musl/-/oxide-linux-arm64-musl-4.1.3.tgz";
        hash = "sha256-oG8/CCw5VeBgFfRBP6HXq1QAIWUSX7xE+OYMa+sZXqw=";
      };
    };
    "@tailwindcss/oxide-linux-x64-gnu" = {
      out_path = "@tailwindcss/oxide-linux-x64-gnu";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-linux-x64-gnu@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-linux-x64-gnu/-/oxide-linux-x64-gnu-4.1.3.tgz";
        hash = "sha256-t8r94AmUyjAKuF2ZGTelt0au/ils6y5So95+2C0jhgo=";
      };
    };
    "@tailwindcss/oxide-linux-x64-musl" = {
      out_path = "@tailwindcss/oxide-linux-x64-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-linux-x64-musl@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-linux-x64-musl/-/oxide-linux-x64-musl-4.1.3.tgz";
        hash = "sha256-z7y2T0+8AeoPGEmYH7iVSiRMdrsPcHIXUEFpw+4e1Xo=";
      };
    };
    "@tailwindcss/oxide-win32-arm64-msvc" = {
      out_path = "@tailwindcss/oxide-win32-arm64-msvc";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-win32-arm64-msvc@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-win32-arm64-msvc/-/oxide-win32-arm64-msvc-4.1.3.tgz";
        hash = "sha256-srrYcWONgaPIuXLvM3HhXTklC3YWLE2Gb209wDr+E0o=";
      };
    };
    "@tailwindcss/oxide-win32-x64-msvc" = {
      out_path = "@tailwindcss/oxide-win32-x64-msvc";
      binaries = {
      };
      pkg = fetchurl {
        name = "@tailwindcss/oxide-win32-x64-msvc@4.1.3";
        url = "https://registry.npmjs.org/@tailwindcss/oxide-win32-x64-msvc/-/oxide-win32-x64-msvc-4.1.3.tgz";
        hash = "sha256-9YaykHAiezigoRlyZGeUPkk4EIAI4xkRvaRoyJRxCuA=";
      };
    };
    "@types/alpinejs" = {
      out_path = "@types/alpinejs";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/alpinejs@3.13.11";
        url = "https://registry.npmjs.org/@types/alpinejs/-/alpinejs-3.13.11.tgz";
        hash = "sha256-R+5HQDLB86YXkbUgnDcm/SnHZSPjyfDyYNlVy7QuTik=";
      };
    };
    "@types/alpinejs__anchor" = {
      out_path = "@types/alpinejs__anchor";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/alpinejs__anchor@3.13.1";
        url = "https://registry.npmjs.org/@types/alpinejs__anchor/-/alpinejs__anchor-3.13.1.tgz";
        hash = "sha256-zbCEYGCu4viGfmgikq0FhXUtShH8sFetluWAzgHOM2Y=";
      };
    };
    "@types/alpinejs__intersect" = {
      out_path = "@types/alpinejs__intersect";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/alpinejs__intersect@3.13.4";
        url = "https://registry.npmjs.org/@types/alpinejs__intersect/-/alpinejs__intersect-3.13.4.tgz";
        hash = "sha256-aQjur+u1IxVcBEJJfyNg+pUHDJZJLsOCnt3c3abTxAM=";
      };
    };
    "@types/bun" = {
      out_path = "@types/bun";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/bun@1.2.9";
        url = "https://registry.npmjs.org/@types/bun/-/bun-1.2.9.tgz";
        hash = "sha256-676RZuumPabwPPYxsxXk87DmRsnz1bB/4GFOg0iXUVc=";
      };
    };
    "@types/mathjax" = {
      out_path = "@types/mathjax";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/mathjax@0.0.40";
        url = "https://registry.npmjs.org/@types/mathjax/-/mathjax-0.0.40.tgz";
        hash = "sha256-ZtfHomPxt+iAvfb8RnstCbWqWBggE6jEwGUtCZ3PEYU=";
      };
    };
    "@types/node" = {
      out_path = "@types/node";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/node@22.13.5";
        url = "https://registry.npmjs.org/@types/node/-/node-22.13.5.tgz";
        hash = "sha256-jtFd0kKkgHHRYLsT8YUk7PgWSyu0nVwrUHm+IMzNRKQ=";
      };
    };
    "@types/ws" = {
      out_path = "@types/ws";
      binaries = {
      };
      pkg = fetchurl {
        name = "@types/ws@8.5.14";
        url = "https://registry.npmjs.org/@types/ws/-/ws-8.5.14.tgz";
        hash = "sha256-SsUn5egoHI/5PcxHQwwtNupaBcwgePqjR3uTG7Ew37w=";
      };
    };
    "@vue/reactivity" = {
      out_path = "@vue/reactivity";
      binaries = {
      };
      pkg = fetchurl {
        name = "@vue/reactivity@3.1.5";
        url = "https://registry.npmjs.org/@vue/reactivity/-/reactivity-3.1.5.tgz";
        hash = "sha256-nNOO5PQ6ZwR1hbndoPvc/upHhUkPbmd4GWhBR6nUOcQ=";
      };
    };
    "@vue/shared" = {
      out_path = "@vue/shared";
      binaries = {
      };
      pkg = fetchurl {
        name = "@vue/shared@3.1.5";
        url = "https://registry.npmjs.org/@vue/shared/-/shared-3.1.5.tgz";
        hash = "sha256-n7bFNEvu8+hoGhl3GhMoeGs61+i0MelDpU+VTM6OYXs=";
      };
    };
    "alpinejs" = {
      out_path = "alpinejs";
      binaries = {
      };
      pkg = fetchurl {
        name = "alpinejs@3.14.9";
        url = "https://registry.npmjs.org/alpinejs/-/alpinejs-3.14.9.tgz";
        hash = "sha256-l9rXwMgeZZz8jncABV2pdw+Bhhh8uainbvtX4A1c5So=";
      };
    };
    "braces" = {
      out_path = "braces";
      binaries = {
      };
      pkg = fetchurl {
        name = "braces@3.0.3";
        url = "https://registry.npmjs.org/braces/-/braces-3.0.3.tgz";
        hash = "sha256-HNGOhiyGQLRWixQlp99O4DD/IB1FstqPnyItKYdJT/w=";
      };
    };
    "bun-types" = {
      out_path = "bun-types";
      binaries = {
      };
      pkg = fetchurl {
        name = "bun-types@1.2.9";
        url = "https://registry.npmjs.org/bun-types/-/bun-types-1.2.9.tgz";
        hash = "sha256-0yWVHVa+c0pdHtrZdc6wHT7SXgQyeTiR8uYXd7GVpdA=";
      };
    };
    "caniuse-lite" = {
      out_path = "caniuse-lite";
      binaries = {
      };
      pkg = fetchurl {
        name = "caniuse-lite@1.0.30001713";
        url = "https://registry.npmjs.org/caniuse-lite/-/caniuse-lite-1.0.30001713.tgz";
        hash = "sha256-xCfiOkGjX842sj7iFRloZcy84qREzvkDDCHkrPc9HWA=";
      };
    };
    "detect-libc" = {
      out_path = "detect-libc";
      binaries = {
        "detect-libc" = "../detect-libc/./bin/detect-libc.js";
      };
      pkg = fetchurl {
        name = "detect-libc@1.0.3";
        url = "https://registry.npmjs.org/detect-libc/-/detect-libc-1.0.3.tgz";
        hash = "sha256-eirymLhcs2+HdpMxjhV+TSBxnOjcoCs01owSFCtxIYs=";
      };
    };
    "enhanced-resolve" = {
      out_path = "enhanced-resolve";
      binaries = {
      };
      pkg = fetchurl {
        name = "enhanced-resolve@5.18.1";
        url = "https://registry.npmjs.org/enhanced-resolve/-/enhanced-resolve-5.18.1.tgz";
        hash = "sha256-NOdA+2c/KSvtms70xKxkms0BNurcM1ypBlsxDG93WzM=";
      };
    };
    "fill-range" = {
      out_path = "fill-range";
      binaries = {
      };
      pkg = fetchurl {
        name = "fill-range@7.1.1";
        url = "https://registry.npmjs.org/fill-range/-/fill-range-7.1.1.tgz";
        hash = "sha256-Gmw6NEbLHWk89rYnheL22PlZYwWv9HTdKzGy2NQnWTw=";
      };
    };
    "graceful-fs" = {
      out_path = "graceful-fs";
      binaries = {
      };
      pkg = fetchurl {
        name = "graceful-fs@4.2.11";
        url = "https://registry.npmjs.org/graceful-fs/-/graceful-fs-4.2.11.tgz";
        hash = "sha256-OWE3SqFh5v7YDW9Laq8vt+r9LJ5b406GWzdfARTdCZw=";
      };
    };
    "htmx-ext-preload" = {
      out_path = "htmx-ext-preload";
      binaries = {
      };
      pkg = fetchurl {
        name = "htmx-ext-preload@2.1.1";
        url = "https://registry.npmjs.org/htmx-ext-preload/-/htmx-ext-preload-2.1.1.tgz";
        hash = "sha256-4PbKzAZyxoJSwbmMLE2VdQh8mQwjVqOgqIEo6RCDK0k=";
      };
    };
    "htmx.org" = {
      out_path = "htmx.org";
      binaries = {
      };
      pkg = fetchurl {
        name = "htmx.org@2.0.4";
        url = "https://registry.npmjs.org/htmx.org/-/htmx.org-2.0.4.tgz";
        hash = "sha256-o0mI0bnwBaRY1ZOqjXSGU3y4eLisO4JwPbEmiYMSPsw=";
      };
    };
    "is-extglob" = {
      out_path = "is-extglob";
      binaries = {
      };
      pkg = fetchurl {
        name = "is-extglob@2.1.1";
        url = "https://registry.npmjs.org/is-extglob/-/is-extglob-2.1.1.tgz";
        hash = "sha256-jF1ChhRq1i/BCWmBcAzhwioWdwiSb8oB+cp0+btQvBk=";
      };
    };
    "is-glob" = {
      out_path = "is-glob";
      binaries = {
      };
      pkg = fetchurl {
        name = "is-glob@4.0.3";
        url = "https://registry.npmjs.org/is-glob/-/is-glob-4.0.3.tgz";
        hash = "sha256-P+RT+xk7tY9vBQXfsRUSMJNTgLW1Xh+YZCYcKq/BvsY=";
      };
    };
    "is-number" = {
      out_path = "is-number";
      binaries = {
      };
      pkg = fetchurl {
        name = "is-number@7.0.0";
        url = "https://registry.npmjs.org/is-number/-/is-number-7.0.0.tgz";
        hash = "sha256-e3XBBXGYz5dpaQmpvuF2ycXpvLWwO/Ps7y9ITe+t1R4=";
      };
    };
    "jiti" = {
      out_path = "jiti";
      binaries = {
        "jiti" = "../jiti/lib/jiti-cli.mjs";
      };
      pkg = fetchurl {
        name = "jiti@2.4.2";
        url = "https://registry.npmjs.org/jiti/-/jiti-2.4.2.tgz";
        hash = "sha256-HMryAyqVrmEG6ByGDmKKNzSaqQ8aQ478k3Xrh7I1LAs=";
      };
    };
    "lightningcss" = {
      out_path = "lightningcss";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss@1.29.2";
        url = "https://registry.npmjs.org/lightningcss/-/lightningcss-1.29.2.tgz";
        hash = "sha256-7KMxryRhM+0c0UlRQ8vaLkVG9JkW0lBWGikrqQEnI4Y=";
      };
    };
    "lightningcss-darwin-arm64" = {
      out_path = "lightningcss-darwin-arm64";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-darwin-arm64@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-darwin-arm64/-/lightningcss-darwin-arm64-1.29.2.tgz";
        hash = "sha256-CHAqvJGTFngargr3fTszw4VOIwHK7mVRj8R3vmDx838=";
      };
    };
    "lightningcss-darwin-x64" = {
      out_path = "lightningcss-darwin-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-darwin-x64@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-darwin-x64/-/lightningcss-darwin-x64-1.29.2.tgz";
        hash = "sha256-w+Pa1C1fInJtA6WGM80/+z9aSwOa1fGdycNecfY9DTs=";
      };
    };
    "lightningcss-freebsd-x64" = {
      out_path = "lightningcss-freebsd-x64";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-freebsd-x64@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-freebsd-x64/-/lightningcss-freebsd-x64-1.29.2.tgz";
        hash = "sha256-nBx8CWlmF/UPUu46CRq1UKCWaBghLQRnoKmBcjiZVj4=";
      };
    };
    "lightningcss-linux-arm-gnueabihf" = {
      out_path = "lightningcss-linux-arm-gnueabihf";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-linux-arm-gnueabihf@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-linux-arm-gnueabihf/-/lightningcss-linux-arm-gnueabihf-1.29.2.tgz";
        hash = "sha256-8C9ss25w1NvsvamP6UPKjbX9V8rcZxQ0HRBoO/sgXuY=";
      };
    };
    "lightningcss-linux-arm64-gnu" = {
      out_path = "lightningcss-linux-arm64-gnu";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-linux-arm64-gnu@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-linux-arm64-gnu/-/lightningcss-linux-arm64-gnu-1.29.2.tgz";
        hash = "sha256-+gOuP0f2WtLzP2abZ2naRivND9J39cuSvAq5k/304UM=";
      };
    };
    "lightningcss-linux-arm64-musl" = {
      out_path = "lightningcss-linux-arm64-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-linux-arm64-musl@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-linux-arm64-musl/-/lightningcss-linux-arm64-musl-1.29.2.tgz";
        hash = "sha256-p6FshGT4GjiRFcyM7bjQcbreDPDNiCsc+ykJzxIo5Q0=";
      };
    };
    "lightningcss-linux-x64-gnu" = {
      out_path = "lightningcss-linux-x64-gnu";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-linux-x64-gnu@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-linux-x64-gnu/-/lightningcss-linux-x64-gnu-1.29.2.tgz";
        hash = "sha256-QYuqh+bt1MZqzb3ERlmAO+poLCnB0YsYoq+Qos60zwQ=";
      };
    };
    "lightningcss-linux-x64-musl" = {
      out_path = "lightningcss-linux-x64-musl";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-linux-x64-musl@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-linux-x64-musl/-/lightningcss-linux-x64-musl-1.29.2.tgz";
        hash = "sha256-5ZkF0IouHSz5lHRQXpRHjKuCpKZsoNRBA/SOxEYMGCI=";
      };
    };
    "lightningcss-win32-arm64-msvc" = {
      out_path = "lightningcss-win32-arm64-msvc";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-win32-arm64-msvc@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-win32-arm64-msvc/-/lightningcss-win32-arm64-msvc-1.29.2.tgz";
        hash = "sha256-oYuE4SiLDlqPwr0xalxZHQYUOg5gWVHUW30CgVZc7N0=";
      };
    };
    "lightningcss-win32-x64-msvc" = {
      out_path = "lightningcss-win32-x64-msvc";
      binaries = {
      };
      pkg = fetchurl {
        name = "lightningcss-win32-x64-msvc@1.29.2";
        url = "https://registry.npmjs.org/lightningcss-win32-x64-msvc/-/lightningcss-win32-x64-msvc-1.29.2.tgz";
        hash = "sha256-zhE5YOvJLP9pTdtgfsrTc8wUHeTdJbH3eSnTsbHpVjA=";
      };
    };
    "lightningcss/detect-libc" = {
      out_path = "lightningcss/node_modules/detect-libc";
      binaries = {
      };
      pkg = fetchurl {
        name = "detect-libc@2.0.3";
        url = "https://registry.npmjs.org/detect-libc/-/detect-libc-2.0.3.tgz";
        hash = "sha256-KKvxElsvHuNN14cm7RjvVrJDqA1BgYgif7eqcRi6ybc=";
      };
    };
    "micromatch" = {
      out_path = "micromatch";
      binaries = {
      };
      pkg = fetchurl {
        name = "micromatch@4.0.8";
        url = "https://registry.npmjs.org/micromatch/-/micromatch-4.0.8.tgz";
        hash = "sha256-x40F8gr/O8lfAksZ5T9ymYz8X2nEKxJsSd/4sbJxvEY=";
      };
    };
    "mri" = {
      out_path = "mri";
      binaries = {
      };
      pkg = fetchurl {
        name = "mri@1.2.0";
        url = "https://registry.npmjs.org/mri/-/mri-1.2.0.tgz";
        hash = "sha256-njg5kJY8FnuqRxzJKKAusXaO4v+r+BFwgqCbUhP5pXU=";
      };
    };
    "node-addon-api" = {
      out_path = "node-addon-api";
      binaries = {
      };
      pkg = fetchurl {
        name = "node-addon-api@7.1.1";
        url = "https://registry.npmjs.org/node-addon-api/-/node-addon-api-7.1.1.tgz";
        hash = "sha256-sQRV0VqXfAzRehyw62eeA9k5+O+NQwLrM+H3jazHH4I=";
      };
    };
    "picocolors" = {
      out_path = "picocolors";
      binaries = {
      };
      pkg = fetchurl {
        name = "picocolors@1.1.1";
        url = "https://registry.npmjs.org/picocolors/-/picocolors-1.1.1.tgz";
        hash = "sha256-067bKAeWe36zf9EbA7fjcB5yWvecZQ0IEv9yE/j4gtk=";
      };
    };
    "picomatch" = {
      out_path = "picomatch";
      binaries = {
      };
      pkg = fetchurl {
        name = "picomatch@2.3.1";
        url = "https://registry.npmjs.org/picomatch/-/picomatch-2.3.1.tgz";
        hash = "sha256-GxTunshnwJDXtSx3GT2D53kQVTs9GLL4bdK3tV6CwR8=";
      };
    };
    "preline" = {
      out_path = "preline";
      binaries = {
      };
      pkg = fetchurl {
        name = "preline@2.7.0";
        url = "https://registry.npmjs.org/preline/-/preline-2.7.0.tgz";
        hash = "sha256-q403TWQSjHgPrg8zRtW3LyOd4WY+aEeO08Yc85WhaXc=";
      };
    };
    "tailwindcss" = {
      out_path = "tailwindcss";
      binaries = {
      };
      pkg = fetchurl {
        name = "tailwindcss@4.1.3";
        url = "https://registry.npmjs.org/tailwindcss/-/tailwindcss-4.1.3.tgz";
        hash = "sha256-C8GQMwaG6quO9mKBc/YHXBv+ZZpjT7qCl2XPK6Zv4O4=";
      };
    };
    "tapable" = {
      out_path = "tapable";
      binaries = {
      };
      pkg = fetchurl {
        name = "tapable@2.2.1";
        url = "https://registry.npmjs.org/tapable/-/tapable-2.2.1.tgz";
        hash = "sha256-s77+kT0Z5+Io+JBfxtxW84Qv+Xo8OnH7HYxCtT4YRuI=";
      };
    };
    "to-regex-range" = {
      out_path = "to-regex-range";
      binaries = {
      };
      pkg = fetchurl {
        name = "to-regex-range@5.0.1";
        url = "https://registry.npmjs.org/to-regex-range/-/to-regex-range-5.0.1.tgz";
        hash = "sha256-pjolXsIKk1w6MzSmZF+SRds24jGd3x62vA6rK/5bQtc=";
      };
    };
    "typescript" = {
      out_path = "typescript";
      binaries = {
        "tsc" = "../typescript/bin/tsc";
        "tsserver" = "../typescript/bin/tsserver";
      };
      pkg = fetchurl {
        name = "typescript@5.7.3";
        url = "https://registry.npmjs.org/typescript/-/typescript-5.7.3.tgz";
        hash = "sha256-gM/KElS6uOgdY5F45C1kBthW+6bjTK1g0atQ7m5ffrs=";
      };
    };
    "undici-types" = {
      out_path = "undici-types";
      binaries = {
      };
      pkg = fetchurl {
        name = "undici-types@6.20.0";
        url = "https://registry.npmjs.org/undici-types/-/undici-types-6.20.0.tgz";
        hash = "sha256-coyp/P9nY3Lk3NZIteJvu9sonsK89nXnYBzCE0pejW4=";
      };
    };
  };

  # Build the node modules directory
  nodeModules = runCommand "node-modules" {
    nativeBuildInputs = [ 
      libarchive 
      makeWrapper
    ];
  } ''
    mkdir -p $out/node_modules/.bin

    # Extract a given package to it's destination
    extract() {
      local pkg=$1
      local dest=$2
      
      mkdir -p "$dest"
      bsdtar --extract \
        --file "$pkg" \
        --directory "$dest" \
        --strip-components=1 \
        --no-same-owner \
        --no-same-permissions

      chmod -R a+X "$dest"
    }

    # Process each package
    ${lib.concatStringsSep "\n" (lib.mapAttrsToList (name: pkg: ''
      echo "Installing package ${name}..."

      mkdir -p "$out/node_modules/${pkg.out_path}"
      extract "${pkg.pkg}" "$out/node_modules/${pkg.out_path}"
      
      # Handle binaries if they exist
      ${lib.concatStringsSep "\n" (lib.mapAttrsToList (binName: binPath: ''
        ln -sf "${binPath}" "$out/node_modules/.bin/${binName}"
      '') pkg.binaries)}
    '') packages)}

    # Force bun instead of node for script execution
    makeWrapper ${bun}/bin/bun $out/bin/node
    export PATH="$out/bin:$PATH"

    patchShebangs $out/node_modules
  '';

in {
  inherit nodeModules packages;
}
