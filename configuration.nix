{
  services.tor.enable = true;
  services.tor.client.enable = false; # required for HiddenServiceSingleHopMode
  services.tor.relay.onionServices = {
    "nix-locate" = {
      # get hostname: sudo cat /var/lib/tor/onion/nix-locate/hostname
      map = [{ port = 80; target = { port = 8080; }; }];
      version = 3;
      settings = {
        # NON ANONYMOUS hidden service. use tor only for NAT punching
        # FIXME this requires tor.client.enable = false https://github.com/NixOS/nixpkgs/pull/48625
        HiddenServiceSingleHopMode = true;
        HiddenServiceNonAnonymousMode = true;
        SocksPort = 0;
      };
    };
  };
}
