{
  installShellFiles,
  buildGoModule,
  lib,
}:
buildGoModule rec {
  pname = "nao";
  version = "3.2.3";

  src = builtins.path {
    name = "nao";
    path = ./.;
  };

  vendorSha256 = "sha256-MTVJWksGWva+Xet+T2aIOXzkxB7w9raJVwa/p1bwkOo=";
  doCheck = false;

  buildTarget = "./cmd/nao";
  ldflags = ["-X main.version=${version}"];

  nativeBuildInputs = [
    installShellFiles
  ];

  postInstall = ''
    installShellCompletion --cmd nao \
      --bash <($out/bin/nao --debug completion bash) \
      --fish <($out/bin/nao --debug completion fish) \
      --zsh <($out/bin/nao --debug completion zsh)
  '';

  meta = with lib; {
    description = "A CLI tool to take notes without worrying about the path where the file is";
    homepage = "https://github.com/luisnquin/nao";
    license = licenses.mit;
    maintainers = with maintainers; [luisnquin];
  };
}
