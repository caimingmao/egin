package main

import (
    "github.com/daodao97/egin/config"
    "github.com/daodao97/egin/pkg"
)

func main() {
    boot := pkg.Bootstrap{
        RouteMap: config.Routes(),
    }
    boot.Start()
}
