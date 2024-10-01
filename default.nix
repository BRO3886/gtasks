{ lib, buildGoModule }:

buildGoModule {
  pname = "gtasks";
  version = "latest"; # Or specify a particular version

  # Fetching the source from GitHub repository
  src = ./.;

  vendorHash = "sha256-eZfgB91pTORpwO5uTMCiVxRP6BlVwQDMzpfScfRUkWI=";
  # Optional, add meta information
  meta = with lib; {
    description = "A command-line task manager written in Go";
    homepage = "https://github.com/BRO3886/gtasks";
    license = licenses.mit;
    maintainers = with maintainers; [ niksingh710 ];
  };
}

