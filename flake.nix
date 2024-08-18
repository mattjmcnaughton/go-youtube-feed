{
  inputs = {
    nixpkgs.url = github:nixos/nixpkgs/nixos-24.05;
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            # allowUnfree = true;
          };
        };

        customNeovim = pkgs.neovim.override {
          configure = {
            customRC = ''
              set number
              syntax enable
              filetype plugin indent on

              set expandtab tabstop=2 shiftwidth=2 softtabstop=2

              augroup filetype_specific_indent
                autocmd!

                " Set proper tabs for Golang.
                autocmd FileType go setlocal noexpandtab tabstop=8 softtabstop=0 shiftwidth=8
                " Run `go fmt` on save. Depends on `vim-go` being installed.
                autocmd BufWritePre *.go :silent! :GoFmt
              augroup END

            '';

            packages.myVimPackage = with pkgs.vimPlugins; {
              start = [
                vim-go
                fzf-vim
                vim-commentary
              ];
            };
          };
        };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = [
            pkgs.just
            pkgs.go
            customNeovim
          ];

          shellHook = ''
            export EDITOR="nvim"
            alias vim="nvim"
            echo "Welcome to the dev env for go-youtube-feed!"
          '';
        };
      });
}
