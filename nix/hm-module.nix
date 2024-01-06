self: {
  config,
  pkgs,
  lib,
  ...
}:
with lib; let
  inherit (pkgs.stdenv.hostPlatform) system;
  cfg = config.programs.nao;
in {
  options = {
    programs.nao = {
      enable = mkEnableOption "nao";
      config = let
        editorModule = types.submodule {
          options = {
            name = {
              type = types.str;
              default = "nano";
            };

            extraArgs = {
              type = types.listOf types.str;
              default = [];
            };
          };
        };

        configModule = types.submodule {
          options = {
            editor = mkOption {
              type = editorModule;
              default = "nano";
            };

            theme = mkOption {
              type = types.enum [
                "default"
                "beach-day"
                "party"
                "nord"
                "no-theme"
                "rose-pine"
                "rose-pine-dawn"
                "rose-pine-moon"
              ];
              default = "default";
            };

            readOnlyOnConflict = mkOption {
              type = types.bool;
              default = false;
            };
          };
        };
      in
        mkOption {
          type = configModule;
          default = {};
        };
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [
      self.packages.${system}.default
    ];

    xdg.configFile = {
      "nao/config.yml".source = (pkgs.formats.yaml {}).generate "nao-config" cfg.config;
    };
  };
}
