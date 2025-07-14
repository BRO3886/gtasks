{ lib, buildGoModule, fetchFromGitHub, installShellFiles }:

buildGoModule {
  pname = "gtasks";
  version = "latest"; # Or specify a particular version

  # Fetching the source from GitHub repository
  src = fetchFromGitHub {
    owner = "BRO3886";
    repo = "gtasks";
    rev = "0.10.0";
    sha256 = "sha256-t+D++0YKZhTguZAv3bHqNt71CgSIboAy6gDgPKsV/JE=";
  };

  nativeBuildInputs = [
    installShellFiles
  ];

  postInstall = ''
    installShellCompletion --cmd gtasks \
      --bash <($out/bin/gtasks completion bash) \
      --fish <($out/bin/gtasks completion fish) \
      --zsh <($out/bin/gtasks completion zsh)
  '';
  vendorHash = "sha256-eZfgB91pTORpwO5uTMCiVxRP6BlVwQDMzpfScfRUkWI=";
  # Optional, add meta information
  meta = with lib; {
    description = "A command-line task manager written in Go";
    homepage = "https://github.com/BRO3886/gtasks";
    license = licenses.mit;
    maintainers = with maintainers; [
      niksingh710
    ];
  };
}
