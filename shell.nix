{ pkgs ? import ( fetchTarball "https://github.com/NixOS/nixpkgs/archive/4ecab3273592f27479a583fb6d975d4aba3486fe.tar.gz" ) {} }:

pkgs.mkShell {

  buildInputs = [
    pkgs.wget
    pkgs.curl
    pkgs.jq
    pkgs.binutils
    
    pkgs.which
    pkgs.htop
    pkgs.zlib
    
    pkgs.git
    pkgs.httpie
    pkgs.mysql80
    pkgs.gping
    pkgs.helix

    pkgs.daemonize
    pkgs.docker
    pkgs.postgresql_15
    pkgs.goose
    pkgs.go
    
  ];
    
  shellHook = ''
    echo "Entering the dev environment"
    sudo env "PATH=$PATH" daemonize -u root `which dockerd`     
    sudo usermod -aG docker $USER    
     
    go version
    
    # exit hook
    trap "sudo killall dockerd" EXIT
    '';

  
}
