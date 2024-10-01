{
  installShellFiles,
  buildGoModule,
  lib,
}:
buildGoModule rec {
  pname = "nao";
  version = "3.3.0";

  src = builtins.path {
    name = "nao";
    path = ./.;
  };

  vendorHash = "sha256-vpZt2SYjS6KsVA+Vee+r2UvRpDDGv0XoF0DHW9VGsbQ=";
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
